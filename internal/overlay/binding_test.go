package overlay

import (
	"strings"
	"testing"
)

func TestBinding(t *testing.T) {
	own := strings.Repeat("ab", 32)
	noise := strings.Repeat("cd", 32)
	m, err := BindingMessage("Alice.KAS", own, noise, 1, 2)
	if err != nil {
		t.Fatal(err)
	}
	want := "kns-session/v1\nname=alice.kas\nowner=" + own + "\nnoise=" + noise + "\nseq=1\nexp=2"
	if m != want {
		t.Fatal(m)
	}
	if _, err := BindingMessage("a\nnoise=ff", own, noise, 1, 2); err == nil {
		t.Fatal("newline")
	}
}

func TestParseURI(t *testing.T) {
	n, p, ok := ParseURI("kns://Alice.kas/pay")
	if !ok || n != "alice.kas" || p != "/pay" {
		t.Fatal(n, p, ok)
	}
	n, p, ok = ParseURI("kns://bob.kas")
	if !ok || n != "bob.kas" || p != "" {
		t.Fatal(n, p)
	}
	if _, _, ok := ParseURI("https://x"); ok {
		t.Fatal("https")
	}
}

func TestParseCap(t *testing.T) {
	tok, ok := ParseCap("/cap/abc")
	if !ok || tok != "abc" {
		t.Fatal(tok, ok)
	}
}

func TestInvoiceURI(t *testing.T) {
	u := InvoiceURI("kaspa:qq", "alice.kas", "/api&amount=9", "100000000")
	if u != "kaspa:qq?amount=1&label=alice.kas&message=%2Fapi%26amount%3D9" {
		t.Fatal(u)
	}
}
