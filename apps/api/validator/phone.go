package validator

import (
	"regexp"
	"strings"
)

var phoneRegex = regexp.MustCompile(
	`^(\+234|0)(7[01]|8[01]|9[01])\d{8}$`,
)

func IsValidPhone(phone string) bool {
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	return phoneRegex.MatchString(phone)
}

func ToE164(phone string) string {
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")

	if strings.HasPrefix(phone, "0") {
		return "+234" + phone[1:]
	}
	if strings.HasPrefix(phone, "234") {
		return "+" + phone
	}
	return phone
}
