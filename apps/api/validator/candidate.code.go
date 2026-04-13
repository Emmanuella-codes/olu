package validator

import (
	"regexp"
	"strings"
)

// Candidate codes are 1 uppercase letter + 1-2 digits, e.g. A1, B12.
var candidateCodeRegex = regexp.MustCompile(`^[A-Z]\d{1,2}$`)

func IsValidCandidateCode(code string) bool {
	return candidateCodeRegex.MatchString(strings.ToUpper(strings.TrimSpace(code)))
}

func NormalizeCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}
