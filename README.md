# KNS implementer kit

**URL for the KNS team:** https://stp-kas.github.io/kns-spec/

Letter in this repo: [FOR-KNS.md](FOR-KNS.md).

For the KNS team and wallets. Two objects. Do not mix them.

| Layer | On chain today | Uniqueness | Wallet |
| --- | --- | --- | --- |
| **Inscription** | `kns` commit-reveal envelope | Official indexer, first valid reveal | **KasWare** `buildScript({ type: "KNS" })` |
| **Covenant** | P2SH of `KasName.sil` | **Not** consensus. Script-hash lock only | Fund / later spend the P2SH |

Live uniqueness of `alice.kas` is still the [KNS indexer](https://api.knsdomains.org/mainnet). A KIP-20 `covenant_id` is hashed from an outpoint. It does not encode the label. Anyone can genesis another UTXO that writes `alice` in state. Nodes accept both.

Sister demo (Web4 UI, not this spec): [STP-KAS/kns](https://github.com/STP-KAS/kns).

## Implement this

Official docs: [inscriptions](https://kns-2.gitbook.io/kns-docs-1/inscriptions/overview) · [wallets](https://kns-2.gitbook.io/kns-docs-1/supporting-wallet) · [indexer API](https://kns-2.gitbook.io/kns-docs-1/kns-indexer-api) · [simply-kaspa-indexer](https://github.com/supertypo/simply-kaspa-indexer)

**Keep it real:** [`REAL.md`](REAL.md) — what is live vs paper, and the non-crypto work (wallet chrome, indexer spec, pinning, no lookup logs). Architecture map (not shipped): [`WEB4.md`](WEB4.md).

1. **Inscribe** with a [supporting wallet](WALLETS.md). KasWare ([`KASWARE.md`](KASWARE.md)) or Kastle **extension** ([`KASTLE.md`](KASTLE.md)). Kastle **mobile cannot inscribe**. Envelope: [`PROTOCOL.md`](PROTOCOL.md). Check the indexer first ([`INDEXER.md`](INDEXER.md)).
2. **Resolve** with `api.knsdomains.org`. URL-encode names. Warn before sending KAS to a resolved address.
3. **Optional elevate** to a Name UTXO: compile [`contracts/v1/KasName.sil`](contracts/v1/KasName.sil) with official **silverc v1.0.0** (Ori, 9 Sep 2026). Own-UTXO only. Continuation keeps the same sompi (fees from a sibling input). No `readInputState` of a foreign covenant ([silverscript#234](https://github.com/kaspanet/silverscript/pull/234) still unmerged). Battle-test notes: [`BATTLETEST.md`](BATTLETEST.md).

```powershell
go test ./...
go run ./cmd/kns-spec prove
go run ./cmd/kns-spec check kns.kas
go run ./cmd/kns-spec plan example
go run ./cmd/kns-spec resolve kns.kas
go run ./cmd/kns-spec overlay kns.kas
go run ./cmd/kns-spec vectors
go run ./cmd/kns-spec bind alice.kas <xonly> <noise-x25519>
```

## Proven mainnet txs

Checked 2026-09-08 against `api.kaspa.org`, `api.knsdomains.org`, and `indexer.kaspa.com` (the API behind [covenants.kaspa.com](https://covenants.kaspa.com/covenants)).

### 1. KNS inscription (KasWare-compatible envelope)

`kns.kas` reveal spends a P2SH whose redeem script is:

```
<xonly_pubkey> OP_CHECKSIG OP_FALSE OP_IF <kns> <0> <{"op":"create","p":"domain","v":"kns"}> OP_ENDIF
```

| | |
| --- | --- |
| name | `kns.kas` |
| inscription id | `223233acf8e5abea9291627a5edd859439c4553100c71309b86599f02414a532i0` |
| reveal tx | [explorer](https://explorer.kaspa.org/txs/223233acf8e5abea9291627a5edd859439c4553100c71309b86599f02414a532) |
| owner | `kaspa:qzt9yuqceqvt2vk9dz7ddzayaa5flnenkymec59xvzm55ln3k72vgecxjhnjp` |
| created | 2025-01-27 |

That reveal is a **pubkey spend with a `kns` envelope**, not a Toccata covenant program. The covenant explorer has no row for this tx. The indexer does.

### 2. Covenant deploys on covenants.kaspa.com

These are **unrevealed P2SH deploys**. Explorer action = `deploy`, status = `unrevealed`, identity = `scriptHashFallback` (no KIP-20 `covenant_id` on the UTXO yet). First output is the name lock.

| Name | KAS | Tx | Covenant |
| --- | ---: | --- | --- |
| bakery.kas | 0.5 | [095622d4…](https://covenants.kaspa.com/tx/095622d4aa4f46cac8a4012d289ed78a7bd16e9acd47d7e94ee17ae0d601f1ca) | [9baa5cb8…](https://covenants.kaspa.com/covenants/9baa5cb8f06876c1ba858c844e03d5d02da0419673761d4dea88e21c2432a685) |
| michelleobamaisaman.kas | 100 | [abf7d9e1…](https://covenants.kaspa.com/tx/abf7d9e1199e7d3c2fd7121b261aa3284eb062a19e5c2f1afefc0257f2d61c2d) | [f7fbf09a…](https://covenants.kaspa.com/covenants/f7fbf09ab0c90eb89228528590fa08ca4eb4e1ce73d62db57fb85f5626c69fd6) |
| trump.kas | 100 | [22ca97e1…](https://covenants.kaspa.com/tx/22ca97e171336859031fdc12c649e50711c330a43a7c8350f7038ab75abb1084) | [5a13bcd0…](https://covenants.kaspa.com/covenants/5a13bcd0132897fbd946853544c91ce9c3a57b1c7dfb9ce69c54bb3e745f37b1) |
| opus.kas | 100 | [bb8e054d…](https://covenants.kaspa.com/tx/bb8e054d8cf3f9f75dc99d76910669b48c8cff36d9d1fe1b378a007c653fadf3) | [72c55a09…](https://covenants.kaspa.com/covenants/72c55a09aea7f15bfad99d1efc148597bc4eead8fd53d07627f187942d99ea51) |
| opus.dei.kas | 100 | [34d410ef…](https://covenants.kaspa.com/tx/34d410ef5df5d75256334f6521c62f1654c624987e5ac0c049ed68c54061e7b7) | [6f7b81bd…](https://covenants.kaspa.com/covenants/6f7b81bdf9269b7edf495065630b2ec9c79b6143f27fff0c8101988ca7602084) |

`trump.kas` and `bakery.kas` already exist as **KNS inscriptions owned by other keys**. The P2SH above is a second object. `opus.kas` has no inscription (indexer: domain not found). Subnames (`opus.dei.kas`) are not an inscription op.

Machine-readable: [`proofs/proofs.json`](proofs/proofs.json). Re-check: `go run ./cmd/kns-spec prove`.

## What this repo is not

- Not the official KNS product ([app.knsdomains.org](https://app.knsdomains.org), [@knsdomain](https://x.com/knsdomain)).
- Not a deployed global registrar. No consensus uniqueness for `.kas`.
- Not a seed prompt. Never ask for a seed.

## License

MIT. Inscription protocol belongs to KNS. `KasName.sil` is compiled with official **silverc v1.0.0** (`3ed9733`, Ori / someone235). Template hash `c8c06c1abe007e97f78b3f1701a443a41f54c65b878113c0d6cf3ed4b47fa79b` (value-conservation + transfer clears pay/vault). Windows zip SHA256 `3e0d660c15a9e7ac90f3960da24d348b076b1891481bfe758db18accc8a102e1`.
