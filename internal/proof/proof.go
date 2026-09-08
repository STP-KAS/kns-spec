package proof

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/STP-KAS/kns-spec/internal/kns"
)

type File struct {
	Checked      string        `json:"checked"`
	Inscriptions []Inscription `json:"inscriptions"`
	Covenants    []Covenant    `json:"covenants"`
}

type Inscription struct {
	Name          string `json:"name"`
	Payload       string `json:"payload"`
	InscriptionID string `json:"inscriptionId"`
	TxID          string `json:"txid"`
	Owner         string `json:"owner"`
	Explorer      string `json:"explorer"`
}

type Covenant struct {
	Name                 string `json:"name"`
	TxID                 string `json:"txid"`
	P2SH                 string `json:"p2sh"`
	ScriptHash           string `json:"scriptHash"`
	CovenantIDHex        string `json:"covenantIdHex"`
	ClassificationStatus string `json:"classificationStatus"`
	Action               string `json:"action"`
	AmountSompi          int64  `json:"amountSompi"`
	CovenantURL          string `json:"covenantUrl"`
	TxURL                string `json:"txUrl"`
}

type Result struct {
	Name   string `json:"name"`
	Kind   string `json:"kind"`
	OK     bool   `json:"ok"`
	Detail string `json:"detail"`
	URL    string `json:"url,omitempty"`
}

