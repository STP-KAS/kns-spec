// Package overlay is the name-addressed layer above KNS inscriptions.
// Kaspa settles. The name locates. The user machine runs the dApp.
package overlay

import (
	"encoding/json"
	"strings"
)

const (
	KeyKAS      = "kas"
	KeyPay      = "pay"
	KeyIPFS     = "ipfs"
	KeyKFS      = "kfs"
	KeyArweave  = "arweave"
	KeyContent  = "contenthash"
	KeyWebsite  = "website"
	KeyRedirect = "redirectUrl"
	KeyPeer     = "peer"
	KeyOnion    = "onion"
	KeyNoise    = "noise"
	KeyAgent    = "agent"
	KeyLane     = "lane"
	KeyVault    = "vault"
	KeyVaultCom = "vaultCommit"
)

type Records struct {
	KAS         string `json:"kas,omitempty"`
	Pay         string `json:"pay,omitempty"`
	IPFS        string `json:"ipfs,omitempty"`
	KFS         string `json:"kfs,omitempty"`
	Arweave     string `json:"arweave,omitempty"`
	ContentHash string `json:"contenthash,omitempty"`
	Website     string `json:"website,omitempty"`
	Redirect    string `json:"redirectUrl,omitempty"`
	Peer        string `json:"peer,omitempty"`
	Onion       string `json:"onion,omitempty"`
	Noise       string `json:"noise,omitempty"`
	Agent       string `json:"agent,omitempty"`
	Lane        string `json:"lane,omitempty"`
	Vault       string `json:"vault,omitempty"`
	VaultCommit string `json:"vaultCommit,omitempty"`
}

func FromMap(m map[string]any) Records {
	str := func(k string) string {
		if m == nil {
			return ""
		}
		v, ok := m[k]
		if !ok || v == nil {
			return ""
		}
		s, ok := v.(string)
		if !ok {
			return ""
		}
		return strings.TrimSpace(s)
	}
	return Records{
		KAS:         str(KeyKAS),
		Pay:         str(KeyPay),
		IPFS:        str(KeyIPFS),
		KFS:         str(KeyKFS),
		Arweave:     str(KeyArweave),
		ContentHash: str(KeyContent),
		Website:     str(KeyWebsite),
		Redirect:    str(KeyRedirect),
		Peer:        str(KeyPeer),
		Onion:       str(KeyOnion),
		Noise:       str(KeyNoise),
		Agent:       str(KeyAgent),
		Lane:        str(KeyLane),
		Vault:       str(KeyVault),
		VaultCommit: str(KeyVaultCom),
	}
}

func Parse(raw []byte) (Records, error) {
	var r Records
	if len(raw) == 0 {
		return r, nil
	}
	err := json.Unmarshal(raw, &r)
	return r, err
}

func (r Records) PayAddress() string {
	if s := strings.TrimSpace(r.KAS); s != "" {
		return s
	}
	return strings.TrimSpace(r.Pay)
}

func (r Records) App() string {
	for _, v := range []string{r.Redirect, r.Website, r.IPFS, r.KFS, r.Arweave, r.ContentHash} {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

func (r Records) Session() string {
	if s := strings.TrimSpace(r.Peer); s != "" {
		return s
	}
	if s := strings.TrimSpace(r.Onion); s != "" {
		return s
	}
	return strings.TrimSpace(r.Noise)
}

func (r Records) Private() bool {
	return strings.TrimSpace(r.Vault) != "" || strings.TrimSpace(r.VaultCommit) != ""
}

// ResolveWarning is the official KNS integration-note class of text.
const ResolveWarning = "⚠️ Before sending a transaction, please ensure that the resolved address is correct. Once a transfer is sent, it cannot be reversed. Any input errors may result in the loss of assets. Please verify carefully."

// ProposedKeys are the extra profile keys KNS should index (still text inscriptions).
var ProposedKeys = []string{KeyIPFS, KeyKFS, KeyContent, KeyPeer, KeyOnion, KeyAgent, KeyKAS, KeyNoise}

func PayURI(addr string) string {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return ""
	}
	if strings.HasPrefix(addr, "kaspa:") || strings.HasPrefix(addr, "kaspatest:") {
		return addr
	}
	return "kaspa:" + addr
}

func URI(name, path string) string {
	name = strings.TrimSpace(strings.ToLower(name))
	if !strings.HasSuffix(name, ".kas") {
		name += ".kas"
	}
	path = strings.TrimSpace(path)
	if path == "" || path == "/" {
		return "kns://" + name
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return "kns://" + name + path
}
