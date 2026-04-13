package utils

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"testing"
)


func TestNormalizePhone(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "local leading zero", input: "08012345678", want: "2348012345678"},
		{name: "already international", input: "+2348012345678", want: "2348012345678"},
		{name: "spaces stripped", input: "0801 234 5678", want: "2348012345678"},
		{name: "hyphens stripped", input: "0801-234-5678", want: "2348012345678"},
		{name: "spaces and hyphens", input: "0801 234-5678", want: "2348012345678"},
		{name: "plus stripped but no leading zero", input: "+2348012345678", want: "2348012345678"},
		{name: "no prefix transformation needed", input: "2348012345678", want: "2348012345678"},
		{name: "empty", input: "", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizePhone(tt.input)
			if got != tt.want {
				t.Fatalf("normalizePhone(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestHashVoterIdentity_IsSHA256Hex(t *testing.T) {
	h := HashVoterIdentity("+2348012345678", "test-salt")
	// SHA-256 produces 32 bytes = 64 hex characters.
	if len(h) != 64 {
		t.Fatalf("expected 64-char hex, got %d chars: %q", len(h), h)
	}
	for _, c := range h {
		if !strings.ContainsRune("0123456789abcdef", c) {
			t.Fatalf("non-hex character %q in hash %q", c, h)
		}
	}
}

func TestHashVoterIdentity_Deterministic(t *testing.T) {
	phone := "+2348012345678"
	salt := "app-salt"
	h1 := HashVoterIdentity(phone, salt)
	h2 := HashVoterIdentity(phone, salt)
	if h1 != h2 {
		t.Fatalf("hash is not deterministic: %q vs %q", h1, h2)
	}
}

func TestHashVoterIdentity_DifferentSaltsProduceDifferentHashes(t *testing.T) {
	phone := "+2348012345678"
	h1 := HashVoterIdentity(phone, "salt-a")
	h2 := HashVoterIdentity(phone, "salt-b")
	if h1 == h2 {
		t.Fatal("different salts should produce different hashes")
	}
}

func TestHashVoterIdentity_DifferentPhonesProduceDifferentHashes(t *testing.T) {
	salt := "app-salt"
	h1 := HashVoterIdentity("+2348012345678", salt)
	h2 := HashVoterIdentity("+2348099999999", salt)
	if h1 == h2 {
		t.Fatal("different phones should produce different hashes")
	}
}

// verify the hash value is computed as SHA-256("salt:normalizedPhone").
func TestHashVoterIdentity_KnownValue(t *testing.T) {
	phone := "08012345678"
	salt := "test-salt"
	// normalizePhone strips leading 0 → "2348012345678"
	input := fmt.Sprintf("%s:%s", salt, "2348012345678")
	sum := sha256.Sum256([]byte(input))
	want := fmt.Sprintf("%x", sum)

	got := HashVoterIdentity(phone, salt)
	if got != want {
		t.Fatalf("HashVoterIdentity(%q, %q) = %q, want %q", phone, salt, got, want)
	}
}

// phone formats that differ only in formatting must hash to the same value.
func TestHashVoterIdentity_PhoneNormalizationEquivalence(t *testing.T) {
	salt := "app-salt"
	variants := []string{
		"+2348012345678",
		"08012345678",
		"0801 234 5678",
		"0801-234-5678",
	}
	base := HashVoterIdentity(variants[0], salt)
	for _, v := range variants[1:] {
		got := HashVoterIdentity(v, salt)
		if got != base {
			t.Fatalf("HashVoterIdentity(%q) = %q, want %q (same as %q)", v, got, base, variants[0])
		}
	}
}
