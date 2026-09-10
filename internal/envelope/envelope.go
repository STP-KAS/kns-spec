// Package envelope is the KNS inscription JSON and fee table.
package envelope

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	ProtocolID     = "kns"
	ProtocolDomain = "domain"
	TLD            = "kas"
	MaxLabelRunes  = 255
	MaxInscription = 520
	TextFeeKAS     = 1
)

const (
	MainnetFee = "kaspa:qyp4nvaq3pdq7609z09fvdgwtc9c7rg07fuw5zgeee7xpr085de59eseqfcmynn"
	TN10Fee    = "kaspatest:qq9h47etjv6x8jgcla0ecnp8mgrkfxm70ch3k60es5a50ypsf4h6sak3g0lru"
)

var asciiLabel = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,253}[a-z0-9])?$`)

type Op struct {
	Op string `json:"op"`
	P  string `json:"p,omitempty"`
	V  string `json:"v,omitempty"`
	S  string `json:"s,omitempty"`
	ID string `json:"id,omitempty"`
	To string `json:"to,omitempty"`
}

func Create(label string) ([]byte, error) {
	label = strings.ToLower(strings.TrimSpace(label))
	label = strings.TrimSuffix(label, ".kas")
	if !ValidLabel(label) {
		return nil, fmt.Errorf("invalid label %q", label)
	}
	return json.Marshal(Op{Op: "create", P: ProtocolDomain, V: label})
}

func Transfer(inscriptionID, to string) ([]byte, error) {
	if inscriptionID == "" || to == "" {
		return nil, fmt.Errorf("id and to required")
	}
	return json.Marshal(Op{Op: "transfer", P: ProtocolDomain, ID: inscriptionID, To: to})
}

func List(inscriptionID string) ([]byte, error) {
	if inscriptionID == "" {
		return nil, fmt.Errorf("id required")
	}
	return json.Marshal(Op{Op: "list", P: ProtocolDomain, ID: inscriptionID})
}

func Send(inscriptionID string) ([]byte, error) {
	if inscriptionID == "" {
		return nil, fmt.Errorf("id required")
	}
	return json.Marshal(Op{Op: "send", ID: inscriptionID})
}

func ScriptSketch(payload []byte) string {
	return fmt.Sprintf("<xonly_pubkey> OP_CHECKSIG OP_FALSE OP_IF <%s> <0> <%s> OP_ENDIF", ProtocolID, payload)
}

func PriceKAS(label string) int {
	label = strings.ToLower(strings.TrimSpace(label))
	label = strings.TrimSuffix(label, ".kas")
	n := VisualLength(label)
	switch {
	case n <= 2:
		return 4200
	case n == 3:
		return 2100
	case n == 4:
		return 525
	default:
		return 35
	}
}

func VisualLength(s string) int {
	n := 0
	for _, r := range s {
		if unicode.Is(unicode.Mn, r) || unicode.Is(unicode.Me, r) || r == 0x200D || r == 0xFE0F {
			continue
		}
		n++
	}
	if n == 0 {
		return utf8.RuneCountInString(s)
	}
	return n
}

func ValidLabel(label string) bool {
	if label == "" || utf8.RuneCountInString(label) > MaxLabelRunes {
		return false
	}
	// L1 create `v` may contain dots (official FAQ: "multi-dot domain").
	// Those are flat string assets, not parent-child subnames.
	for _, part := range strings.Split(label, ".") {
		if !validSegment(part) {
			return false
		}
	}
	return true
}

func validSegment(part string) bool {
	if part == "" || utf8.RuneCountInString(part) > MaxLabelRunes {
		return false
	}
	ascii := true
	for i := 0; i < len(part); i++ {
		if part[i] > 127 {
			ascii = false
			break
		}
	}
	if ascii {
		return asciiLabel.MatchString(part)
	}
	for _, r := range part {
		if unicode.IsControl(r) || unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

// Club is the official numeric collection, or "".
// 99 = 0.kas–99.kas (leading zeros excluded except 0.kas).
// 999 = 100.kas–999.kas. 10k = 1000.kas–9999.kas.
func Club(label string) string {
	label = strings.ToLower(strings.TrimSpace(label))
	label = strings.TrimSuffix(label, ".kas")
	if strings.Contains(label, ".") || !numOnly.MatchString(label) {
		return ""
	}
	if label == "0" {
		return "99"
	}
	if strings.HasPrefix(label, "0") {
		return ""
	}
	switch len(label) {
	case 1, 2:
		return "99"
	case 3:
		return "999"
	case 4:
		return "10k"
	default:
		return ""
	}
}

var numOnly = regexp.MustCompile(`^[0-9]+$`)

// WalletNeedKAS is the official extra-balance rule so inscribe does not fail
// "insufficient funds": domain × 1.05, text × 2.
func WalletNeedKAS(priceKAS int, text bool) float64 {
	if text {
		return float64(priceKAS) * 2
	}
	return float64(priceKAS) * 1.05
}

// LabelHash is the 32-byte constructor arg for KasName.sil.
// Not ENS namehash (no recursive parent node). Normalize like KNS: ASCII
// lower-case, no ".kas". Prefix pins the domain so two implementations match.
func LabelHash(label string) [32]byte {
	label = strings.ToLower(strings.TrimSpace(label))
	label = strings.TrimSuffix(label, ".kas")
	sum := sha256.Sum256([]byte("kns/v1/" + label))
	return sum
}

func LabelHashHex(label string) string {
	h := LabelHash(label)
	return fmt.Sprintf("%x", h[:])
}

func FeeAddress(network string) string {
	if strings.EqualFold(network, "tn10") || strings.EqualFold(network, "testnet-10") {
		return TN10Fee
	}
	return MainnetFee
}
