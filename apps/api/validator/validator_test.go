package validator

import "testing"

func TestIsValidPhone(t *testing.T) {
	tests := []struct {
		name  string
		phone string
		want  bool
	}{
		{name: "local mtn", phone: "08012345678", want: true},
		{name: "international", phone: "+2348012345678", want: true},
		{name: "spaces and hyphen", phone: "0801 234-5678", want: true},
		{name: "hyphen separated local", phone: "0801-234-5678", want: true},
		{name: "unsupported prefix", phone: "06012345678", want: false},
		{name: "too short", phone: "080123", want: false},
		{name: "too long", phone: "080123456789", want: false},
		{name: "missing plus on country code is invalid for validation", phone: "2348012345678", want: false},
		{name: "plus234 wrong prefix", phone: "+23406012345678", want: false},
		{name: "whitespace only", phone: "   ", want: false},
		{name: "empty", phone: "", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidPhone(tt.phone); got != tt.want {
				t.Fatalf("IsValidPhone(%q) = %v, want %v", tt.phone, got, tt.want)
			}
		})
	}
}

func TestSupportedPhonePrefixes(t *testing.T) {
	// covers all four supported prefix groups: 70, 71, 80, 81, 90, 91
	for _, phone := range []string{
		"07012345678", "07112345678",
		"08012345678", "08112345678",
		"09012345678", "09112345678",
	} {
		t.Run(phone, func(t *testing.T) {
			if !IsValidPhone(phone) {
				t.Fatalf("expected %q to be valid", phone)
			}
		})
	}
}

func TestToE164(t *testing.T) {
	tests := []struct {
		name  string
		phone string
		want  string
	}{
		{name: "local", phone: "08012345678", want: "+2348012345678"},
		{name: "local with spaces", phone: "0801 234-5678", want: "+2348012345678"},
		{name: "local with hyphens", phone: "0801-234-5678", want: "+2348012345678"},
		{name: "country code without plus", phone: "2348012345678", want: "+2348012345678"},
		{name: "already e164", phone: "+2348012345678", want: "+2348012345678"},
		{name: "unrecognized format returned as-is", phone: "5551234567", want: "5551234567"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToE164(tt.phone); got != tt.want {
				t.Fatalf("ToE164(%q) = %q, want %q", tt.phone, got, tt.want)
			}
		})
	}
}

func TestCandidateCodeValidationAndNormalization(t *testing.T) {
	normCases := []struct {
		input string
		want  string
	}{
		{input: " a12 ", want: "A12"},
		{input: "A1", want: "A1"},   // already uppercase, no-op
		{input: "aB1", want: "AB1"}, // mixed case
	}
	for _, tc := range normCases {
		t.Run("normalize_"+tc.input, func(t *testing.T) {
			if got := NormalizeCode(tc.input); got != tc.want {
				t.Fatalf("NormalizeCode(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}

	tests := []struct {
		code string
		want bool
	}{
		{code: "A1", want: true},
		{code: "a12", want: true},
		{code: "Z1", want: true},   // Z is a valid letter
		{code: "Z99", want: true},  // upper boundary: letter + 2 digits
		{code: "A0", want: true},   // digit zero allowed
		{code: "A100", want: false}, // three digits rejected
		{code: "A123", want: false},
		{code: "AA1", want: false},
		{code: "1A", want: false},
		{code: "A!", want: false},
		{code: "   ", want: false},
		{code: "", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			if got := IsValidCandidateCode(tt.code); got != tt.want {
				t.Fatalf("IsValidCandidateCode(%q) = %v, want %v", tt.code, got, tt.want)
			}
		})
	}
}

func TestPartyValidationAndNormalization(t *testing.T) {
	normCases := []struct {
		input string
		want  string
	}{
		{input: " APC ", want: "apc"},
		{input: " Lp ", want: "lp"},
		{input: "   ", want: ""},
	}
	for _, tc := range normCases {
		t.Run("normalize_"+tc.input, func(t *testing.T) {
			if got := NormalizeParty(tc.input); got != tc.want {
				t.Fatalf("NormalizeParty(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}

	tests := []struct {
		party string
		want  bool
	}{
		{party: "APC", want: true},
		{party: " pdp ", want: true},
		{party: "lp", want: true},
		{party: "a", want: true}, // Accord — single char, easy to overlook
		{party: "unknown", want: false},
		{party: "   ", want: false},
		{party: "", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.party, func(t *testing.T) {
			if got := IsValidParty(tt.party); got != tt.want {
				t.Fatalf("IsValidParty(%q) = %v, want %v", tt.party, got, tt.want)
			}
		})
	}
}

// Sweeps all 20 parties in the allowlist to catch any map typos.
func TestAllValidPartiesAreAccepted(t *testing.T) {
	all20 := []string{
		"app", "a", "apm", "apga", "apc",
		"adp", "adc", "aac", "aa", "bp",
		"dla", "lp", "nrm", "nnpp", "pdp",
		"prp", "sdp", "ypp", "yp", "zlp",
	}
	if len(all20) != 20 {
		t.Fatalf("expected 20 parties, got %d — update this list", len(all20))
	}
	for _, party := range all20 {
		t.Run(party, func(t *testing.T) {
			if !IsValidParty(party) {
				t.Fatalf("expected party %q to be valid", party)
			}
		})
	}
}
