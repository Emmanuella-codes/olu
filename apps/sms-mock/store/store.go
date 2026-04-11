package store

import (
	"fmt"
	"sync"
	"time"
)

type Channel string
type SourceChannel string

const (
	ChannelInbound Channel = "inbound"
	ChannelOTP     Channel = "otp"
	ChannelConfirm Channel = "confirm"
	ChannelReject  Channel = "reject"
	ChannelGeneric Channel = "generic"

	SourceChannelWeb SourceChannel = "web"
	SourceChannelSMS SourceChannel = "sms"
)

type Message struct {
	ID      string    `json:"id"`
	To      string    `json:"to"`
	From    string    `json:"from"`
	Body    string    `json:"body"`
	Channel Channel   `json:"channel"`
	Source  string    `json:"source"`
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
	return s.addMessage(to, from, body, detectChannel(body), s.detectSourceChannel(to), extractOTP(body))
}

func (s *Store) AddInbound(from, body string) Message {
	return s.addMessage("OLU", from, body, ChannelInbound, SourceChannelSMS, "")
}

func (s *Store) addMessage(to, from, body string, channel Channel, source SourceChannel, otpCode string) Message {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.counter++
	msg := Message{
		ID:      fmt.Sprintf("msg_%04d", s.counter),
		To:      to,
		From:    from,
		Body:    body,
		Channel: channel,
		Source:  string(source),
		OTPCode: otpCode,
		SentAt:  time.Now(),
	}
	s.messages = append(s.messages, msg)
	if len(s.messages) > maxMessages {
		s.messages = s.messages[len(s.messages)-maxMessages:]
	}
	return msg
}

func (s *Store) detectSourceChannel(phone string) SourceChannel {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for i := len(s.messages) - 1; i >= 0; i-- {
		msg := s.messages[i]
		if matchesPhone(msg.From, phone) && msg.Channel == ChannelInbound {
			return SourceChannelSMS
		}
		if matchesPhone(msg.To, phone) && msg.Channel == ChannelOTP {
			return SourceChannelWeb
		}
	}

	return SourceChannelWeb
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
		if matchesPhone(m.To, phone) || matchesPhone(m.From, phone) {
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
		if matchesPhone(m.To, phone) || matchesPhone(m.From, phone) {
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
		if (matchesPhone(m.To, phone) || matchesPhone(m.From, phone)) &&
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

func matchesPhone(value, phone string) bool {
	return value == phone || normalizePhone(value) == normalizePhone(phone)
}
