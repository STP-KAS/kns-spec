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
	KNS          string          `json:"kns"`
	Records      overlay.Records `json:"records"`
	Warn         string          `json:"warn"`
	ProfileError string          `json:"profileError,omitempty"`
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
	rec := overlay.Records{}
	if p, err := c.Profile(own.AssetID); err != nil {
		s0 := &Snapshot{Name: name, Owner: own.Owner, AssetID: own.AssetID, Warn: overlay.ResolveWarning, ProfileError: err.Error()}
		s0.Pay = own.Owner
		s0.PayURI = overlay.PayURI(own.Owner)
		return s0, nil
	} else if p != nil {
		rec = overlay.FromMap(p.Profile)
	}
	if rec.KAS == "" {
		rec.KAS = own.Owner
	}
	pay := rec.PayAddress()
	if pay == "" {
		pay = own.Owner
		rec.KAS = own.Owner
	}
	s := &Snapshot{
		Name:    name,
		Owner:   own.Owner,
		AssetID: own.AssetID,
		Pay:     pay,
		PayURI:  overlay.PayURI(pay),
		Web:     rec.Web(),
		Run:     rec.Run(),
		Session: rec.Session(),
		KNS:     overlay.URI(name, ""),
		Records: rec,
		Warn:    overlay.ResolveWarning,
	}
	if a, err := c.Asset(name); err == nil && a != nil {
		s.TxID = a.TransactionID
		st := strings.ToLower(a.Status)
		s.Verified = st == "verified" || st == "valid" || st == "active"
	}
	if p, err := c.Primary(own.Owner); err == nil && p != nil && p.Domain != nil && p.Domain.FullName != "" {
		fwd, ferr := c.Owner(p.Domain.FullName)
		if ferr == nil && strings.EqualFold(fwd.Owner, own.Owner) {
			s.Primary = p.Domain.FullName
			s.PrimaryOK = strings.EqualFold(p.Domain.FullName, name)
			if p.Domain.IsVerified {
				s.Verified = true
			}
		}
	}
	return s, nil
}
