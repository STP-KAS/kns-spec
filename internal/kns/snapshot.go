package kns

import (
	"strings"

	"github.com/STP-KAS/kns-spec/internal/overlay"
)

const ProfileKeys = "redirectUrl,bio,avatarUrl,website,x,github,telegram,discord,email,banner,ipfs,kfs,contenthash,peer,onion,agent,kas,noise"

type Snapshot struct {
	Name      string           `json:"name"`
	Owner     string           `json:"owner"`
	AssetID   string           `json:"assetId"`
	TxID      string           `json:"txid,omitempty"`
	Verified  bool             `json:"verified"`
	Pay       string           `json:"pay"`
	PayURI    string           `json:"payUri"`
	Web       string           `json:"web"`
	Run       string           `json:"run"`
	Session   string           `json:"session"`
	Primary   string           `json:"primary,omitempty"`
	PrimaryOK bool             `json:"primaryOk"`
	KNS       string           `json:"kns"`
	Records   overlay.Records  `json:"records"`
	Warn      string           `json:"warn"`
}

func (c *Client) Snapshot(name string) (*Snapshot, error) {
	name = strings.ToLower(strings.TrimSpace(name))
	if !strings.HasSuffix(name, ".kas") {
		name += ".kas"
	}
	own, err := c.Owner(name)
	if err != nil {
		return nil, err
	}
	rec := overlay.Records{KAS: own.Owner}
	if p, err := c.Profile(own.AssetID); err == nil && p != nil {
		got := overlay.FromMap(p.Profile)
		if got.KAS == "" {
			got.KAS = own.Owner
		}
		rec = got
	}
	pay := rec.PayAddress()
	if pay == "" {
		pay = own.Owner
		rec.KAS = own.Owner
	}
	s := &Snapshot{
		Name:     name,
		Owner:    own.Owner,
		AssetID:  own.AssetID,
		Pay:      pay,
		PayURI:   overlay.PayURI(pay),
		Web:      rec.Web(),
		Run:      rec.Run(),
		Session:  rec.Session(),
		KNS:      overlay.URI(name, ""),
		Records:  rec,
		Warn:     overlay.ResolveWarning,
		Verified: own.Owner != "",
	}
	if a, err := c.Asset(name); err == nil && a != nil {
		s.TxID = a.TransactionID
	}
	if p, err := c.Primary(own.Owner); err == nil && p != nil && p.Domain != nil {
		s.Primary = p.Domain.FullName
		s.PrimaryOK = strings.EqualFold(p.Domain.FullName, name)
	}
	return s, nil
}
