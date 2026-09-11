package overlay

import "testing"

func TestBinding(t *testing.T) {
	m := BindingMessage("Alice.KAS", "ab", "cd", 1, 2)
	if m != "kns-session/v1\nname=alice.kas\nowner=ab\nnoise=cd\nseq=1\nexp=2" {
		t.Fatal(m)
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
	u := InvoiceURI("kaspa:qq", "alice.kas", "/api", "100000000")
	if u != "kaspa:qq?label=alice.kas&message=/api&amount=1" {
		t.Fatal(u)
	}
}
