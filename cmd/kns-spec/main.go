package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/STP-KAS/kns-spec/internal/envelope"
	"github.com/STP-KAS/kns-spec/internal/kns"
	"github.com/STP-KAS/kns-spec/internal/proof"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		usage()
		os.Exit(2)
	}
	switch args[0] {
	case "prove":
		prove()
	case "envelope":
		if len(args) < 3 {
			fail("usage: kns-spec envelope create|transfer|list|send …")
		}
		envelop(args[1], args[2:])
	case "resolve":
		if len(args) < 2 {
			fail("usage: kns-spec resolve name.kas")
		}
		resolve(args[1])
	default:
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprint(os.Stderr, `kns-spec — implementer kit for KNS inscriptions + proven covenants

  kns-spec prove
  kns-spec envelope create <label>
  kns-spec envelope transfer <inscriptionId> <kaspa:addr>
  kns-spec envelope list <inscriptionId>
  kns-spec envelope send <inscriptionId>
  kns-spec resolve <name.kas>
`)
}

func prove() {
	f, err := proof.Load("")
	if err != nil {
		fail(err.Error())
	}
	results := proof.Verify(f)
	bad := 0
	for _, r := range results {
		mark := "ok"
		if !r.OK {
			mark = "FAIL"
			bad++
		}
		fmt.Printf("%s  %-12s  %-28s  %s\n", mark, r.Kind, r.Name, r.Detail)
		if r.URL != "" {
			fmt.Printf("    %s\n", r.URL)
		}
	}
	if bad > 0 {
		os.Exit(1)
	}
}

func envelop(op string, args []string) {
	var (
		b   []byte
		err error
	)
	switch op {
	case "create":
		b, err = envelope.Create(args[0])
	case "transfer":
		if len(args) < 2 {
			fail("envelope transfer <id> <to>")
		}
		b, err = envelope.Transfer(args[0], args[1])
	case "list":
		b, err = envelope.List(args[0])
	case "send":
		b, err = envelope.Send(args[0])
	default:
		fail("unknown op " + op)
	}
	if err != nil {
		fail(err.Error())
	}
	fmt.Println(string(b))
	fmt.Println(envelope.ScriptSketch(b))
	if op == "create" {
		fmt.Printf("price %d KAS → %s (reveal output 0)\n", envelope.PriceKAS(args[0]), envelope.FeeAddress("mainnet"))
		fmt.Println("KasWare: buildScript({ type: \"KNS\", data: <json> }) then submitCommitReveal")
	}
}

func resolve(name string) {
	if !strings.HasSuffix(strings.ToLower(name), ".kas") {
		name += ".kas"
	}
	c := kns.New("")
	own, err := c.Owner(name)
	if err != nil {
		fail(err.Error())
	}
	a, _ := c.Asset(name)
	out := map[string]any{
		"name":  own.Asset,
		"owner": own.Owner,
		"id":    own.AssetID,
		"warn":  "Indexer resolution. Verify the address before sending. Not consensus.",
	}
	if a != nil {
		out["txid"] = a.TransactionID
		out["explorer"] = "https://explorer.kaspa.org/txs/" + a.TransactionID
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(out)
}

func fail(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
