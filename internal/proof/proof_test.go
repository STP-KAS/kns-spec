package proof

import "testing"

func TestProofsFile(t *testing.T) {
	f, err := Load("../../proofs/proofs.json")
	if err != nil {
		f, err = Load("")
	}
	if err != nil {
		t.Fatal(err)
	}
	if len(f.Inscriptions) == 0 || len(f.Covenants) == 0 {
		t.Fatal("empty proofs")
	}
	if f.Inscriptions[0].Name != "kns.kas" {
		t.Fatal(f.Inscriptions[0].Name)
	}
	if f.Covenants[0].Name != "bakery.kas" {
		t.Fatal(f.Covenants[0].Name)
	}
}

func TestLiveProofs(t *testing.T) {
	if testing.Short() {
		t.Skip("live")
	}
	f, err := Load("../../proofs/proofs.json")
	if err != nil {
		f, err = Load("")
	}
	if err != nil {
		t.Fatal(err)
	}
	results := Verify(f)
	fail := 0
	for _, r := range results {
		if !r.OK {
			fail++
			t.Errorf("%s %s: %s", r.Kind, r.Name, r.Detail)
			continue
		}
		t.Logf("ok %s %s %s", r.Kind, r.Name, r.Detail)
	}
	if fail > 0 {
		t.Fatalf("%d proof(s) failed", fail)
	}
}
