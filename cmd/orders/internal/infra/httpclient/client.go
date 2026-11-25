package httpclient

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/edmilsonmedeiross/go-order-generator/cmd/orders/internal/config"
)

type Client struct {
	client     *http.Client
	timeout    time.Duration
	maxRetries int
}

func New(cfg config.HTTPConfig) *Client {
	return &Client{client: &http.Client{}, timeout: cfg.Timeout, maxRetries: cfg.MaxRetries}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(req.Context(), c.timeout)
	for attempt := 0; ; attempt++ {
		response, err := c.client.Do(req.Clone(ctx))
		if ctx.Err() != nil {
			if response != nil {
				response.Body.Close()
			}
			cancel()
			return nil, ctx.Err()
		}
		var networkErr net.Error
		retryable := req.Method == http.MethodGet && attempt < c.maxRetries &&
			((err != nil && errors.As(err, &networkErr) && (networkErr.Timeout() || networkErr.Temporary())) ||
				(err == nil && (response.StatusCode == 429 || response.StatusCode == 502 || response.StatusCode == 503 || response.StatusCode == 504)))
		if !retryable {
			if err != nil {
				cancel()
				return nil, err
			}
			response.Body = &cancelBody{ReadCloser: response.Body, cancel: cancel}
			return response, nil
		}
		if response != nil {
			response.Body.Close()
		}
		timer := time.NewTimer(time.Duration(attempt+1) * 100 * time.Millisecond)
		select {
		case <-ctx.Done():
			timer.Stop()
			cancel()
			return nil, ctx.Err()
		case <-timer.C:
		}
	}
}

type cancelBody struct {
	io.ReadCloser
	cancel context.CancelFunc
}

func (b *cancelBody) Close() error { err := b.ReadCloser.Close(); b.cancel(); return err }
