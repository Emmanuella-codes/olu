package store

import (
	"fmt"
	"sync"
	"testing"
)

// --- Add ---

func TestAdd_MessageStoredWithCorrectFields(t *testing.T) {
	s := New()
	msg := s.Add("+2348012345678", "OLU", "Your voting code is 123456")

	if msg.To != "+2348012345678" {
		t.Fatalf("expected To +2348012345678, got %q", msg.To)
	}
	if msg.From != "OLU" {
		t.Fatalf("expected From OLU, got %q", msg.From)
	}
	if msg.Channel != ChannelOTP {
		t.Fatalf("expected otp channel, got %q", msg.Channel)
	}
	if msg.OTPCode != "123456" {
		t.Fatalf("expected OTPCode 123456, got %q", msg.OTPCode)
	}
}

func TestAdd_IDIncrements(t *testing.T) {
	s := New()
	m1 := s.Add("phone1", "OLU", "msg one")
	m2 := s.Add("phone2", "OLU", "msg two")
	if m1.ID != "msg_0001" {
		t.Fatalf("expected msg_0001, got %q", m1.ID)
	}
	if m2.ID != "msg_0002" {
		t.Fatalf("expected msg_0002, got %q", m2.ID)
	}
}

// --- AddInbound ---

func TestAddInbound_ToIsOLU(t *testing.T) {
	s := New()
	msg := s.AddInbound("+2348012345678", "VOTE A1")

	if msg.To != "OLU" {
		t.Fatalf("expected To OLU, got %q", msg.To)
	}
	if msg.Channel != ChannelInbound {
		t.Fatalf("expected inbound channel, got %q", msg.Channel)
	}
}

// --- All ---

func TestAll_ReturnsDefensiveCopy(t *testing.T) {
	s := New()
	s.Add("+2348012345678", "OLU", "msg one")

	all := s.All()
	all[0].Body = "mutated"

	// Internal state should be unaffected.
	internal := s.All()
	if internal[0].Body == "mutated" {
		t.Fatal("All() returned a reference to internal slice")
	}
}

func TestAll_EmptyStore(t *testing.T) {
	s := New()
	if msgs := s.All(); len(msgs) != 0 {
		t.Fatalf("expected empty, got %d messages", len(msgs))
	}
}

// --- ByPhone ---

func TestByPhone_MatchesNormalizedForms(t *testing.T) {
	s := New()
	s.Add("+2348012345678", "OLU", "hello")
	s.Add("+2349099999999", "OLU", "other")

	// Query with local format — should still match the +234 stored message.
	msgs := s.ByPhone("08012345678")
	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(msgs))
	}
}

func TestByPhone_NoMatchReturnsEmpty(t *testing.T) {
	s := New()
	s.Add("+2348012345678", "OLU", "hello")

	msgs := s.ByPhone("08099999999")
	if len(msgs) != 0 {
		t.Fatalf("expected 0 messages, got %d", len(msgs))
	}
}

// --- Latest ---

func TestLatest_ReturnsNewestMessage(t *testing.T) {
	s := New()
	s.Add("+2348012345678", "OLU", "first")
	s.Add("+2348012345678", "OLU", "second")

	m := s.Latest("+2348012345678")
	if m == nil {
		t.Fatal("expected a message, got nil")
	}
	if m.Body != "second" {
		t.Fatalf("expected body %q, got %q", "second", m.Body)
	}
}

func TestLatest_ReturnsNilWhenNoMatch(t *testing.T) {
	s := New()
	if m := s.Latest("+2348099999999"); m != nil {
		t.Fatalf("expected nil, got %+v", m)
	}
}

// --- LatestOTP ---

func TestLatestOTP_ReturnsCode(t *testing.T) {
	s := New()
	s.Add("+2348012345678", "OLU", "Your voting code is 987654")

	code := s.LatestOTP("+2348012345678")
	if code != "987654" {
		t.Fatalf("expected 987654, got %q", code)
	}
}

func TestLatestOTP_EmptyWhenNoOTPMessage(t *testing.T) {
	s := New()
	s.Add("+2348012345678", "OLU", "Vote confirmed! Confirmation ID: abc")

	code := s.LatestOTP("+2348012345678")
	if code != "" {
		t.Fatalf("expected empty, got %q", code)
	}
}

