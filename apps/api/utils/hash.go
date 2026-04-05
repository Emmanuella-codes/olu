package utils

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

// returns a hex-encoded SHA-256 hash of the voter's
// phone number. We store only the hash — never the raw phone number.
// A salt ties the hash to this application so rainbow tables from other
// systems are useless against our voter_registry.

func HashVoterIdentity(phone, salt string) string {
	normalized := normalizePhone(phone)
	input := fmt.Sprintf("%s:%s", salt, normalized)
	sum := sha256.Sum256([]byte(input))
	return fmt.Sprintf("%x", sum)
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
