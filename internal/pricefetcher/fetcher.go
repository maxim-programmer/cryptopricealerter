package pricefetcher

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

type fetcher struct {
	httpClient *http.Client
	baseURL    []string
}

func NewFetcher() *fetcher {
	return &fetcher{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: []string{"https://api.coingecko.com/api/v3/simple/price?vs_currencies=usd&ids=", "&x_cg_demo_api_key=", "&vs_currencies=usd"},
	}
}

func (f *fetcher) GetPrices(ctx context.Context, symbols []string, m map[string]Price, api_key string) (map[string]Price, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", f.baseURL[0]+strings.Join(symbols, ",")+f.baseURL[1]+api_key+f.baseURL[2], nil)
	if err != nil {
		return nil, err
	}

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &m); err != nil {
		return nil, err
	}

	return m, nil
}
