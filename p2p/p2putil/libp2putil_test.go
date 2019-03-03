/*
 * @file
 * @copyright defined in aergo/LICENSE.txt
 */

package p2putil

import (
	"net"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/aergoio/aergo/p2p/p2pcommon"
	peer "github.com/libp2p/go-libp2p-peer"

	"github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/assert"
)

func TestFromMultiAddr(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		wantIp   net.IP // verify one of them
		wantPort int
		wantErr  bool
	}{
		{"TIP4peerAddr", "/ip4/192.168.0.58/tcp/11002/p2p/16Uiu2HAmHuBgtnisgPLbujFvxPNZw3Qvpk3VLUwTzh5C67LAZSFh", net.ParseIP("192.168.0.58"), 11002, false},
		{"TIP4MissingID", "/ip4/192.168.0.58/tcp/11002", net.ParseIP("192.168.0.58"), -1, true},
		{"TIP4MissingPort", "/ip4/192.168.0.58/p2p/16Uiu2HAmHuBgtnisgPLbujFvxPNZw3Qvpk3VLUwTzh5C67LAZSFh", net.ParseIP("192.168.0.58"), 11002, true},
		{"TIP6peerAddr", "/ip6/FE80::0202:B3FF:FE1E:8329/tcp/11003/p2p/16Uiu2HAmHuBgtnisgPLbujFvxPNZw3Qvpk3VLUwTzh5C67LAZSFh", net.ParseIP("FE80::0202:B3FF:FE1E:8329"), 11003, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ma, _ := ParseMultiaddrWithResolve(tt.str)
			got, err := FromMultiAddr(ma)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromMultiAddr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				ip := net.ParseIP(got.IPAddress)
				if !reflect.DeepEqual(ip, tt.wantIp) {
					t.Errorf("FromMultiAddr() = %v, want %v", ip.String(), tt.wantIp)
				}
				if !reflect.DeepEqual(got.Port, uint32(tt.wantPort)) {
					t.Errorf("FromMultiAddr() = %v, want %v", got.Port, tt.wantPort)
				}
			}

			got2, err := ParseMultiAddrString(tt.str)
			if !reflect.DeepEqual(got2, got) {
				t.Errorf("result of FromMultiAddr and ParseMultiAddrString differ %v, want %v", got2, got)
			}
		})
	}
}

func TestPeerMeta_ToMultiAddr(t *testing.T) {
	type fields struct {
		IPAddress string
		Port      uint32
		ID        peer.ID
	}
	tests := []struct {
		name         string
		fields       fields
		wantErr      bool
		wantProtocol int
	}{
		{"t1", fields{"192.168.1.2", 2, "id0002"}, false, multiaddr.P_IP4},
		{"t2", fields{"0.0.0.0", 2223, "id2223"}, false, multiaddr.P_IP4},
		{"t3", fields{"2001:0db8:85a3:08d3:1319:8a2e:0370:7334", 444, "id0002"}, false, multiaddr.P_IP6},
		{"t4", fields{"::ffff:192.0.1.2", 444, "id0002"}, false, multiaddr.P_IP4},
		//{"t5", fields{"www.aergo.io", 444, "id0002"}, false, multiaddr.P_IP4},
		{"t6", fields{"no1.blocko.com", 444, "id0002"}, true, multiaddr.P_IP4},
		{"tErr", fields{"dw::ffff:192.0.1.2", 444, "id0002"}, true, multiaddr.P_IP4},
		{"tErr2", fields{"dw!.com", 444, "id0002"}, true, multiaddr.P_IP4},
		{"tErr3", fields{".google.com", 444, "id0002"}, true, multiaddr.P_IP4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := p2pcommon.PeerMeta{
				ID:        tt.fields.ID,
				IPAddress: tt.fields.IPAddress,
				Port:      tt.fields.Port,
			}
			got, err := PeerMetaToMultiAddr(m)
			if (err != nil) != tt.wantErr {
				t.Errorf("PeerMeta.PeerMetaToMultiAddr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				val, err := got.ValueForProtocol(tt.wantProtocol)
				assert.Nil(t, err)
				expIp := net.ParseIP(tt.fields.IPAddress)
				if expIp == nil {
					ips, _ := ResolveHostDomain(tt.fields.IPAddress)
					expIp = ips[0]
				}
				assert.Equal(t, expIp, net.ParseIP(val))
			}
		})
	}
}

