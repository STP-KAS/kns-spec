# KasWare: inscribe a .kas name

Docs: [buildScript](https://docs.kasware.xyz/wallet/developer-documentation/kaspa/kaspa-krc20) (KRC-20 page also covers KNS). Inject: `window.kasware`. Chrome extension + Android APK. Not iOS.

Do not call `requestAccounts()` on page load. User click only.

## Connect

```js
const accounts = await window.kasware.requestAccounts();
const address = accounts[0];
```

Quiet resume: `getAccounts()`.

## Create `example.kas`

```js
const payload = JSON.stringify({ op: "create", p: "domain", v: "example" });

const { script, p2shAddress } = await window.kasware.buildScript({
  type: "KNS",
  data: payload,
});

const entries = await window.kasware.getUtxoEntries();
const [address] = await window.kasware.getAccounts();
const network = await window.kasware.getNetwork(); // kaspa_mainnet → "mainnet"

const fee = "kaspa:qyp4nvaq3pdq7609z09fvdgwtc9c7rg07fuw5zgeee7xpr085de59eseqfcmynn";
const priceKas = 35; // 5+ graphemes; see PROTOCOL.md

const { commitIds, revealIds } = await window.kasware.submitCommitReveal(
  {
    priorityEntries: [],
    entries,
    outputs: [{ address: p2shAddress, amount: 1 }],
    changeAddress: address,
    priorityFee: 0.01,
  },
  {
    outputs: [{ address: fee, amount: priceKas }],
    changeAddress: address,
    priorityFee: 0.02,
  },
  script,
  "mainnet"
);

// revealIds[0] is the inscription tx. Inscription id is usually `${revealIds[0]}i0`.
```

Amounts in `submitCommitReveal` are **KAS**, not sompi.

`buildScript` is safe (no popup). `submitCommitReveal` requires approval. It waits for commit confirmation, then reveals.

Split flow if you need a delay: `submitCommit` then `submitReveal`.

## Transfer

Same `type: "KNS"`, payload:

```json
{"op":"transfer","p":"domain","id":"<inscriptionId>","to":"kaspa:q…"}
```

Sender must be the current owner.

## Check first

```js
const res = await fetch("https://api.knsdomains.org/mainnet/api/v1/domains/check", {
  method: "POST",
  headers: { "content-type": "application/json" },
  body: JSON.stringify({ domainNames: ["example.kas"], address }),
});
```

If `available` is false, do not inscribe. The fee is still paid on a losing reveal.

## Working page

Hosted: [kasware-create.html](https://stp-kas.github.io/kns-spec/kasware-create.html). Local copy: [`examples/kasware-create.html`](examples/kasware-create.html). Never asks for a seed.

Kastle (two popups, namespace `kns`): [`KASTLE.md`](KASTLE.md).
