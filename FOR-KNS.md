# For the KNS team

**Share this URL:** https://stp-kas.github.io/kns-spec/

Official KNS (not us): [inscriptions](https://kns-2.gitbook.io/kns-docs-1/inscriptions/overview) · [wallets](https://kns-2.gitbook.io/kns-docs-1/supporting-wallet) · [indexer API](https://kns-2.gitbook.io/kns-docs-1/kns-indexer-api) · [simply-kaspa-indexer](https://github.com/supertypo/simply-kaspa-indexer)

Repo: https://github.com/STP-KAS/kns-spec

TN10 snapshot test holder set (addresses + `.kas` names, Testnet-10 only): [snapshot-tn10/](snapshot-tn10/).

How we recommend testing a snapshot / claim with it: [snapshot-tn10/TESTING-SNAPSHOT.md](snapshot-tn10/TESTING-SNAPSHOT.md).

This is a handoff, not a fork of your product and not a claim that covenants already unique `.kas` names.

## What we are asking you to implement

1. **Keep inscribing with KasWare or Kastle.** KasWare: `buildScript({ type: "KNS", data })` then `submitCommitReveal`. Kastle: `commitReveal("mainnet", "kns", data)` — two popups. Reveal output 0 still pays your protocol fee address. Kastle’s high-level `commitReveal` does not document that fee output; the official inscribe tool must attach it.
2. **Keep uniqueness on the indexer.** First valid reveal wins. Consensus will not reject a second `alice.kas`. A KIP-20 `covenant_id` is hashed from an outpoint. It does not encode the label.
3. **Treat a Name UTXO as optional elevation.** `contracts/v1/KasName.sil` is compiled with official **Silverscript v1.0.0** (Ori, 9 Sep 2026, `3ed9733`). Own UTXO only. Continuation keeps the same sompi (fees from a sibling input). Do not `readInputState` a foreign covenant. Argent is Sutton’s multi-actor layer **above** this and is not release-ready. Notes: [BATTLETEST.md](BATTLETEST.md).
4. **Add overlay profile keys** (`ipfs`, `kfs`, `contenthash`, `peer`, `onion`, `agent`, `kas`, `noise`) to Edit Profile and the Profile API. Same 1 KAS text inscriptions as website. Checklist: [CONFORMANCE.md](CONFORMANCE.md). Keys: [PROFILE.md](PROFILE.md). Schema: [schemas/overlay-records.schema.json](schemas/overlay-records.schema.json). Not a hard fork. Not your L2.

## Two objects. Do not mix them.

| Layer | Live today | Uniqueness | Wallet |
| --- | --- | --- | --- |
| Inscription | `kns` commit-reveal envelope | Official indexer, FCFS | KasWare `type: "KNS"` |
| Covenant | P2SH of `KasName.sil` | **Not** consensus. Script-hash lock | Fund / later spend the P2SH |

`trump.kas` and `bakery.kas` already exist as **your** inscriptions on other keys. The P2SH rows we funded are a second object. `opus.kas` has no inscription. Subnames (`opus.dei.kas`) are not an inscription op.

## Proven mainnet

### Inscription (your envelope, on L1)

`kns.kas` reveal:

```
<xonly_pubkey> OP_CHECKSIG OP_FALSE OP_IF <kns> <0> <{"op":"create","p":"domain","v":"kns"}> OP_ENDIF
```

- id: `223233acf8e5abea9291627a5edd859439c4553100c71309b86599f02414a532i0`
- tx: https://explorer.kaspa.org/txs/223233acf8e5abea9291627a5edd859439c4553100c71309b86599f02414a532
- owner: `kaspa:qzt9yuqceqvt2vk9dz7ddzayaa5flnenkymec59xvzm55ln3k72vgecxjhnjp`

That tx is **not** a Toccata covenant program. Your indexer has it. [covenants.kaspa.com](https://covenants.kaspa.com/covenants) does not.

### Covenant deploys (on covenants.kaspa.com)

Unrevealed P2SH, action `deploy`, identity `scriptHashFallback`.

| Name | KAS | Covenant |
| --- | ---: | --- |
| bakery.kas | 0.5 | https://covenants.kaspa.com/covenants/9baa5cb8f06876c1ba858c844e03d5d02da0419673761d4dea88e21c2432a685 |
| michelleobamaisaman.kas | 100 | https://covenants.kaspa.com/covenants/f7fbf09ab0c90eb89228528590fa08ca4eb4e1ce73d62db57fb85f5626c69fd6 |
| trump.kas | 100 | https://covenants.kaspa.com/covenants/5a13bcd0132897fbd946853544c91ce9c3a57b1c7dfb9ce69c54bb3e745f37b1 |
| opus.kas | 100 | https://covenants.kaspa.com/covenants/72c55a09aea7f15bfad99d1efc148597bc4eead8fd53d07627f187942d99ea51 |
| opus.dei.kas | 100 | https://covenants.kaspa.com/covenants/6f7b81bdf9269b7edf495065630b2ec9c79b6143f27fff0c8101988ca7602084 |

Re-check: `go run ./cmd/kns-spec prove`

## Pages on the URL

- https://stp-kas.github.io/kns-spec/ — this letter
- https://stp-kas.github.io/kns-spec/protocol.html — envelope, fees, indexer
- https://stp-kas.github.io/kns-spec/kasware.html — KasWare calls
- https://stp-kas.github.io/kns-spec/kastle.html — Kastle calls (fee output 0 still required)
- https://stp-kas.github.io/kns-spec/kasware-create.html — working inscribe page
- https://stp-kas.github.io/kns-spec/proofs.html — txs

## Resolve warning (your own integration note)

Before sending KAS to a resolved `.kas` name, show the address. The map is indexer-derived. Transfers cannot be reversed.

## What this is not

Not [app.knsdomains.org](https://app.knsdomains.org). Not a deployed global registrar. Not a seed prompt.
