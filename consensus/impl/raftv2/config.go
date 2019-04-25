package raftv2

import (
	"errors"
	"fmt"
	"github.com/aergoio/aergo/chain"
	"github.com/aergoio/aergo/config"
	"github.com/libp2p/go-libp2p-peer"
	"net"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	ErrNotIncludedRaftMember = errors.New("this node isn't included in initial raft members")
	ErrInvalidRaftID         = errors.New("invalid raft raftID")
	ErrDupRaftUrl            = errors.New("duplicated raft bp urls")
	ErrRaftEmptyTLSFile      = errors.New("cert or key file name is empty")
	ErrNotHttpsURL           = errors.New("url scheme is not https")
	ErrURLInvalidScheme      = errors.New("url has invalid scheme")
	ErrURLInvalidPort        = errors.New("url must have host:port style")
	ErrInvalidRaftBPID       = errors.New("raft bp raftID is not ordered. raftID must start with 1 and be sorted")
	ErrDupBP                 = errors.New("raft bp description is duplicated")
	ErrInvalidRaftPeerID     = errors.New("peerID of current raft bp is not equals to p2p configure")
)

const (
	DefaultMarginChainDiff = 1
	DefaultTickMS          = time.Millisecond * 30
)

func (bf *BlockFactory) InitCluster(cfg *config.Config) error {
	useTls := true
	var err error

	raftConfig := cfg.Consensus.Raft
	if raftConfig == nil {
		panic("raftconfig is not set. please set raftName, raftBPs.")
	}

	//set default
	if raftConfig.Tick != 0 {
		RaftTick = time.Duration(raftConfig.Tick * 1000000)
	}

	lenBPs := len(raftConfig.BPs)

	chainID, err := chain.Genesis.ID.Bytes()
	if err != nil {
		return err
	}

	bf.bpc = NewCluster(chainID, bf, raftConfig.Name, uint16(lenBPs), chain.Genesis.Timestamp)

	if useTls, err = validateTLS(raftConfig); err != nil {
		logger.Error().Err(err).
			Str("key", raftConfig.KeyFile).
			Str("cert", raftConfig.CertFile).
			Msg("failed to validate tls config for raft")
		return err
	}

	if raftConfig.ListenUrl != "" {
		if err := isValidURL(raftConfig.ListenUrl, useTls); err != nil {
			logger.Error().Err(err).Msg("failed to validate listen url for raft")
			return err
		}
	}

	if err = bf.bpc.AddInitialMembers(raftConfig, useTls); err != nil {
		logger.Error().Err(err).Msg("failed to validate bpurls, bpid config for raft")
		return err
	}

	if err = bf.bpc.SetThisNode(); err != nil {
		return err
	}

	RaftSkipEmptyBlock = raftConfig.SkipEmpty

	logger.Info().Bool("skipempty", RaftSkipEmptyBlock).Int64("rafttick(nanosec)", RaftTick.Nanoseconds()).Float64("interval(sec)", bf.blockInterval.Seconds()).Msg(bf.bpc.toString())

	return nil
}

func validateTLS(raftCfg *config.RaftConfig) (bool, error) {
	if len(raftCfg.CertFile) == 0 && len(raftCfg.KeyFile) == 0 {
		return false, nil
	}

	//두 파일이 모두 설정되어 있는지 확인
	//실제 file에 존재하는지 확인
	if len(raftCfg.CertFile) == 0 || len(raftCfg.KeyFile) == 0 {
		logger.Error().Str("raftcertfile", raftCfg.CertFile).Str("raftkeyfile", raftCfg.KeyFile).
			Msg(ErrRaftEmptyTLSFile.Error())
		return false, ErrRaftEmptyTLSFile
	}

	if len(raftCfg.CertFile) != 0 {
		if _, err := os.Stat(raftCfg.CertFile); err != nil {
			logger.Error().Err(err).Msg("not exist certificate file for raft")
			return false, err
		}
	}

	if len(raftCfg.KeyFile) != 0 {
		if _, err := os.Stat(raftCfg.KeyFile); err != nil {
			logger.Error().Err(err).Msg("not exist Key file for raft")
			return false, err
		}
	}

	return true, nil
}

func isValidURL(urlstr string, useTls bool) error {
	var urlobj *url.URL
	var err error

	if urlobj, err = parseToUrl(urlstr); err != nil {
		logger.Error().Str("url", urlstr).Err(err).Msg("raft bp urlstr is not vaild form")
		return err
	}

	if useTls && urlobj.Scheme != "https" {
		logger.Error().Str("urlstr", urlstr).Msg("raft bp urlstr shoud use https protocol")
		return ErrNotHttpsURL
	}

	return nil
}

func (cl *Cluster) AddInitialMembers(raftCfg *config.RaftConfig, useTls bool) error {
	lenBPs := len(raftCfg.BPs)
	if lenBPs == 0 {
		return fmt.Errorf("config of raft bp is empty")
	}

	// validate each bp
	for _, raftBP := range raftCfg.BPs {
		trimUrl := strings.TrimSpace(raftBP.Url)

		if err := isValidURL(trimUrl, useTls); err != nil {
			return err
		}

		peerID, err := peer.IDB58Decode(raftBP.P2pID)
		if err != nil {
			return fmt.Errorf("invalid raft peerID %s", raftBP.P2pID)
		}

		m := newMember(raftBP.Name, trimUrl, peerID, cl.chainID, cl.chainTimestamp)

		if err := cl.configMembers.add(m, cl.NodeName); err != nil {
			return err
		}
	}

	return nil
}

func (cl *Cluster) SetThisNode() error {
	var member *Member

	if member = cl.configMembers.getMemberByName(cl.NodeName); member == nil {
		return ErrNotIncludedRaftMember
	}

	cl.NodeID = member.ID

	return nil
}

func parseToUrl(urlstr string) (*url.URL, error) {
	var urlObj *url.URL
	var err error

	if urlObj, err = url.Parse(urlstr); err != nil {
		return nil, err
	}

	if urlObj.Scheme != "http" && urlObj.Scheme != "https" {
		return nil, ErrURLInvalidScheme
	}

	if _, _, err := net.SplitHostPort(urlObj.Host); err != nil {
		return nil, ErrURLInvalidPort
	}

	return urlObj, nil
}
