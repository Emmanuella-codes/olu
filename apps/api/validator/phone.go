package validator

import (
	"regexp"
	"strings"
)

var nigerianPhoneRegex = regexp.MustCompile(
	`^(\+234|0)(7[01]|8[01]|9[01])\d{8}$`,
)

func IsValidNigerianPhone(phone string) bool {
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	return nigerianPhoneRegex.MatchString(phone)
}
