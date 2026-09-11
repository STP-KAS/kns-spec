# Profile keys

Official dashboard fields today (each a **separate text inscription**, 1 KAS, follow the domain on transfer):

| Key | Official name | Live on kns.kas |
| --- | --- | --- |
| `avatarUrl` | Avatar | yes |
| `website` | Website URL | `https://app.knsdomains.org/` |
| `banner` | Banner | |
| `x` | X handle | `knsdomain` |
| `github` | GitHub | |
| `telegram` | Telegram | `knsdomains` |
| `discord` | Discord | |
| `email` | Contact email | |
| `redirectUrl` | .limo redirect | `https://app.knsdomains.org/` |
| `bio` | Bio | |

API: `GET /api/v1/domain/{assetId}/profile?keys=redirectUrl,bio,avatarUrl,website,x,github,telegram,discord,email,banner`

## Add these (same mechanism — no new consensus)

| Key | Why |
| --- | --- |
| `ipfs` | CID for the dApp. `ipfs://bafy…` |
| `kfs` | Kaspa File Storage pointer |
| `contenthash` | ENS-style contenthash if you already store one |
| `peer` | libp2p multiaddr(s), comma-separated |
| `onion` | `onion://` or `i2p://` |
| `agent` | CID or URL of a callable agent card |
| `kas` | Pay address if it is **not** the inscription owner |
| `noise` | Overlay static key if not the Kaspa schnorr |

App resolve order: `redirectUrl` → `website` → `ipfs` → `kfs` → `arweave` → `contenthash`.

`website` / `redirectUrl` stay the **old web**. `ipfs` + `peer` are the overlay. Both can exist on the same name.

Schema: [`schemas/overlay-records.schema.json`](schemas/overlay-records.schema.json).
