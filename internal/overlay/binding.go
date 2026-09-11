package overlay

import (
	"fmt"
	"strconv"
	"strings"
)

// BindingMessage is what a wallet signs with KIP-5 Schnorr (not ECDSA).
// noise is an X25519 public key hex. Never the Kaspa spend key.
func BindingMessage(name, ownerXonly, noiseX25519 string, seq int, expUnix int64) string {
	name = strings.ToLower(strings.TrimSpace(name))
	if !strings.HasSuffix(name, ".kas") {
		name += ".kas"
	}
	return fmt.Sprintf("kns-session/v1\nname=%s\nowner=%s\nnoise=%s\nseq=%d\nexp=%d",
		name, strings.TrimSpace(ownerXonly), strings.TrimSpace(noiseX25519), seq, expUnix)
}

func ParseURI(raw string) (name, path string, ok bool) {
	s := strings.TrimSpace(raw)
	s = strings.TrimPrefix(s, "kns://")
	if s == "" || strings.Contains(s, "://") {
		return "", "", false
	}
	name, path, _ = strings.Cut(s, "/")
	name = strings.ToLower(name)
	if !strings.HasSuffix(name, ".kas") {
		name += ".kas"
	}
	if path != "" {
		path = "/" + path
	}
	return name, path, name != ".kas"
}

func ParseCap(path string) (token string, ok bool) {
	path = strings.TrimPrefix(path, "/")
	if !strings.HasPrefix(path, "cap/") {
		return "", false
	}
	token = strings.TrimPrefix(path, "cap/")
	return token, token != ""
}

func InvoiceURI(payTo, name, resource, sompi string) string {
	u := PayURI(payTo)
	if name == "" && resource == "" && sompi == "" {
		return u
	}
	q := []string{}
	if name != "" {
		q = append(q, "label="+name)
	}
	if resource != "" {
		q = append(q, "message="+resource)
	}
	if sompi != "" {
		if n, err := strconv.ParseInt(sompi, 10, 64); err == nil {
			kas := float64(n) / 1e8
			q = append(q, fmt.Sprintf("amount=%g", kas))
		}
	}
	if len(q) == 0 {
		return u
	}
	return u + "?" + strings.Join(q, "&")
}
