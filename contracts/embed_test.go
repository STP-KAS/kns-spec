package contracts

import (
	"encoding/json"
	"testing"
)

func TestKasNameArtifact(t *testing.T) {
	if len(Artifact) < 100 {
		t.Fatal("missing KasName.json")
	}
	var art struct {
		Contracts map[string]struct {
			Compiled struct {
				TemplateHash []byte `json:"template_hash"`
			} `json:"compiled"`
		} `json:"contracts"`
	}
	if err := json.Unmarshal(Artifact, &art); err != nil {
		t.Fatal(err)
	}
	h := art.Contracts["KasName"].Compiled.TemplateHash
	if len(h) != 32 {
		t.Fatalf("template_hash len %d", len(h))
	}
}
