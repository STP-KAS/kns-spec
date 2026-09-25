# KNS snapshot test holder set (Testnet-10)

Test addresses and the KNS names they hold, for the KNS team's snapshot / claim testing. **Kaspa Testnet-10 only.** tKAS is worthless testnet coin. Not mainnet. Not Kaspa core. Not official KNS. Experimental. Not advice.

Test log, scripts and per-name tx ids: [STP-KAS/kns-tn10-testing](https://github.com/STP-KAS/kns-tn10-testing).

## Files

| File | Rows | Columns |
| --- | --- | --- |
| [addresses.csv](addresses.csv) | 25,700 + header | `index,address,derivation_path` |
| [holders-with-names.csv](holders-with-names.csv) | 25,700 + header | `index,address,name_count,names` |
| [TESTING-SNAPSHOT.md](TESTING-SNAPSHOT.md) | — | how we recommend testing a snapshot with this set |

- `addresses.csv`: all snapshot test addresses, indices 0–25699. `kaspatest:` addresses only. Derivation template BIP44 `m/44'/111111'/0'/0/{i}` (one testnet mnemonic, account 0). Copied unchanged from the test box and checked row by row: index, `kaspatest:` address, matching path, nothing else. 25,700 unique addresses.
- `holders-with-names.csv`: the same addresses, each with the `.kas` names it created on TN10. `names` is semicolon-separated. Addresses with no names have `name_count` 0.

## How names were joined

- Source: the box's per-create result records (`api/smoke-result-<label>.json`, one per commit + reveal written by the create scripts). The public fields of the same records are exported to [kns-tn10-testing/results/inscriptions.csv](https://github.com/STP-KAS/kns-tn10-testing/blob/main/results/inscriptions.csv). Only records with a reveal tx id count. Deduplicated by name.
- Join key: the paying address of the create (`payer`) = the snapshot address. A create's payer is the name's owner at inscription time. No transfers were made.
- Check: for every `stp-snap-w###-d###` and `stp-s3-w#####-d#` name, the wallet number in the name equals the address `index`. 0 mismatches. All R1 names sit on indices 20700–25699, as planned.
- As of **25 Sep 2026 06:01 CEST** (Europe/Brussels; newest create in the set 06:00 CEST):

| | Addresses | Names |
| --- | ---: | ---: |
| Phase A `stp-snap-*` (up to 100 per address) | 357 | 28,837 |
| S3 `stp-s3-*` (1–3 per address) | 5,426 | 13,556 |
| R1 random names, 3–10 chars (1 per address; 119 three-char, 122 four-char) | 910 | 910 |
| **Total with at least one name** | **6,693** | **43,303** |
| 0 names | 19,007 | 0 |

- 248 addresses hold exactly 100 names; 965 hold 1; 2,614 hold 2; 2,760 hold 3.
- Runs are still going (one slow runner, at most 2 creates per 10 min since ~02:00 CEST 25 Sep). Counts are a minimum.
- Ownership was not re-queried for every name. A read-only sample of 45 names (15 Phase A, 15 S3, 15 R1) against `GET https://api.knsdomains.org/tn10/api/v1/<name>/owner` at ~06:05 CEST: 45/45 owner = the listed address.
- Bulk `stp-bulk-*` (1,341) and smoke (1) names were paid from a separate address that is not in this set, so they are not listed here. It is a useful non-snapshot control.
- SHA-256 of the files as published: `addresses.csv` `d92812fe10fd82f75f21c68841abcb5a2d1ddaa857c0da0bacc501079c074357`, `holders-with-names.csv` `6bdb4f8cb1c2cdf46f751a055a6072aa861dba7d8adf3f8276d36808ba1b7c4d`.
- Previous version (24 Sep 23:20 CEST: 1,048 addresses, 7,614 names) is in git history at commit `0a4d07f`. Every name in it is still in this file.

## Testing with this set

Our recommendation for how a snapshot works best for testing, and what we can do on our side: [TESTING-SNAPSHOT.md](TESTING-SNAPSHOT.md).

## Safety

Keys and the mnemonic are never published, here or in kns-tn10-testing. Nobody from this repo will DM you or ask for a seed, key or password.
