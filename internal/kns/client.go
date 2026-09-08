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
