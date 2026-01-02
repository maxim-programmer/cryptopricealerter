package pricefetcher

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type fetcher struct {
	httpClient *http.Client
}

func NewFetcher() *fetcher {
	return &fetcher{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (f *fetcher) GetPrices(ctx context.Context, symbols []string, api_key string) (map[string]Price, error) {
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd&x_cg_demo_api_key=%s", strings.Join(symbols, ","), api_key)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("status " + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	m := make(map[string]Price)
	if err := json.Unmarshal(body, &m); err != nil {
		return nil, err
	}

	return m, nil
}
