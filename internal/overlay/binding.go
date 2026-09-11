package overlay

import (
	"encoding/hex"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

func hex32(s string) (string, bool) {
	s = strings.ToLower(strings.TrimSpace(s))
	b, err := hex.DecodeString(s)
	if err != nil || len(b) != 32 {
		return "", false
	}
	return s, true
}

// BindingMessage is KIP-5 Schnorr payload. noise is X25519 public hex, not the spend key.
func BindingMessage(name, ownerXonly, noiseX25519 string, seq int, expUnix int64) (string, error) {
	name = strings.ToLower(strings.TrimSpace(name))
	if strings.ContainsAny(name, "\n= ") {
		return "", fmt.Errorf("bad name")
	}
	if !strings.HasSuffix(name, ".kas") {
		name += ".kas"
	}
	owner, ok := hex32(ownerXonly)
	if !ok {
		return "", fmt.Errorf("owner must be 32-byte hex")
	}
	noise, ok := hex32(noiseX25519)
	if !ok {
		return "", fmt.Errorf("noise must be 32-byte hex")
	}
	return fmt.Sprintf("kns-session/v1\nname=%s\nowner=%s\nnoise=%s\nseq=%d\nexp=%d",
		name, owner, noise, seq, expUnix), nil
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
	if u == "" {
		return ""
	}
	q := url.Values{}
	if name != "" {
		q.Set("label", name)
	}
	if resource != "" {
		q.Set("message", resource)
	}
	if sompi != "" {
		n, err := strconv.ParseInt(sompi, 10, 64)
		if err != nil || n < 0 {
			return ""
		}
		q.Set("amount", sompiToKAS(n))
	}
	enc := q.Encode()
	if enc == "" {
		return u
	}
	return u + "?" + enc
}

func sompiToKAS(n int64) string {
	neg := ""
	if n < 0 {
		neg = "-"
		n = -n
	}
	whole := n / 100000000
	frac := n % 100000000
	if frac == 0 {
		return fmt.Sprintf("%s%d", neg, whole)
	}
	s := fmt.Sprintf("%s%d.%08d", neg, whole, frac)
	s = strings.TrimRight(s, "0")
	return strings.TrimRight(s, ".")
}
