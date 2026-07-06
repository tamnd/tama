// Package gen talks to the OpenAI-compatible endpoint that generates course
// content: GET {base}/models for connectivity checks and POST
// {base}/chat/completions for completions. M1 ships the client; the prompts
// and pack building land in M8.
package gen

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/tamnd/tama/pkg/config"
)

// maxAttempts bounds retries on 429 and 5xx responses.
const maxAttempts = 3

// Client is a thin, retrying HTTP client for the OpenAI-compatible endpoint.
type Client struct {
	base  string
	key   string
	model string
	hc    *http.Client

	// backoffBase seeds the jittered retry wait; tests shrink it.
	backoffBase time.Duration
}

// New builds a client from the resolved [llm] config. The API key rides in
// an Authorization: Bearer header only when non-empty.
func New(cfg config.LLM) *Client {
	return &Client{
		base:  strings.TrimRight(cfg.BaseURL, "/"),
		key:   cfg.APIKey,
		model: cfg.Model,
		hc: &http.Client{
			// Total per-request budget; connect-to-first-byte is bounded
			// separately by the transport below.
			Timeout: cfg.RequestTimeout,
			Transport: &http.Transport{
				DialContext:           (&net.Dialer{Timeout: cfg.ConnectTimeout}).DialContext,
				ResponseHeaderTimeout: cfg.ConnectTimeout,
			},
		},
		backoffBase: 500 * time.Millisecond,
	}
}

// Model is one entry from GET /models.
type Model struct {
	ID string `json:"id"`
}

// Ping lists the endpoint's models, proving connectivity and auth. It backs
// `tama gen --dry-run`.
func (c *Client) Ping(ctx context.Context) ([]Model, error) {
	body, err := c.do(ctx, http.MethodGet, "/models", nil)
	if err != nil {
		return nil, err
	}
	var out struct {
		Data []Model `json:"data"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("gen: decode models: %w", err)
	}
	return out.Data, nil
}

// Message is one chat turn.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// CompleteRequest shapes one completion call. JSONMode asks the endpoint for
// a json_object response format.
type CompleteRequest struct {
	Messages    []Message
	Temperature float64
	MaxTokens   int
	JSONMode    bool
}

// Complete runs one non-streaming completion and returns the first choice's
// content. Streaming is deferred to M8.
func (c *Client) Complete(ctx context.Context, req CompleteRequest) (string, error) {
	if c.model == "" {
		return "", fmt.Errorf("gen: no model configured, set TAMA_LLM_MODEL or [llm] model")
	}
	payload := map[string]any{
		"model":       c.model,
		"messages":    req.Messages,
		"temperature": req.Temperature,
	}
	if req.MaxTokens > 0 {
		payload["max_tokens"] = req.MaxTokens
	}
	if req.JSONMode {
		payload["response_format"] = map[string]string{"type": "json_object"}
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	body, err := c.do(ctx, http.MethodPost, "/chat/completions", raw)
	if err != nil {
		return "", err
	}
	var out struct {
		Choices []struct {
			Message Message `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return "", fmt.Errorf("gen: decode completion: %w", err)
	}
	if len(out.Choices) == 0 {
		return "", fmt.Errorf("gen: completion returned no choices")
	}
	return out.Choices[0].Message.Content, nil
}

// do issues one request with up to maxAttempts tries. 429 and 5xx retry with
// jittered backoff, honoring Retry-After when the endpoint sends one.
func (c *Client) do(ctx context.Context, method, path string, body []byte) ([]byte, error) {
	var lastErr error
	for attempt := 0; attempt < maxAttempts; attempt++ {
		if attempt > 0 {
			if err := sleep(ctx, c.wait(attempt, lastErr)); err != nil {
				return nil, err
			}
		}

		var r io.Reader
		if body != nil {
			r = bytes.NewReader(body)
		}
		req, err := http.NewRequestWithContext(ctx, method, c.base+path, r)
		if err != nil {
			return nil, err
		}
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		if c.key != "" {
			req.Header.Set("Authorization", "Bearer "+c.key)
		}

		resp, err := c.hc.Do(req)
		if err != nil {
			// Network errors surface immediately; the retry loop is for
			// the endpoint telling us to back off.
			return nil, fmt.Errorf("gen: %s %s: %w", method, path, err)
		}
		payload, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("gen: read response: %w", err)
		}

		if resp.StatusCode == http.StatusOK {
			return payload, nil
		}
		lastErr = &httpError{status: resp.StatusCode, retryAfter: parseRetryAfter(resp.Header.Get("Retry-After"))}
		if resp.StatusCode != http.StatusTooManyRequests && resp.StatusCode < 500 {
			return nil, fmt.Errorf("gen: %s %s: endpoint returned %d: %s", method, path, resp.StatusCode, strings.TrimSpace(string(payload)))
		}
	}
	return nil, fmt.Errorf("gen: %s %s: gave up after %d attempts: %w", method, path, maxAttempts, lastErr)
}

type httpError struct {
	status     int
	retryAfter time.Duration
}

func (e *httpError) Error() string {
	return fmt.Sprintf("endpoint returned %d", e.status)
}

// wait picks the backoff before the given attempt: Retry-After verbatim if
// present, otherwise base*2^n with up to 50% jitter.
func (c *Client) wait(attempt int, lastErr error) time.Duration {
	if he, ok := lastErr.(*httpError); ok && he.retryAfter > 0 {
		return he.retryAfter
	}
	d := c.backoffBase << (attempt - 1)
	return d + time.Duration(rand.Int64N(int64(d)/2+1))
}

func parseRetryAfter(v string) time.Duration {
	if v == "" {
		return 0
	}
	if secs, err := strconv.Atoi(v); err == nil && secs >= 0 {
		return time.Duration(secs) * time.Second
	}
	if at, err := http.ParseTime(v); err == nil {
		return time.Until(at)
	}
	return 0
}

func sleep(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}
