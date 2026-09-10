# Kastle: inscribe a .kas name

Official KNS [supporting wallet](https://kns-2.gitbook.io/kns-docs-1/supporting-wallet): **extension can inscribe**; **mobile cannot**. Inject: `window.kastle`. Docs: [Kastle Wallet API](https://docs.kastle.cc/readme/how-to-integrate/kastle-wallet-api.md). Official KNS tutorial: two popups (commit, then reveal). KNS does not support ECDSA addresses.

Do not call `connect()` on page load. User click only.

## Connect

```js
const ok = await window.kastle.connect();
const { address, publicKey } = await window.kastle.getAccount();
const network = await window.kastle.getNetwork(); // "mainnet"
```

## Create `example.kas`

`commitReveal` puts the protocol id in the envelope as `namespace`. For KNS that id is `kns` (same bytes KasWare uses for `type: "KNS"`). For KRC-20 Kastle uses `"kasplex"`.

```js
const payload = JSON.stringify({ op: "create", p: "domain", v: "example" });

const { commitTxId, revealTxId } = await window.kastle.commitReveal(
  "mainnet",
  "kns",
  payload
);
// revealTxId is the inscription tx. Id is usually `${revealTxId}i0`.
```

Kastle opens **two** confirmations. That is expected.

## Fee (do not skip)

KNS requires reveal **output 0** to pay:

`kaspa:qyp4nvaq3pdq7609z09fvdgwtc9c7rg07fuw5zgeee7xpr085de59eseqfcmynn`

`commitReveal` options are only priority fees (`commitPriorityFee`, `revealPriorityFee`, sompi strings). It does **not** document a custom fee output. The official KNS inscribe tool still has to attach that output. If you only call `commitReveal`, check the reveal tx: if output 0 is not the protocol address, the indexer will not verify the domain.

Until Kastle documents a kns fee output, use the official inscribe UI or build the reveal with `buildTransaction` / `signAndBroadcastTx` so output 0 is the fee.

## Transfer

Same namespace `kns`:

```json
{"op":"transfer","p":"domain","id":"<inscriptionId>","to":"kaspa:q…"}
```

## Check first

Same as KasWare: `POST /api/v1/domains/check`. If `available` is false, do not inscribe.
