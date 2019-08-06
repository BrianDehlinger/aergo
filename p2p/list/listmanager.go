/*
 * @file
 * @copyright defined in aergo/LICENSE.txt
 */

package list

import (
	"errors"
	"github.com/aergoio/aergo-lib/log"
	"github.com/aergoio/aergo/config"
	"github.com/aergoio/aergo/contract/enterprise"
	"github.com/aergoio/aergo/p2p/p2pcommon"
	"github.com/aergoio/aergo/types"
	"net"
	"strings"
	"sync"
	"time"
)

// variables that are used internally
var (
	NotFoundError = errors.New("ban status not found")
	UndefinedTime = time.Unix(0, 0)
	FarawayFuture = time.Date(99999, 1, 1, 0, 0, 0, 0, time.UTC)
)

const (
	localListFile = "list.json"
)

type listManagerImpl struct {
	logger    *log.Logger
	chainAcc  types.ChainAccessor
	prm       p2pcommon.PeerRoleManager
	publicNet bool

	entries []enterprise.WhiteListEntry
	rwLock  sync.RWMutex
	authDir string

	stopScheduler chan interface{}
}

func NewListManager(conf *config.AuthConfig, authDir string, chainAcc types.ChainAccessor, prm p2pcommon.PeerRoleManager, logger *log.Logger, publicNet bool) p2pcommon.ListManager {
	bm := &listManagerImpl{
		logger:    logger,
		chainAcc:  chainAcc,
		prm:       prm,
		publicNet: publicNet,

		authDir:       authDir,
		stopScheduler: make(chan interface{}),
	}

	return bm
}

func (lm *listManagerImpl) Start() {
	lm.logger.Debug().Msg("starting up list manager")

	lm.RefineList()
}

func (lm *listManagerImpl) Stop() {
	lm.logger.Debug().Msg("stopping list manager")
}

func (lm *listManagerImpl) IsBanned(addr string, pid types.PeerID) (bool, time.Time) {
	// malformed ip address is banned
	ip := net.ParseIP(addr)
	if ip == nil {
		return true, FarawayFuture
	}

	// empty entry is
	if len(lm.entries) == 0 {
		return false, FarawayFuture
	}

	// bps are automatically allowed
	if lm.prm.GetRole(pid) == p2pcommon.BlockProducer {
		return false, FarawayFuture
	}

	// finally check peer is in list
	for _, ent := range lm.entries {
		if ent.Contains(ip, pid) {
			return false, FarawayFuture
		}
	}
	return true, FarawayFuture
}

func (lm *listManagerImpl) RefineList() {
	if lm.publicNet {
		lm.logger.Info().Msg("network is public, apply default policy instead (allow all)")
		lm.entries = make([]enterprise.WhiteListEntry, 0)
	} else {
		wl, err := lm.chainAcc.GetEnterpriseConfig(enterprise.P2PWhite)
		if err != nil || len(wl.Values) == 0 {
			lm.logger.Info().Msg("no whitelist found. apply default policy instead (allow all)")
			//ent, _ := NewWhiteListEntry(":")
			//lm.entries = append(lm.entries, ent)
			lm.entries = make([]enterprise.WhiteListEntry, 0)
		} else if !wl.GetOn() {
			lm.logger.Info().Msg("whitelist conf is disabled. apply default policy instead (allow all)")
			lm.entries = make([]enterprise.WhiteListEntry, 0)
		} else {
			entries := make([]enterprise.WhiteListEntry, 0, len(wl.Values))
			for _, v := range wl.Values {
				ent, err := enterprise.NewWhiteListEntry(v)
				if err != nil {
					panic("invalid whitelist entry " + v)
				}
				entries = append(entries, ent)
			}
			lm.entries = entries
			lm.logger.Info().Str("entries", strings.Join(wl.Values, " , ")).Msg("loaded whitelist entries")
		}
	}
}

func (lm *listManagerImpl) Summary() map[string]interface{} {
	// There can be a little error
	sum := make(map[string]interface{})
	entries := make([]string, 0, len(lm.entries))
	for _, e := range lm.entries {
		entries = append(entries, e.String())
	}
	sum["whitelist"] = entries
	return sum
}