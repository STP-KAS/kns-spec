# KNS snapshot test holder set (Testnet-10)

Test addresses and the KNS names they hold, for the KNS team's snapshot / claim testing. **Kaspa Testnet-10 only.** tKAS is worthless testnet coin. Not mainnet. Not Kaspa core. Not official KNS. Experimental. Not advice.

Test log, scripts and per-name tx ids: [STP-KAS/kns-tn10-testing](https://github.com/STP-KAS/kns-tn10-testing).

## Files

| File | Rows | Columns |
| --- | --- | --- |
| [addresses.csv](addresses.csv) | 25,700 + header | `index,address,derivation_path` |
| [holders-with-names.csv](holders-with-names.csv) | 25,700 + header | `index,address,name_count,names` |

- `addresses.csv`: all snapshot test addresses, indices 0–25699. `kaspatest:` addresses only. Derivation template BIP44 `m/44'/111111'/0'/0/{i}` (one testnet mnemonic, account 0). Copied unchanged from the test box and checked row by row: index, `kaspatest:` address, matching path, nothing else. 25,700 unique addresses.
- `holders-with-names.csv`: the same addresses, each with the `.kas` names it created on TN10. `names` is semicolon-separated. Addresses with no names have `name_count` 0.

## How names were joined

- Source: the create results in [kns-tn10-testing/results/inscriptions.csv](https://github.com/STP-KAS/kns-tn10-testing/blob/main/results/inscriptions.csv), one row per commit + reveal recorded by the create scripts.
- Join key: the paying address of the create (`payer_address`) = the snapshot address. A create's payer is the name's owner at inscription time. No transfers were made.
- As of **24 Sep 2026 23:20 CEST** (Europe/Brussels): **1,048** addresses hold at least one name; **7,614** names in total (Phase A `stp-snap-*` 6,207, S3 `stp-s3-*` 791, R1 random names 616). The other 24,652 addresses have 0 so far.
- Runs were still going at that time. Counts are a floor. Indexer ownership was not re-queried for this file; check `GET https://api.knsdomains.org/tn10/api/v1/<name>/owner`.
- Bulk `stp-bulk-*` and smoke names were paid from a separate address that is not in this set, so they are not listed here.

## Safety

Keys and the mnemonic are never published, here or in kns-tn10-testing. Nobody from this repo will DM you or ask for a seed, key or password.
