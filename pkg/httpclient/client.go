package httpclient

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

type HTTPClient interface {
	Post(ctx context.Context, url string, contentType string, body []byte, headers map[string]string) (*http.Response, error)
	Get(ctx context.Context, url string, headers map[string]string) (*http.Response, error)
	Put(ctx context.Context, url string, contentType string, body []byte, headers map[string]string) (*http.Response, error)
	Delete(ctx context.Context, url string, headers map[string]string) (*http.Response, error)
	Patch(ctx context.Context, url string, contentType string, body []byte, headers map[string]string) (*http.Response, error)
}

type Client struct {
	HTTP   *http.Client
	tracer trace.Tracer
}

func NewClient(tracer trace.Tracer) *Client {
	return &Client{
		HTTP: &http.Client{
			Transport: http.DefaultTransport,
		},
		tracer: tracer,
	}
}

func (c *Client) Post(ctx context.Context, url string, contentType string, body []byte, headers map[string]string) (*http.Response, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))

	_, span := c.tracer.Start(ctx, req.Method+" "+req.URL.String())
	defer span.End()

	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to create POST request: %w", err)
	}

	req.Header.Set("Content-Type", contentType)
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	span.SetAttributes(
		semconv.HTTPMethodKey.String(req.Method),
		semconv.HTTPURLKey.String(req.URL.String()),
		semconv.HTTPTargetKey.String(req.URL.Path),
		semconv.HTTPHostKey.String(req.Host),
		semconv.HTTPUserAgentKey.String(req.UserAgent()),
	)

	res, err := c.HTTP.Do(req)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to execute POST request: %w", err)
	}
	defer res.Body.Close()

	span.SetAttributes(
		semconv.HTTPStatusCodeKey.Int(res.StatusCode),
	)

	return res, nil
}

func (c *Client) Get(ctx context.Context, url string, headers map[string]string) (*http.Response, error) {
	ctx, span := c.tracer.Start(ctx, "HTTP GET")
	defer span.End()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to create GET request: %w", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	span.SetAttributes(
		semconv.HTTPMethodKey.String(req.Method),
		semconv.HTTPURLKey.String(req.URL.String()),
		semconv.HTTPTargetKey.String(req.URL.Path),
		semconv.HTTPHostKey.String(req.Host),
		semconv.HTTPUserAgentKey.String(req.UserAgent()),
	)

	res, err := c.HTTP.Do(req)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to execute GET request: %w", err)
	}
	defer res.Body.Close()

	span.SetAttributes(
		semconv.HTTPStatusCodeKey.Int(res.StatusCode),
	)

	return res, nil
}

func (c *Client) Put(ctx context.Context, url string, contentType string, body []byte, headers map[string]string) (*http.Response, error) {
	ctx, span := c.tracer.Start(ctx, "HTTP PUT")
	defer span.End()

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(body))
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to create PUT request: %w", err)
	}

	req.Header.Set("Content-Type", contentType)
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	span.SetAttributes(
		semconv.HTTPMethodKey.String(req.Method),
		semconv.HTTPURLKey.String(req.URL.String()),
		semconv.HTTPTargetKey.String(req.URL.Path),
		semconv.HTTPHostKey.String(req.Host),
		semconv.HTTPUserAgentKey.String(req.UserAgent()),
	)

	res, err := c.HTTP.Do(req)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to execute PUT request: %w", err)
	}
	defer res.Body.Close()

	span.SetAttributes(
		semconv.HTTPStatusCodeKey.Int(res.StatusCode),
	)

	return res, nil
}

func (c *Client) Delete(ctx context.Context, url string, headers map[string]string) (*http.Response, error) {
	ctx, span := c.tracer.Start(ctx, "HTTP DELETE")
	defer span.End()

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to create DELETE request: %w", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	span.SetAttributes(
		semconv.HTTPMethodKey.String(req.Method),
		semconv.HTTPURLKey.String(req.URL.String()),
		semconv.HTTPTargetKey.String(req.URL.Path),
		semconv.HTTPHostKey.String(req.Host),
		semconv.HTTPUserAgentKey.String(req.UserAgent()),
	)

	res, err := c.HTTP.Do(req)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to execute DELETE request: %w", err)
	}
	defer res.Body.Close()

	span.SetAttributes(
		semconv.HTTPStatusCodeKey.Int(res.StatusCode),
	)

	return res, nil
}

func (c *Client) Patch(ctx context.Context, url string, contentType string, body []byte, headers map[string]string) (*http.Response, error) {
	ctx, span := c.tracer.Start(ctx, "HTTP PATCH")
	defer span.End()

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewBuffer(body))
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to create PATCH request: %w", err)
	}

	req.Header.Set("Content-Type", contentType)
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	span.SetAttributes(
		semconv.HTTPMethodKey.String(req.Method),
		semconv.HTTPURLKey.String(req.URL.String()),
		semconv.HTTPTargetKey.String(req.URL.Path),
		semconv.HTTPHostKey.String(req.Host),
		semconv.HTTPUserAgentKey.String(req.UserAgent()),
	)

	res, err := c.HTTP.Do(req)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to execute PATCH request: %w", err)
	}
	defer res.Body.Close()

	span.SetAttributes(
		semconv.HTTPStatusCodeKey.Int(res.StatusCode),
	)

	return res, nil
}
