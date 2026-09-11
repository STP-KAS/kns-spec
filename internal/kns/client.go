package kns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const Mainnet = "https://api.knsdomains.org/mainnet"

type Client struct {
	Base string
	HTTP *http.Client
}

func New(base string) *Client {
	if base == "" {
		base = Mainnet
	}
	return &Client{
		Base: strings.TrimRight(base, "/"),
		HTTP: &http.Client{Timeout: 18 * time.Second},
	}
}

type Envelope[T any] struct {
	Success bool   `json:"success"`
	Data    T      `json:"data"`
	Message string `json:"message"`
}

type Owner struct {
	ID      string `json:"id"`
	AssetID string `json:"assetId"`
	Asset   string `json:"asset"`
	Owner   string `json:"owner"`
}

func (c *Client) Owner(domain string) (*Owner, error) {
	var env Envelope[Owner]
	path := "/api/v1/" + url.PathEscape(domain) + "/owner"
	if err := c.get(path, &env); err != nil {
		return nil, err
	}
	if !env.Success || env.Data.Owner == "" {
		return nil, fmt.Errorf("not found")
	}
	return &env.Data, nil
}

type Asset struct {
	ID            string `json:"id"`
	AssetID       string `json:"assetId"`
	Asset         string `json:"asset"`
	Owner         string `json:"owner"`
	TransactionID string `json:"transactionId"`
	Status        string `json:"status"`
}

type CheckItem struct {
	Domain           string `json:"domain"`
	Available        bool   `json:"available"`
	IsReservedDomain bool   `json:"isReservedDomain"`
}

const CheckDummy = "kaspa:qzt9yuqceqvt2vk9dz7ddzayaa5flnenkymec59xvzm55ln3k72vgecxjhnjp"

func (c *Client) Check(names []string, address string) ([]CheckItem, error) {
	if address == "" {
		address = CheckDummy
	}
	var env Envelope[struct {
		Domains []CheckItem `json:"domains"`
	}]
	err := c.Post("/api/v1/domains/check", map[string]any{
		"domainNames": names,
		"address":     address,
	}, &env)
	if err != nil {
		return nil, err
	}
	if !env.Success {
		return nil, fmt.Errorf("check failed: %s", env.Message)
	}
	return env.Data.Domains, nil
}

type Profile struct {
	AssetID string         `json:"assetId"`
	Owner   string         `json:"owner"`
	Name    string         `json:"name"`
	TLD     string         `json:"tld"`
	Profile map[string]any `json:"profile"`
}

func (c *Client) Profile(assetID string) (*Profile, error) {
	var env Envelope[Profile]
	path := "/api/v1/domain/" + url.PathEscape(assetID) + "/profile?keys=" + url.QueryEscape(ProfileKeys)
	if err := c.get(path, &env); err != nil {
		return nil, err
	}
	if !env.Success {
		return nil, fmt.Errorf("profile failed: %s", env.Message)
	}
	return &env.Data, nil
}

type Primary struct {
	OwnerAddress  string `json:"ownerAddress"`
	InscriptionID string `json:"inscriptionId"`
	Domain        *struct {
		FullName   string `json:"fullName"`
		IsVerified bool   `json:"isVerified"`
	} `json:"domain"`
}

func (c *Client) Primary(owner string) (*Primary, error) {
	var env Envelope[Primary]
	path := "/api/v1/primary-name/" + url.PathEscape(owner)
	if err := c.get(path, &env); err != nil {
		return nil, err
	}
	if !env.Success {
		return nil, fmt.Errorf("primary not set")
	}
	return &env.Data, nil
}

func (c *Client) Asset(name string) (*Asset, error) {
	q := url.Values{}
	q.Set("asset", name)
	q.Set("type", "domain")
	q.Set("pageSize", "5")
	var env Envelope[struct {
		Assets []Asset `json:"assets"`
	}]
	if err := c.get("/api/v1/assets?"+q.Encode(), &env); err != nil {
		return nil, err
	}
	for i := range env.Data.Assets {
		if strings.EqualFold(env.Data.Assets[i].Asset, name) {
			return &env.Data.Assets[i], nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func (c *Client) get(path string, out any) error {
	req, err := http.NewRequest(http.MethodGet, c.Base+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "kns-spec/1.0")
	res, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode >= 400 {
		return fmt.Errorf("kns %s: %s", path, res.Status)
	}
	return json.Unmarshal(body, out)
}

func (c *Client) Post(path string, in, out any) error {
	raw, err := json.Marshal(in)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, c.Base+path, bytes.NewReader(raw))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "kns-spec/1.0")
	res, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode >= 400 {
		return fmt.Errorf("kns %s: %s", path, res.Status)
	}
	return json.Unmarshal(body, out)
}
