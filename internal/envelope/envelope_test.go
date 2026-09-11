package envelope

import (
	"encoding/json"
	"os"
	"testing"
)

func TestVectorsFile(t *testing.T) {
	raw, err := os.ReadFile("../../schemas/vectors.json")
	if err != nil {
		t.Skip(err)
	}
	var v struct {
		CreateKns    string         `json:"createKns"`
		MainnetFee   string         `json:"mainnetFee"`
		LabelHashKns string         `json:"labelHashKns"`
		Prices       map[string]int `json:"prices"`
	}
	if err := json.Unmarshal(raw, &v); err != nil {
		t.Fatal(err)
	}
	got, err := Create("kns")
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != v.CreateKns {
		t.Fatalf("create %s want %s", got, v.CreateKns)
	}
	if FeeAddress("mainnet") != v.MainnetFee {
		t.Fatal("fee")
	}
	if LabelHashHex("kns") != v.LabelHashKns {
		t.Fatal("hash")
	}
	if PriceKAS("a") != v.Prices["1"] || PriceKAS("alice") != v.Prices["5"] {
		t.Fatal("prices")
	}
}

func TestCreatePayload(t *testing.T) {
	b, err := Create("KNS.kas")
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `{"op":"create","p":"domain","v":"kns"}` {
		t.Fatalf("got %s", b)
	}
}

func TestProvenKnsCreate(t *testing.T) {
	b, err := Create("kns")
	if err != nil {
		t.Fatal(err)
	}
	want := `{"op":"create","p":"domain","v":"kns"}`
	if string(b) != want {
		t.Fatalf("proven kns.kas payload changed: %s", b)
	}
}

func TestOps(t *testing.T) {
	id := "223233acf8e5abea9291627a5edd859439c4553100c71309b86599f02414a532i0"
	to := "kaspa:qzt9yuqceqvt2vk9dz7ddzayaa5flnenkymec59xvzm55ln3k72vgecxjhnjp"
	tr, err := Transfer(id, to)
	if err != nil {
		t.Fatal(err)
	}
	if string(tr) != `{"op":"transfer","p":"domain","id":"`+id+`","to":"`+to+`"}` {
		t.Fatalf("transfer %s", tr)
	}
	ls, err := List("id0")
	if err != nil {
		t.Fatal(err)
	}
	if string(ls) != `{"op":"list","p":"domain","id":"id0"}` {
		t.Fatalf("list %s", ls)
	}
	sd, err := Send("id0")
	if err != nil {
		t.Fatal(err)
	}
	if string(sd) != `{"op":"send","id":"id0"}` {
		t.Fatalf("send %s", sd)
	}
}

func TestGraphemeFamily(t *testing.T) {
	if VisualLength("👨‍👩‍👧‍👦") != 1 {
		t.Fatalf("family %d want 1 (graphemer)", VisualLength("👨‍👩‍👧‍👦"))
	}
	if PriceKAS("👨‍👩‍👧‍👦") != 4200 {
		t.Fatal("family price")
	}
}

func TestPrice(t *testing.T) {
	if PriceKAS("a") != 4200 || PriceKAS("ab") != 4200 {
		t.Fatal("1-2")
	}
	if PriceKAS("abc") != 2100 || PriceKAS("abcd") != 525 {
		t.Fatal("3-4")
	}
	if PriceKAS("alice") != 35 || PriceKAS("example.kas") != 35 {
		t.Fatal("5+")
	}
}

func TestMultiDotL1(t *testing.T) {
	b, err := Create("abc.def")
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `{"op":"create","p":"domain","v":"abc.def"}` {
		t.Fatalf("got %s", b)
	}
	if _, err := Create(".abc"); err == nil {
		t.Fatal("empty segment")
	}
}

func TestClub(t *testing.T) {
	if Club("0") != "99" || Club("7") != "99" || Club("99") != "99" {
		t.Fatal("99")
	}
	if Club("01") != "" || Club("100") != "999" || Club("1000") != "10k" {
		t.Fatal("leading zero / 999 / 10k")
	}
	if Club("abc.def") != "" {
		t.Fatal("multi-dot is not a numeric club")
	}
}

func TestWalletNeed(t *testing.T) {
	if WalletNeedKAS(35, false) != 36.75 {
		t.Fatal("domain 1.05")
	}
	if WalletNeedKAS(1, true) != 2 {
		t.Fatal("text 2x")
	}
}

func TestSketch(t *testing.T) {
	b, _ := Create("example")
	s := ScriptSketch(b)
	if s != `<xonly_pubkey> OP_CHECKSIG OP_FALSE OP_IF <kns> <0> <{"op":"create","p":"domain","v":"example"}> OP_ENDIF` {
		t.Fatal(s)
	}
}

func TestLabelHashStable(t *testing.T) {
	a := LabelHash("KNS.kas")
	b := LabelHash("kns")
	if a != b {
		t.Fatal("normalize")
	}
	if LabelHash("alice") == LabelHash("bob") {
		t.Fatal("distinct")
	}
	if LabelHashHex("kns") != "bf2c1d2ba1a39f872cf89a5cbf0edae491f1a25d50bcb3b983973a98fab6c906" {
		t.Fatalf("pin %s", LabelHashHex("kns"))
	}
}

func TestFee(t *testing.T) {
	if FeeAddress("mainnet") != MainnetFee || FeeAddress("tn10") != TN10Fee {
		t.Fatal("fee")
	}
}
