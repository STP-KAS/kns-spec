# Supporting wallets (official)

Source: [KNS supporting wallet](https://kns-2.gitbook.io/kns-docs-1/supporting-wallet). Tutorials under that page.

| Wallet | Type | Inscribe | Transfer `.kas` | Receive `.kas` |
| --- | --- | --- | --- | --- |
| **KasWare** | Extension | yes | KAS, KRC-20, KRC-721, domain, text | same |
| **Kastle** | Extension | yes (two popups) | KAS, KRC-20, KRC-721, domain | same |
| **Kastle** | Mobile | **no** | KAS, KRC-20, KRC-721 | + `.kas` domain |
| **Kurncy** | Mobile | yes | KAS, KRC-20, KRC-721, domain | same |
| **Kasanova** | Mobile | yes | KAS, KRC-20, KRC-721, domain | same |

- KasWare inject: `window.kasware` — [KASWARE.md](KASWARE.md). Tutorial: [kasware-wallet-tutorial](https://kns-2.gitbook.io/kns-docs-1/supporting-wallet/kasware-wallet-tutorial.md).
- Kastle inject: `window.kastle` — [KASTLE.md](KASTLE.md). Extension only for inscribe. Tutorial: [kastle-wallet-tutorial](https://kns-2.gitbook.io/kns-docs-1/supporting-wallet/kastle-wallet-tutorial.md).
- Kurncy: in-app. [@KurncySolutions](https://x.com/KurncySolutions). Tutorial: [kurncy-wallet-tutorial](https://kns-2.gitbook.io/kns-docs-1/supporting-wallet/kurncy-wallet-tutorial.md).
- Kasanova: in-app. [@KasanovaWallet](https://x.com/KasanovaWallet). Tutorial: [kasanova-wallet-mobile-tutorial](https://kns-2.gitbook.io/kns-docs-1/supporting-wallet/kasanova-wallet-mobile-tutorial.md).

Official must-reads:

- If the **receiving** wallet is not a supporting wallet, the domain is **not visible** and cannot be managed or transferred out.
- **KNS does not support ECDSA addresses.**
- Always show the resolved `kaspa:` address before send. Transfers cannot be reversed.

FAQ still says “only KasWare” in places. Use **this table**, not the FAQ, for wallet support.

dApp inject on the open web is still only KasWare and Kastle extension. Kurncy and Kasanova inscribe inside their apps. [KCC-0012](https://github.com/kaspanet/kccs/pull/24) (wallet provider discovery) is **Draft**. Do not wait for it.
