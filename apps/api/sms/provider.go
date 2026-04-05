package sms

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/emmanuella-codes/olu/config"
)

type Message struct {
	To   string
	Body string
	From string
	Type string
}

type Result struct {
	MessageID  string
	Status     string
	RawMessage string
}

type Provider interface {
	Send(ctx context.Context, msg Message) (*Result, error)
	Name() string
}

func Build(cfg *config.Config) (Provider, error) {
	switch strings.ToLower(strings.TrimSpace(cfg.SMSProvider)) {
	case "", "mock":
		return NewMock(cfg.SMSBaseURL, cfg.SMSFrom, time.Duration(cfg.SMSTimeoutSec)*time.Second), nil
	default:
		return nil, fmt.Errorf("sms: unsupported provider %q", cfg.SMSProvider)
	}
}
