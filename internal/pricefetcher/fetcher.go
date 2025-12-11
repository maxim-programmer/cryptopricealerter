package pricefetcher

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Fetcher struct {
	HTTPClient *http.Client
	BaseURL    string
}

func NewFetcher() *Fetcher {
	return &Fetcher{
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		BaseURL: "https://api.coingecko.com/api/v3/simple/price?vs_currencies=usd&ids=",
	}
}

func (f *Fetcher) GetPrices(ctx context.Context, symbols string, m map[string]Price, api_key string) (map[string]Price, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", f.BaseURL+symbols+"&x_cg_demo_api_key="+api_key+"&vs_currencies=usd", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}
