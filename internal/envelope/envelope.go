// Package envelope is the KNS inscription JSON and fee table.
package envelope

import (
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
	if label == "" || strings.Contains(label, ".") || utf8.RuneCountInString(label) > MaxLabelRunes {
		return false
	}
	ascii := true
	for i := 0; i < len(label); i++ {
		if label[i] > 127 {
			ascii = false
			break
		}
	}
	if ascii {
		return asciiLabel.MatchString(label)
	}
	for _, r := range label {
		if unicode.IsControl(r) || unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

func FeeAddress(network string) string {
	if strings.EqualFold(network, "tn10") || strings.EqualFold(network, "testnet-10") {
		return TN10Fee
	}
	return MainnetFee
}
