package store

import (
	"fmt"
	"sync"
	"time"
)

type Channel string

const (
	ChannelOTP     Channel = "otp"
	ChannelConfirm Channel = "confirm"
	ChannelReject  Channel = "reject"
	ChannelGeneric Channel = "generic"
)

type Message struct {
	ID      string    `json:"id"`
	To      string    `json:"to"`
	From    string    `json:"from"`
	Body    string    `json:"body"`
	Channel Channel   `json:"channel"`
	OTPCode string    `json:"otp_code,omitempty"`
	SentAt  time.Time `json:"sent_at"`
}

const maxMessages = 500

type Store struct {
	mu       sync.RWMutex
	messages []Message
	counter  int
}

func New() *Store {
	return &Store{messages: []Message{}}
}

func (s *Store) Add(to, from, body string) Message {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.counter++
	msg := Message{
		ID:      fmt.Sprintf("msg_%04d", s.counter),
		To:      to,
		From:    from,
		Body:    body,
		Channel: detectChannel(body),
		OTPCode: extractOTP(body),
		SentAt:  time.Now(),
	}
	s.messages = append(s.messages, msg)
	if len(s.messages) > maxMessages {
		s.messages = s.messages[len(s.messages)-maxMessages:]
	}
	return msg
}
func (s *Store) All() []Message {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Message, len(s.messages))
	copy(out, s.messages)
	return out
}

func (s *Store) ByPhone(phone string) []Message {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []Message{}
	for _, m := range s.messages {
		if m.To == phone || normalizePhone(m.To) == normalizePhone(phone) {
			out = append(out, m)
		}
	}
	return out
}

func (s *Store) Latest(phone string) *Message {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i := len(s.messages) - 1; i >= 0; i-- {
		m := s.messages[i]
		if m.To == phone || normalizePhone(m.To) == normalizePhone(phone) {
			return &m
		}
	}
	return nil
}

func (s *Store) LatestOTP(phone string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i := len(s.messages) - 1; i >= 0; i-- {
		m := s.messages[i]
		if (m.To == phone || normalizePhone(m.To) == normalizePhone(phone)) &&
			m.Channel == ChannelOTP &&
			m.OTPCode != "" {
			return m.OTPCode
		}
	}
	return ""
}

func (s *Store) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.messages = []Message{}
	s.counter = 0
}

func (s *Store) Stats() map[string]int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	stats := map[string]int{"total": len(s.messages)}
	for _, m := range s.messages {
		stats[string(m.Channel)]++
	}
	return stats
}
