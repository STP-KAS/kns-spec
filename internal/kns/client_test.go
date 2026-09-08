package kns

import "testing"

func TestLiveKnsKas(t *testing.T) {
	if testing.Short() {
		t.Skip("live")
	}
	c := New("")
	own, err := c.Owner("kns.kas")
	if err != nil {
		t.Skip(err)
	}
	if own.Owner != "kaspa:qzt9yuqceqvt2vk9dz7ddzayaa5flnenkymec59xvzm55ln3k72vgecxjhnjp" {
		t.Fatalf("owner %s", own.Owner)
	}
	if own.AssetID != "223233acf8e5abea9291627a5edd859439c4553100c71309b86599f02414a532i0" {
		t.Fatalf("id %s", own.AssetID)
	}
	a, err := c.Asset("kns.kas")
	if err != nil {
		t.Fatal(err)
	}
	if a.TransactionID != "223233acf8e5abea9291627a5edd859439c4553100c71309b86599f02414a532" {
		t.Fatalf("txid %s", a.TransactionID)
	}
}