func TestParseMultiaddrWithResolve(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		wantIps  []net.IP // verify one of them
		wantPort string
		wantErr  bool
	}{
		{"TIP4peerAddr", "/ip4/192.168.0.58/tcp/11002/p2p/16Uiu2HAmHuBgtnisgPLbujFvxPNZw3Qvpk3VLUwTzh5C67LAZSFh", []net.IP{net.ParseIP("192.168.0.58")}, "11002", false},
		{"TIP4AndPort", "/ip4/192.168.0.58/tcp/11002", []net.IP{net.ParseIP("192.168.0.58")}, "11002", false},
		{"TIP6peerAddr", "/ip6/FE80::0202:B3FF:FE1E:8329/tcp/11003/p2p/16Uiu2HAmHuBgtnisgPLbujFvxPNZw3Qvpk3VLUwTzh5C67LAZSFh", []net.IP{net.ParseIP("FE80::0202:B3FF:FE1E:8329")}, "11003", false},
		//FIXME
		// skip case, since it depend on external environment. uncomment it if really need, and comment back after finishing test
		//{"TDomainName", "/dns/aergo.io/tcp/11004/p2p/16Uiu2HAmHuBgtnisgPLbujFvxPNZw3Qvpk3VLUwTzh5C67LAZSFh", []net.IP{net.ParseIP("104.20.161.59"), net.ParseIP("104.20.160.59")}, "11004", false},
		{"TInvalidDomain", "/dns/nowhere.a.aergo.io/tcp/11004/p2p/16Uiu2HAmHuBgtnisgPLbujFvxPNZw3Qvpk3VLUwTzh5C67LAZSFh", []net.IP{}, "", true},
		{"TWrongProto", "/ipx/192.168.0.58/tcp/11004/p2p/16Uiu2HAmHuBgtnisgPLbujFvxPNZw3Qvpk3VLUwTzh5C67LAZSFh", []net.IP{}, "", true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMultiaddrWithResolve(tt.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMultiaddrWithResolve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Errorf("ParseMultiaddrWithResolve() is nil , want not nil")
				}

				ipStr, err := got.ValueForProtocol(multiaddr.P_IP4)
				if err != nil {
					ipStr, _ = got.ValueForProtocol(multiaddr.P_IP6)
				}
				ip := net.ParseIP(ipStr)
				found := false
				for _, wantIp := range tt.wantIps {
					if reflect.DeepEqual(ip, wantIp) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("ParseMultiaddrWithResolve() = %v, want %v", ip, tt.wantIps)
				}
				port, err := got.ValueForProtocol(multiaddr.P_TCP)
				if !reflect.DeepEqual(port, tt.wantPort) {
					t.Errorf("ParseMultiaddrWithResolve() = %v, want %v", port, tt.wantPort)
				}
			}
		})
	}
}

func TestGenerateKeyFile(t *testing.T) {
	// this test should not be run by root
	testDir := "_tmp"
	existDir := "_tmp/holder.key"
	if _, err := os.Stat(testDir); !os.IsNotExist(err) {
		t.Skip("can't test since test dir already exists")
	} else {
		err := os.MkdirAll(testDir, os.ModePerm)
		if err != nil {
			t.Skip("can't test. permission error. " + err.Error())
		}
		defer func() {
			os.RemoveAll(testDir)
		}()
		err = os.MkdirAll(existDir, os.ModePerm)
		if err != nil {
			t.Skip("can't test. permission error. " + err.Error())
		}
	}
	type args struct {
		dir    string
		prefix string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"TSucc", args{filepath.Join(testDir, "abc"), "testkey"}, false},
		{"TPermission", args{"/sbin/abc", "testkey"}, true},
		{"TDir", args{testDir, "holder"}, true},
		//{"TNotExist", args{}, true},
		//{"TInvalidKey", args{}, false},
		//{"TInvalidKey", args{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPriv, gotPub, err := GenerateKeyFile(tt.args.dir, tt.args.prefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateKeyFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			file := filepath.Join(tt.args.dir, tt.args.prefix+p2pcommon.DefaultPkKeyExt)
			if !tt.wantErr {
				if !gotPriv.GetPublic().Equals(gotPub) {
					t.Errorf("priv key and pub key check failed")
					return
				}

				ldPriv, ldPub, actErr := LoadKeyFile(file)
				if actErr != nil {
					t.Errorf("LoadKeyFile() should not return error, but get %v, ", actErr)
				}
				if !ldPriv.Equals(gotPriv) {
					t.Errorf("GenerateKeyFile() and LoadKeyFile() private key is differ.")
					return
				}
				if !ldPub.Equals(gotPub) {
					t.Errorf("GenerateKeyFile() and LoadKeyFile() public key is differ.")
					return
				}
			}
		})
	}
}

func TestGetIPFromMultiAddr(t *testing.T) {
	ip4 := net.ParseIP("211.34.56.78")
	ip6 := net.ParseIP("fe80::dcbf:beff:fe87:e30a")
	ip6_2 := net.ParseIP("::ffff:192.0.1.2")
	if ip4 == nil || ip6 == nil || ip6_2 == nil {

	}
	tests := []struct {
		name string
		addrStr string

		want net.IP
	}{
		{"TIP4","/ip4/211.34.56.78",ip4},
		{"TIP6","/ip6/fe80::dcbf:beff:fe87:e30a",ip6},
		{"TIP4withPort","/ip4/211.34.56.78/tcp/1234",ip4},
		{"TIP6withPort","/ip6/fe80::dcbf:beff:fe87:e30a/tcp/1234/p2p/16Uiu2HAmHuBgtnisgPLbujFvxPNZw3Qvpk3VLUwTzh5C67LAZSFh",ip6},
		{"TIP6_2","/ip6/::ffff:192.0.1.2/tcp/1234",ip6_2},
		{"TIP4overIP6","/ip6/211.34.56.78/tcp/1234",ip4},
		{"TNoIP","/https/tcp/1234",nil},
		{"TNoIP2","/unix/abcde/tcp/1234",nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr, err  := multiaddr.NewMultiaddr(tt.addrStr)
			if err != nil {
				t.Errorf("GetIPFromMultiAddr() test is wrong (addr %v, err %v)", addr, err.Error())
				t.FailNow()
			}
			if got := GetIPFromMultiAddr(addr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetIPFromMultiAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}
