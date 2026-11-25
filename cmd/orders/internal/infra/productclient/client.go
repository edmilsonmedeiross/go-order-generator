package productclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/edmilsonmedeiross/go-order-generator/cmd/orders/internal/config"
	"github.com/edmilsonmedeiross/go-order-generator/cmd/orders/internal/usecase"
)

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	httpClient HTTPClient
	baseURL    string
}

func New(httpClient HTTPClient, cfg config.ProductConfig) *Client {
	return &Client{httpClient: httpClient, baseURL: cfg.BaseURL}
}

func (c *Client) GetProduct(ctx context.Context, id string) (usecase.ProductInfo, error) {
	address := strings.TrimRight(c.baseURL, "/") + "/v1/products/" + url.PathEscape(id)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, address, nil)
	if err != nil {
		return usecase.ProductInfo{}, err
	}
	response, err := c.httpClient.Do(req)
	if err != nil {
		return usecase.ProductInfo{}, fmt.Errorf("get product: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return usecase.ProductInfo{}, fmt.Errorf("product service returned HTTP %d", response.StatusCode)
	}
	var product usecase.ProductInfo
	if err := json.NewDecoder(response.Body).Decode(&product); err != nil {
		return usecase.ProductInfo{}, fmt.Errorf("decode product: %w", err)
	}
	return product, nil
}
