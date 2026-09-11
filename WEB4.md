# Web4 `.kas`: the name *is* the computer

Four independent reviews (ENS protocol, privacy overlay, Kaspa runtime, other name systems) collapse to one design.

**KNS already copied ENS’s skin** (normalize, pay-to-name, text profile, primary name, `.limo`).  
**It did not copy ENS’s machinery** (consensus registry, pluggable resolver, hierarchy, contenthash, verifiable resolve).

The extra internet layer is not an L2 EVM and not Tor with a Kaspa sticker. It is four planes:

```
Run      →  local sandbox (HTML/WASM from a CID)
Session  →  Noise + capability. Not TLS-to-Cloudflare.
Locate   →  records on the name (reproducible FCFS)
Settle   →  Kaspa PoW UTXO (inscription today, optional KasName.sil)
```

`https://alice.kas.limo` is the old web. It needs DNS and a CA. `kns://alice.kas` must work with **your** node.

## Do not copy from ENS

| Copying this | Why it is a downgrade |
| --- | --- |
| Rent / expiry | KNS already owns forever. Rent taxes identity. |
| One registry contract | Serial bottleneck + politics. Keep names as data on L1, optionally a **UTXO per name**. |
| Infura as required resolve | `api.knsdomains.org` is already that. Publish FCFS rules so simply-kaspa-indexer can reindex. |
| NameWrapper fuses | No expiry → a burnt fuse is forever. ENSv2 walked this back. |
| DNSSEC as the root | Binds `.kas` to ICANN. Optional *bridge out* only. |
| L2 as where names become real | Makes L1 names souvenirs. Kaspa L1 is the cheap parallel layer. |

## Five product changes (no hard fork)

These are the only things that make a `.kas` name a **dApp address**.

1. **Index overlay keys like `website`.**  
   `ipfs` `kfs` `contenthash` `peer` `onion` `agent` `kas` `noise` — 1 KAS text inscriptions, Profile API `?keys=`. See [PROFILE.md](PROFILE.md).

2. **Publish the indexer as a spec, not a host.**  
   Exact FCFS rules so anyone with [simply-kaspa-indexer](https://github.com/supertypo/simply-kaspa-indexer) can resolve. Official API becomes a convenience. Until then uniqueness is “trust knsdomains.org.”

3. **Reverse resolution with a forward check.**  
   Primary Name API already exists. Display `alice.kas` for an address **only if** Domain API owner == that address. FAQ saying this waits for L2 is a stall.

4. **Wallet-native `kns://`, not only `.limo`.**  
   Open app bytes on the device. `/pay` shows `kaspa:`. `/peer` is session. KasWare / Kastle / Kurncy / Kasanova are the distribution.

5. **Split owner, pay, and session.**  
   ENS: registry owner ≠ `addr` ≠ resolver. KNS today: owner **is** pay. Add `kas` (pay if not owner) and `noise`/`peer` (session if not the spend key). Keep ECDSA out of the happy path.

## Session plane (privacy without folklore)

Kaspa spend is **BIP340 Schnorr**. libp2p PeerId / Noise `identity_sig` is **ECDSA**. Onion v3 is **Ed25519**. Noise DH is **X25519**.

**Do not claim “libp2p identity = Kaspa schnorr.”** The name key **binds** overlay keys:

```
BINDING = "kns-session/v1" || name || owner_xonly || noise_x25519 || [peerid] || [onion] || seq || exp
sig     = KIP-5 Schnorr(owner, BINDING)
```

Derive overlay keys with HKDF from the seed (`kns-overlay-v1` / name). **Never** use the spend key as a Noise static. An online session must not prove the UTXO key is hot except for a one-shot BINDING signature.

| On L1 (eternal, correlatable) | Off-chain MUST |
| --- | --- |
| Owner, transfers, CID, optional BINDING public half | Seed, Noise secret, IPs, ports, tickets, transcripts, vault plaintext, who resolved whom |

Default rendezvous is a **capability** (`kns://alice.kas/cap/…`), not a public `/ip4/` multiaddr. Public `peer` is a shop window. Say so in the UI.

Fake privacy (refuse in wallet copy): “libp2p is anonymous”, “onion on-chain is private”, “402 is private pay”, “IPFS is private”, “Noise without BINDING”.

Honest 402: challenge **over the session**, pay **on L1** to a **fresh** address, verify with **your** node, then mint a macaroon. L1 payment is public. Do not put session ids in tx payload.

## How a dApp runs (v0)

Not Argent (not release-ready). Not vProgs. Not an EVM.

1. Resolve name (local indexer first).
2. Show `kaspa:` address. Official warning.
3. Fetch CID via embedded bitswap/KFS or the Noise session — **not** `ipfs.io`.
4. Run in a wallet webview / WASM sandbox. Origin = `kns://alice.kas`.
5. Calls authenticated by BINDING. Pay with 402 as above.

Covenants authorize **spends**. They are not a world computer. KIP-21 lanes can wait.

## 90 days

| Days | Ship |
| --- | --- |
| 1–15 | `kns://` parser, local resolve, pay UI, no IP in schema |
| 10–30 | BINDING + HKDF overlay keys, Schnorr-only |
| 20–50 | Noise XX in the wallet process |
| 35–60 | Capabilities default; public locator opt-in |
| 45–70 | CID fetch without DNS gateways; local sandbox |
| 55–90 | 402 over Noise; vault ciphertext never on L1 |

Explicit non-goals in 90 days: consensus uniqueness, Argent, mixnet payments, one key for Kaspa+libp2p+onion.

## Other systems: steal / reject

| Steal | Reject |
| --- | --- |
| Unstoppable: payee ≠ owner + IPFS with HTTP fallback | Company registrar / hosted pin as the protocol |
| Handshake: name is a wallet object; proofs later | Browser forks / TLD auctions year one |
| SNS: primary name in the send box; records go stale on transfer | One vendor SDK as the spec |
| ENS CCIP: cheap offchain *profile* edits | Company-rewritable names as the default |
| IPNS: mutable pointer *inside* the session | DHT as the wallet resolve path |
| Farcaster: forever owner, swappable display name | Central reclaim of handles |
| Nostr: `/.well-known/kas.json` as a **hint** | HTTPS file as ownership |

## KNS dashboard this week

Same 1 KAS text fields as website:

`ipfs` · `kfs` · `contenthash` · `peer` · `onion` · `agent` · `kas` · `noise`

Also this week, no fork:

- **`kas` / `addr`:** Send pays this, not the inscription owner (stop doxxing the vault that holds the name).
- **Primary name in KasWare chrome** + show `kaspa:` before sign.
- **One-shot `/resolve?q=`** for wallets.
- **Batch profile save** (N signatures, one sheet) — 520-byte envelope still applies per field.
- **`/.well-known/kas.json`** as a hint ([example](examples/well-known-kas.json)). Identity is the inscription.
- **Clear profile on transfer.**
- Hide `email` behind “this is public forever.”

Plus: Profile API returns **all** keys, not a hardcoded allowlist. Primary name requires forward check. Document FCFS so a third indexer can exist.

`kns://` execution uses **Run** order (`ipfs` → `kfs` → `contenthash`). Never execute `website`.

Working companion (no seed): https://stp-kas.github.io/kns-spec/open.html — resolve any `.kas`, split web vs run, show pay address.

Checklist: [CONFORMANCE.md](CONFORMANCE.md). Overlay: [OVERLAY.md](OVERLAY.md).
