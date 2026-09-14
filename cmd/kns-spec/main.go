package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/STP-KAS/kns-spec/internal/envelope"
	"github.com/STP-KAS/kns-spec/internal/kns"
	"github.com/STP-KAS/kns-spec/internal/overlay"
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
	case "check":
		if len(args) < 2 {
			fail("usage: kns-spec check name.kas")
		}
		check(args[1])
	case "plan":
		if len(args) < 2 {
			fail("usage: kns-spec plan <label>")
		}
		plan(args[1])
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
	case "overlay":
		if len(args) < 2 {
			fail("usage: kns-spec overlay name.kas")
		}
		overlayCmd(args[1])
	case "primary":
		if len(args) < 2 {
			fail("usage: kns-spec primary kaspa:q…")
		}
		primary(args[1])
	case "vectors":
		vectors()
	case "bind":
		if len(args) < 4 {
			fail("usage: kns-spec bind <name.kas> <owner-xonly-hex> <noise-x25519-hex> [seq] [exp]")
		}
		seq, exp := 1, int64(0)
		if len(args) >= 5 {
			fmt.Sscan(args[4], &seq)
		}
		if len(args) >= 6 {
			fmt.Sscan(args[5], &exp)
		}
		msg, err := overlay.BindingMessage(args[1], args[2], args[3], seq, exp)
		if err != nil {
			fail(err.Error())
		}
		fmt.Println(msg)
	case "serve":
		serve()
	default:
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprint(os.Stderr, `kns-spec — implementer kit for KNS inscriptions + proven covenants

  kns-spec prove
  kns-spec check <name.kas>
  kns-spec plan <label>
  kns-spec envelope create <label>
  kns-spec envelope transfer <inscriptionId> <kaspa:addr>
  kns-spec envelope list <inscriptionId>
  kns-spec envelope send <inscriptionId>
  kns-spec resolve <name.kas>
  kns-spec overlay <name.kas>
  kns-spec primary <kaspa:addr>
  kns-spec vectors
  kns-spec bind <name.kas> <owner-xonly> <noise-x25519> [seq] [exp]
  kns-spec serve          # http://127.0.0.1:8083  (KNS_SPEC_ADDR to override)
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

func check(name string) {
	name = withKas(name)
	c := kns.New("")
	items, err := c.Check([]string{name}, "")
	if err != nil {
		fail(err.Error())
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(items)
}

func plan(label string) {
	payload, err := envelope.Create(label)
	if err != nil {
		fail(err.Error())
	}
	name := withKas(label)
	c := kns.New("")
	items, err := c.Check([]string{name}, "")
	avail := "unknown"
	if err == nil && len(items) > 0 {
		if items[0].Available {
			avail = "available"
		} else {
			avail = "taken"
		}
		if items[0].IsReservedDomain {
			avail += " (reserved)"
		}
	} else if err != nil {
		avail = "check error: " + err.Error()
	}
	fmt.Println(string(payload))
	fmt.Println(envelope.ScriptSketch(payload))
	fmt.Printf("name     %s\n", name)
	fmt.Printf("labelHash %s  (sha256 kns/v1/ + label)\n", envelope.LabelHashHex(label))
	fmt.Printf("price    %d KAS\n", envelope.PriceKAS(label))
	fmt.Printf("wallet   ≥ %.2f KAS (official 5%% extra for domain inscribe)\n", envelope.WalletNeedKAS(envelope.PriceKAS(label), false))
	if c := envelope.Club(label); c != "" {
		fmt.Printf("club     %s\n", c)
	}
	fmt.Printf("fee to   %s (reveal output 0)\n", envelope.FeeAddress("mainnet"))
	fmt.Printf("indexer  %s\n", avail)
	fmt.Println("KasWare  buildScript({ type: \"KNS\", data }) then submitCommitReveal")
	fmt.Println("Kastle   connect(); commitReveal(\"mainnet\", \"kns\", data) — two popups")
}

func withKas(name string) string {
	name = strings.TrimSpace(strings.ToLower(name))
	if !strings.HasSuffix(name, ".kas") {
		name += ".kas"
	}
	return name
}

func resolve(name string) {
	dumpSnapshot(name)
}

func overlayCmd(name string) {
	dumpSnapshot(name)
}

func dumpSnapshot(name string) {
	c := kns.New("")
	s, err := c.Snapshot(name)
	if err != nil {
		fail(err.Error())
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(s)
}

func primary(owner string) {
	c := kns.New("")
	p, err := c.Primary(owner)
	if err != nil {
		fail(err.Error())
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(p)
}

func vectors() {
	raw, err := os.ReadFile(findFile("schemas/vectors.json"))
	if err != nil {
		fail(err.Error())
	}
	os.Stdout.Write(raw)
	if len(raw) == 0 || raw[len(raw)-1] != '\n' {
		fmt.Println()
	}
}

func findFile(rel string) string {
	if _, err := os.Stat(rel); err == nil {
		return rel
	}
	wd, err := os.Getwd()
	if err != nil {
		return rel
	}
	dir := wd
	for i := 0; i < 8; i++ {
		p := filepath.Join(dir, rel)
		if _, err := os.Stat(p); err == nil {
			return p
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return rel
}

func fail(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