func Load(path string) (*File, error) {
	if path == "" {
		path = findProofs()
	}
	if path == "" {
		return nil, fmt.Errorf("proofs/proofs.json not found")
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var f File
	if err := json.Unmarshal(raw, &f); err != nil {
		return nil, err
	}
	return &f, nil
}

func findProofs() string {
	var starts []string
	if wd, err := os.Getwd(); err == nil {
		starts = append(starts, wd)
	}
	if exe, err := os.Executable(); err == nil {
		starts = append(starts, filepath.Dir(exe))
	}
	for _, start := range starts {
		dir := start
		for i := 0; i < 8; i++ {
			p := filepath.Join(dir, "proofs", "proofs.json")
			if _, err := os.Stat(p); err == nil {
				return p
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}
	return ""
}

func Verify(f *File) []Result {
	httpClient := &http.Client{Timeout: 18 * time.Second}
	knsClient := kns.New("")
	var out []Result
	for _, ins := range f.Inscriptions {
		out = append(out, verifyInscription(knsClient, httpClient, ins))
	}
	for _, c := range f.Covenants {
		out = append(out, verifyCovenant(httpClient, c))
	}
	return out
}

func verifyInscription(c *kns.Client, httpClient *http.Client, ins Inscription) Result {
	r := Result{Name: ins.Name, Kind: "inscription", URL: ins.Explorer}
	own, err := c.Owner(ins.Name)
	if err != nil {
		r.Detail = "kns indexer: " + err.Error()
		return r
	}
	if !strings.EqualFold(own.Owner, ins.Owner) {
		r.Detail = "owner mismatch " + own.Owner
		return r
	}
	if own.AssetID != ins.InscriptionID {
		r.Detail = "inscription id mismatch " + own.AssetID
		return r
	}
	ok, detail := l1Tx(httpClient, ins.TxID, ins.Payload)
	r.OK = ok
	r.Detail = detail
	return r
}

func verifyCovenant(httpClient *http.Client, c Covenant) Result {
	r := Result{Name: c.Name, Kind: "covenant", URL: c.CovenantURL}
	ok, detail := l1P2SH(httpClient, c.TxID, c.P2SH, c.AmountSompi)
	if !ok {
		r.Detail = detail
		return r
	}
	req, err := http.NewRequest(http.MethodGet, "https://indexer.kaspa.com/covenants/"+c.CovenantIDHex, nil)
	if err != nil {
		r.Detail = err.Error()
		return r
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "kns-spec/1.0")
	res, err := httpClient.Do(req)
	if err != nil {
		r.Detail = "covenant indexer: " + err.Error()
		return r
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode >= 400 {
		r.Detail = fmt.Sprintf("covenant indexer %s", res.Status)
		return r
	}
	var row struct {
		Covenant struct {
			CovenantIDHex        string `json:"covenantIdHex"`
			ClassificationStatus string `json:"classificationStatus"`
			GenesisTxIDHex       string `json:"genesisTxidHex"`
			Address              string `json:"address"`
		} `json:"covenant"`
	}
	if err := json.Unmarshal(body, &row); err != nil {
		r.Detail = "covenant json: " + err.Error()
		return r
	}
	if !strings.EqualFold(row.Covenant.CovenantIDHex, c.CovenantIDHex) {
		r.Detail = "covenant id mismatch"
		return r
	}
	if !strings.EqualFold(row.Covenant.GenesisTxIDHex, c.TxID) {
		r.Detail = "genesis tx mismatch " + row.Covenant.GenesisTxIDHex
		return r
	}
	if row.Covenant.Address != c.P2SH {
		r.Detail = "p2sh mismatch " + row.Covenant.Address
		return r
	}
	r.OK = true
	r.Detail = fmt.Sprintf("covenants.kaspa.com %s %s · L1 %s", row.Covenant.ClassificationStatus, c.Action, detail)
	return r
}

func l1Tx(httpClient *http.Client, txid, payload string) (bool, string) {
	raw, err := getJSON(httpClient, "https://api.kaspa.org/transactions/"+txid)
	if err != nil {
		return false, err.Error()
	}
	var tx struct {
		TransactionID string `json:"transaction_id"`
		IsAccepted    bool   `json:"is_accepted"`
		Inputs        []struct {
			SignatureScript string `json:"signature_script"`
		} `json:"inputs"`
	}
	if err := json.Unmarshal(raw, &tx); err != nil {
		return false, err.Error()
	}
	if !strings.EqualFold(tx.TransactionID, txid) || !tx.IsAccepted {
		return false, "tx not accepted"
	}
	if payload != "" {
		hexPayload := fmt.Sprintf("%x", payload)
		found := false
		for _, in := range tx.Inputs {
			if strings.Contains(strings.ToLower(in.SignatureScript), hexPayload) {
				found = true
				break
			}
		}
		if !found {
			return false, "envelope payload not in signature script"
		}
	}
	return true, "accepted on L1, envelope in reveal"
}

func l1P2SH(httpClient *http.Client, txid, p2sh string, sompi int64) (bool, string) {
	raw, err := getJSON(httpClient, "https://api.kaspa.org/transactions/"+txid)
	if err != nil {
		return false, err.Error()
	}
	var tx struct {
		TransactionID string `json:"transaction_id"`
		IsAccepted    bool   `json:"is_accepted"`
		Outputs       []struct {
			Index                  int    `json:"index"`
			Amount                 int64  `json:"amount"`
			ScriptPublicKeyAddress string `json:"script_public_key_address"`
			ScriptPublicKeyType    string `json:"script_public_key_type"`
		} `json:"outputs"`
	}
	if err := json.Unmarshal(raw, &tx); err != nil {
		return false, err.Error()
	}
	if !strings.EqualFold(tx.TransactionID, txid) || !tx.IsAccepted {
		return false, "tx not accepted"
	}
	if len(tx.Outputs) == 0 || tx.Outputs[0].ScriptPublicKeyAddress != p2sh {
		return false, "output 0 is not the name P2SH"
	}
	if tx.Outputs[0].ScriptPublicKeyType != "scripthash" {
		return false, "output 0 is not scripthash"
	}
	if sompi > 0 && tx.Outputs[0].Amount != sompi {
		return false, fmt.Sprintf("amount %d want %d", tx.Outputs[0].Amount, sompi)
	}
	return true, fmt.Sprintf("accepted · %s · %d sompi", p2sh, tx.Outputs[0].Amount)
}

func getJSON(httpClient *http.Client, u string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "kns-spec/1.0")
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("%s: %s", u, res.Status)
	}
	return body, nil
}
