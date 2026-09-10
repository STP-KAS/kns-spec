package envelope

import "testing"

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

func TestRejectSubname(t *testing.T) {
	if _, err := Create("pay.shop"); err == nil {
		t.Fatal("subname must not be a create payload")
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
