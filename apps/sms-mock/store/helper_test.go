package store

import "testing"

// --- extractOTP ---

func TestExtractOTP_SixDigits(t *testing.T) {
	got := extractOTP("Your voting code is 123456.")
	if got != "123456" {
		t.Fatalf("expected 123456, got %q", got)
	}
}

func TestExtractOTP_CodeAtStart(t *testing.T) {
	got := extractOTP("654321 is your OTP")
	if got != "654321" {
		t.Fatalf("expected 654321, got %q", got)
	}
}

func TestExtractOTP_NoMatch(t *testing.T) {
	got := extractOTP("Vote confirmed! No numeric code here.")
	if got != "" {
		t.Fatalf("expected empty, got %q", got)
	}
}

func TestExtractOTP_FiveDigitsNotMatched(t *testing.T) {
	// Five-digit numbers must not be extracted.
	got := extractOTP("Code: 12345")
	if got != "" {
		t.Fatalf("expected empty for 5-digit number, got %q", got)
	}
}

func TestExtractOTP_SevenDigitsNotMatched(t *testing.T) {
	// Seven-digit numbers must not be extracted.
	got := extractOTP("Code: 1234567")
	if got != "" {
		t.Fatalf("expected empty for 7-digit number, got %q", got)
	}
}

func TestExtractOTP_EmptyBody(t *testing.T) {
	got := extractOTP("")
	if got != "" {
		t.Fatalf("expected empty, got %q", got)
	}
}

func TestExtractOTP_LeadingZero(t *testing.T) {
	// Six-digit codes that start with zero are valid OTPs.
	got := extractOTP("Your OTP is 012345.")
	if got != "012345" {
		t.Fatalf("expected 012345, got %q", got)
	}
}

// --- detectChannel ---

func TestDetectChannel_OTPByVotingCode(t *testing.T) {
	ch := detectChannel("Your voting code is 123456")
	if ch != ChannelOTP {
		t.Fatalf("expected otp, got %q", ch)
	}
}

func TestDetectChannel_OTPByOTPKeyword(t *testing.T) {
	ch := detectChannel("OTP: 654321")
	if ch != ChannelOTP {
		t.Fatalf("expected otp, got %q", ch)
	}
}

func TestDetectChannel_OTPCaseInsensitive(t *testing.T) {
	ch := detectChannel("VOTING CODE 123456")
	if ch != ChannelOTP {
		t.Fatalf("expected otp, got %q", ch)
	}
}

func TestDetectChannel_ConfirmByVoteConfirmed(t *testing.T) {
	ch := detectChannel("Vote confirmed! Confirmation ID: abc123")
	if ch != ChannelConfirm {
		t.Fatalf("expected confirm, got %q", ch)
	}
}

func TestDetectChannel_ConfirmByConfirmationID(t *testing.T) {
	ch := detectChannel("Your confirmation id is XYZ789")
	if ch != ChannelConfirm {
		t.Fatalf("expected confirm, got %q", ch)
	}
}

func TestDetectChannel_RejectByInvalid(t *testing.T) {
	ch := detectChannel("Invalid candidate code. Please try again.")
	if ch != ChannelReject {
		t.Fatalf("expected reject, got %q", ch)
	}
}

func TestDetectChannel_RejectByCouldNot(t *testing.T) {
	ch := detectChannel("We could not process your vote.")
	if ch != ChannelReject {
		t.Fatalf("expected reject, got %q", ch)
	}
}

func TestDetectChannel_RejectByAlready(t *testing.T) {
	ch := detectChannel("You have already voted.")
	if ch != ChannelReject {
		t.Fatalf("expected reject, got %q", ch)
	}
}

func TestDetectChannel_Generic(t *testing.T) {
	ch := detectChannel("Welcome to the election system.")
	if ch != ChannelGeneric {
		t.Fatalf("expected generic, got %q", ch)
	}
}

func TestDetectChannel_Empty(t *testing.T) {
	ch := detectChannel("")
	if ch != ChannelGeneric {
		t.Fatalf("expected generic for empty body, got %q", ch)
	}
}

// --- normalizePhone (store package) ---

func TestNormalizePhone_LocalLeadingZero(t *testing.T) {
	got := normalizePhone("08012345678")
	if got != "2348012345678" {
		t.Fatalf("expected 2348012345678, got %q", got)
	}
}

func TestNormalizePhone_PlusStripped(t *testing.T) {
	got := normalizePhone("+2348012345678")
	if got != "2348012345678" {
		t.Fatalf("expected 2348012345678, got %q", got)
	}
}

func TestNormalizePhone_SpacesAndHyphens(t *testing.T) {
	got := normalizePhone("0801 234-5678")
	if got != "2348012345678" {
		t.Fatalf("expected 2348012345678, got %q", got)
	}
}

func TestNormalizePhone_AlreadyNormalized(t *testing.T) {
	got := normalizePhone("2348012345678")
	if got != "2348012345678" {
		t.Fatalf("expected 2348012345678, got %q", got)
	}
}

func TestNormalizePhone_Empty(t *testing.T) {
	got := normalizePhone("")
	if got != "" {
		t.Fatalf("expected empty, got %q", got)
	}
}
