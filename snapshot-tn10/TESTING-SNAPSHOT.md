# How a snapshot works best for testing: our recommendation

From STP-KAS to the KNS team, for the airdrop / claim work on Testnet-10. **TN10 only.** tKAS is worthless. Not official KNS. Experimental. Not advice. Take what is useful.

Test set: [addresses.csv](addresses.csv) and [holders-with-names.csv](holders-with-names.csv) (counts and as-of time in [README.md](README.md)). Test log: [STP-KAS/kns-tn10-testing](https://github.com/STP-KAS/kns-tn10-testing).

## 1. Fix one freeze point and publish it first

- Pick an explicit **block (DAA score) and time**. Announce it before it happens.
- Snapshot = **indexer state at that point**: one row per name, with its owner address. Not wallet balances, not raw chain scans.
- Export it as a flat file (CSV or JSON) sorted by name, and publish its **SHA-256**. Anyone can then re-download, re-hash and diff. We do this for our own files (hashes in [README.md](README.md)).
- Say which indexer build produced it, so a re-run gives the same file.

## 2. What goes in, what stays out

- **In:** every revealed name with its current owner. That includes addresses with many names (we have 248 with 100 each), short names (3- and 4-char), random-looking names, and names on addresses with 0 or dust balance. Holding a name is the entitlement, not holding KAS.
- **Out:** commits without an accepted reveal. We saw why this matters: a create interrupted between commit and reveal leaves 1 tKAS locked in the P2SH and **no name**. Such stranded commits exist on chain; they must not count.
- **Indexer rules decide.** The TN10 indexer accepts labels that the kns-spec regex rejects (leading or trailing hyphens). Whatever the indexer accepted before the freeze is in the snapshot. Write that down, so nobody argues later.
- **Transfers near the freeze:** state one rule, e.g. "owner as of the last indexed block at or before the freeze score". A transfer revealed after it counts for the new owner only in a later snapshot.

## 3. What this test set gives you

| Category | Where | What a correct claim does |
| --- | --- | --- |
| Many names on one address | `stp-snap-*`, indices 0–699 (up to 100 each) | Entitlement = all names, counted once each |
| 1–3 names per address | `stp-s3-*`, indices 700–20699 | Entitlement = exactly those names |
| One short or random name | R1, indices 20700–25699 (3–10 chars) | Entitlement = that one name; check fee-tier names (3, 4 chars) are not treated differently unless you intend it |
| **0 names** (negative control) | 19,000+ addresses in the same derivation, same wallet | Nothing. Any claim must fail |
| Non-snapshot payer (control) | the bulk/smoke address (1,342 names, not in `addresses.csv`) | Gets its own names, and nothing from the snapshot set |

All addresses come from one BIP44 path (`m/44'/111111'/0'/0/{i}`), so they look alike. A claim tool cannot pass by accident on "different-looking" wallets.

## 4. How to test, in order

1. **Dry-run export and diff.** Export your snapshot for our addresses and diff it against `holders-with-names.csv`. Tell us the freeze score and we re-export ours at that point, so both sides compare the same moment. Any difference is either a bug or a rule to write down.
2. **Claims on a sample, every category.** Pick addresses from each row above. We sign the claims with our keys, on our side. You never need a key.
3. **Edge cases.** Double claim (second one must fail), claim from the wrong address, claim with a name that was transferred after the freeze, claim from a 0-name address.
4. **Scale.** Thousands of claims from thousands of addresses. Our create runs show what to expect: just before we slowed down (25 Sep 02:18 CEST) our runners had made 487 creates in 5 minutes; the public APIs answered with occasional HTTP 429; and a public node that lags makes retries fail with "already spent". Watch the indexer, the API and the claim backend under that kind of load.
5. **Freeze-point behaviour.** Names revealed in the last blocks before the freeze, blocks that arrive late, and reorgs. In one early bulk batch, 49 of 100 creates ended with no result record and an indexer 404 and had to be backfilled; our fixed create code now waits until each reveal is accepted. Do the same at the freeze: wait for acceptance and reconcile, don't assume.
6. **Publish and reconcile.** Totals per category (names, addresses, claimed, rejected) on both sides. They must add up.

## 5. What we can do for you

- Re-export `holders-with-names.csv` on request, at any block or time you name.
- Run claims from any subset of the 25,700 addresses, any size, on your schedule.
- **Stop our create runners before your freeze point.** Tell us the time; we stop new creates before it and confirm nothing is left half-done (no commit without reveal). Right now only one slow runner is active, at most 2 creates per 10 minutes.
- Share the CSVs and per-name tx ids ([results/inscriptions.csv](https://github.com/STP-KAS/kns-tn10-testing/blob/main/results/inscriptions.csv)).

## Keys

The mnemonic and private keys stay with STP-KAS. They are never published, here or anywhere, and never sent to you. We sign; you verify. Nobody from this repo will DM you or ask for a seed, key or password.