func TestLatestOTP_ReturnsNewest(t *testing.T) {
	s := New()
	s.Add("+2348012345678", "OLU", "Your voting code is 111111")
	s.Add("+2348012345678", "OLU", "Your voting code is 222222")

	code := s.LatestOTP("+2348012345678")
	if code != "222222" {
		t.Fatalf("expected 222222, got %q", code)
	}
}

// --- Clear ---

func TestClear_ResetsStoreAndCounter(t *testing.T) {
	s := New()
	s.Add("+2348012345678", "OLU", "msg")
	s.Clear()

	if msgs := s.All(); len(msgs) != 0 {
		t.Fatalf("expected 0 after clear, got %d", len(msgs))
	}

	// Counter resets: next message gets ID msg_0001 again.
	msg := s.Add("+2348012345678", "OLU", "after clear")
	if msg.ID != "msg_0001" {
		t.Fatalf("expected msg_0001 after clear, got %q", msg.ID)
	}
}

// --- Stats ---

func TestStats_CountsChannels(t *testing.T) {
	s := New()
	s.Add("+2348012345678", "OLU", "Your voting code is 123456") // otp
	s.Add("+2348012345678", "OLU", "Vote confirmed! Confirmation ID: X")  // confirm
	s.AddInbound("+2348012345678", "VOTE A1")                              // inbound

	stats := s.Stats()
	if stats["total"] != 3 {
		t.Fatalf("expected total 3, got %d", stats["total"])
	}
	if stats["otp"] != 1 {
		t.Fatalf("expected 1 otp, got %d", stats["otp"])
	}
	if stats["confirm"] != 1 {
		t.Fatalf("expected 1 confirm, got %d", stats["confirm"])
	}
	if stats["inbound"] != 1 {
		t.Fatalf("expected 1 inbound, got %d", stats["inbound"])
	}
}

// --- maxMessages cap ---

func TestAdd_CapsAtMaxMessages(t *testing.T) {
	s := New()
	for i := range maxMessages + 10 {
		s.Add(fmt.Sprintf("+23480%08d", i), "OLU", "msg")
	}
	if msgs := s.All(); len(msgs) != maxMessages {
		t.Fatalf("expected %d messages (cap), got %d", maxMessages, len(msgs))
	}
}

// --- Concurrency ---

func TestAdd_ConcurrentAddsAreSafe(t *testing.T) {
	s := New()
	const goroutines = 50

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := range goroutines {
		go func(i int) {
			defer wg.Done()
			s.Add(fmt.Sprintf("+23480%08d", i), "OLU", "concurrent msg")
		}(i)
	}
	wg.Wait()

	if msgs := s.All(); len(msgs) != goroutines {
		t.Fatalf("expected %d messages, got %d", goroutines, len(msgs))
	}
}

func TestAll_ConcurrentReadAndWrite(t *testing.T) {
	s := New()
	var wg sync.WaitGroup
	const n = 30

	wg.Add(n * 2)
	for i := range n {
		go func(i int) {
			defer wg.Done()
			s.Add(fmt.Sprintf("+23480%08d", i), "OLU", "msg")
		}(i)
		go func() {
			defer wg.Done()
			_ = s.All()
		}()
	}
	wg.Wait()
}

// --- detectSourceChannel ---

func TestDetectSourceChannel_DefaultsToWeb(t *testing.T) {
	s := New()
	// No prior messages for this phone.
	ch := s.detectSourceChannel("+2348012345678")
	if ch != SourceChannelWeb {
		t.Fatalf("expected web, got %q", ch)
	}
}

func TestDetectSourceChannel_SMSAfterInbound(t *testing.T) {
	s := New()
	// Simulate inbound vote from this phone.
	s.AddInbound("+2348012345678", "VOTE A1")
	ch := s.detectSourceChannel("+2348012345678")
	if ch != SourceChannelSMS {
		t.Fatalf("expected sms after inbound, got %q", ch)
	}
}

func TestDetectSourceChannel_WebAfterOTP(t *testing.T) {
	s := New()
	// OTP was sent TO this phone — web voter.
	s.Add("+2348012345678", "OLU", "Your voting code is 123456")
	ch := s.detectSourceChannel("+2348012345678")
	if ch != SourceChannelWeb {
		t.Fatalf("expected web after otp sent, got %q", ch)
	}
}
