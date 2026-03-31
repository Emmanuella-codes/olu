package store

import (
	"regexp"
	"strings"
)

var otpRegex = regexp.MustCompile(`\b(\d{6})\b`)

func extractOTP(body string) string {
	m := otpRegex.FindStringSubmatch(body)
	if len(m) > 1 {
		return m[1]
	}
	return ""
}

// classify msg based on the content
func detectChannel(body string) Channel {
	lower := strings.ToLower(body)
	switch {
	case strings.Contains(lower, "voting code") || strings.Contains(lower, "otp"):
		return ChannelOTP
	case strings.Contains(lower, "vote confirmed") || strings.Contains(lower, "confirmation id"):
		return ChannelConfirm
	case strings.Contains(lower, "invalid") || strings.Contains(lower, "could not") || strings.Contains(lower, "already"):
		return ChannelReject
	default:
		return ChannelGeneric
	}
}

func normalizePhone(phone string) string {
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.TrimPrefix(phone, "+")

	if strings.HasPrefix(phone, "0") {
		phone = "234" + phone[1:]
	}

	return phone
}
