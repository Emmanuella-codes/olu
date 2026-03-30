package validator

import "strings"

var validParties = map[string]bool{
	"app": true, "a": true, "apm": true, "apga": true, "apc": true,
	"adp": true, "adc": true, "aac": true, "aa": true, "bp": true,
	"dla": true, "lp": true, "nrm": true, "nnpp": true, "pdp": true,
	"prp": true, "sdp": true, "ypp": true, "yp": true, "zlp": true,
}

func IsValidParty(party string) bool {
	return validParties[NormalizeParty(party)]
}

func NormalizeParty(party string) string {
	return strings.ToLower(strings.TrimSpace(party))
}
