package overlay

import "testing"

func TestAppOrder(t *testing.T) {
	r := Records{Website: "https://a", IPFS: "ipfs://b"}
	if r.Web() != "https://a" {
		t.Fatal(r.Web())
	}
	if r.Run() != "ipfs://b" {
		t.Fatal("run ignores https")
	}
	r.Redirect = "https://r"
	if r.App() != "https://r" {
		t.Fatal("redirect first for web")
	}
	if r.Run() != "ipfs://b" {
		t.Fatal("run still cid")
	}
}

func TestSessionPrefersPeer(t *testing.T) {
	r := Records{Onion: "onion://x", Peer: "/ip4/1.2.3.4/tcp/9"}
	if r.Session() != "/ip4/1.2.3.4/tcp/9" {
		t.Fatal(r.Session())
	}
}

func TestURI(t *testing.T) {
	if URI("Alice.KAS", "") != "kns://alice.kas" {
		t.Fatal(URI("Alice.KAS", ""))
	}
	if URI("alice.kas", "pay") != "kns://alice.kas/pay" {
		t.Fatal(URI("alice.kas", "pay"))
	}
}

func TestFromMapNil(t *testing.T) {
	r := FromMap(map[string]any{"website": "https://app.knsdomains.org/", "bio": nil, "x": "knsdomain"})
	if r.App() != "https://app.knsdomains.org/" {
		t.Fatal(r.App())
	}
}

func TestPayURI(t *testing.T) {
	if PayURI("kaspa:qq") != "kaspa:qq" || PayURI("qq") != "kaspa:qq" {
		t.Fatal("pay uri")
	}
}

func TestParse(t *testing.T) {
	r, err := Parse([]byte(`{"kas":"kaspa:qq","ipfs":"ipfs://cid","peer":"/ip4/10.0.0.1/tcp/4001"}`))
	if err != nil {
		t.Fatal(err)
	}
	if r.PayAddress() != "kaspa:qq" || r.App() != "ipfs://cid" || r.Session() == "" {
		t.Fatalf("%+v", r)
	}
}
