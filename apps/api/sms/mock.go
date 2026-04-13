package sms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type MockProvider struct {
	baseURL     string
	defaultFrom string
	client      *http.Client
}

func NewMock(baseURL, defaultFrom string, timeout time.Duration) *MockProvider {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}

	return &MockProvider{
		baseURL:     strings.TrimRight(baseURL, "/"),
		defaultFrom: defaultFrom,
		client:      &http.Client{Timeout: timeout},
	}
}

func (p *MockProvider) Name() string {
	return "mock"
}

func (p *MockProvider) Send(ctx context.Context, msg Message) (*Result, error) {
	from := msg.From
	if from == "" {
		from = p.defaultFrom
	}

	payload := struct {
		To   string `json:"to"`
		From string `json:"from"`
		SMS  string `json:"sms"`
	}{
		To:   msg.To,
		From: from,
		SMS:  msg.Body,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("sms mock: marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/api/sms/send", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("sms mock: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sms mock: send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("sms mock: read response: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("sms mock: unexpected status %d: %s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}

	var out struct {
		MessageID string `json:"message_id"`
		Message   string `json:"message"`
	}
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("sms mock: decode response: %w", err)
	}

	return &Result{
		MessageID:  out.MessageID,
		Status:     out.Message,
		RawMessage: string(respBody),
	}, nil
}
