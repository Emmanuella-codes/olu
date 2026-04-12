package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/emmanuella-codes/olu/models"
	"github.com/emmanuella-codes/olu/services"
	"github.com/google/uuid"
)

// --- parseInboundVoteCode ---

func TestParseInboundVoteCode_FullFormat(t *testing.T) {
	code, err := parseInboundVoteCode("VOTE A1")
	if err != nil || code != "A1" {
		t.Fatalf("expected A1, got %q, err=%v", code, err)
	}
}

func TestParseInboundVoteCode_CodeOnly(t *testing.T) {
	code, err := parseInboundVoteCode("A1")
	if err != nil || code != "A1" {
		t.Fatalf("expected A1, got %q, err=%v", code, err)
	}
}

func TestParseInboundVoteCode_LowercaseInput(t *testing.T) {
	code, err := parseInboundVoteCode("vote a1")
	if err != nil || code != "A1" {
		t.Fatalf("expected A1, got %q, err=%v", code, err)
	}
}

func TestParseInboundVoteCode_TwoDigitCode(t *testing.T) {
	code, err := parseInboundVoteCode("VOTE B12")
	if err != nil || code != "B12" {
		t.Fatalf("expected B12, got %q, err=%v", code, err)
	}
}

func TestParseInboundVoteCode_InvalidFormat(t *testing.T) {
	cases := []string{"", "VOTE", "hello world", "VOTE 123", "!!!", "VOTE A"}
	for _, tc := range cases {
		_, err := parseInboundVoteCode(tc)
		if err == nil {
			t.Fatalf("expected error for input %q, got nil", tc)
		}
	}
}

// --- SMSWebhookHandler.InboundVote ---

func newSMSWebhookHandler() *SMSWebhookHandler {
	return NewSMSWebhookHandler(services.NewVoteService(nil, nil, nil), nil)
}

func TestInboundVote_MissingBody(t *testing.T) {
	r := newTestRouter()
	r.POST("/sms/inbound", newSMSWebhookHandler().InboundVote)

	w := performRequest(r, "POST", "/sms/inbound", nil)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body)
	}
}

func TestInboundVote_InvalidPhone(t *testing.T) {
	r := newTestRouter()
	r.POST("/sms/inbound", newSMSWebhookHandler().InboundVote)

	w := performRequest(r, "POST", "/sms/inbound", []byte(`{"from":"not-a-phone","text":"VOTE A1"}`))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body)
	}
}

func TestInboundVote_InvalidSMSFormat(t *testing.T) {
	r := newTestRouter()
	r.POST("/sms/inbound", newSMSWebhookHandler().InboundVote)

	w := performRequest(r, "POST", "/sms/inbound", []byte(`{"from":"08012345678","text":"HELLO"}`))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body)
	}
}

func TestInboundVote_InvalidCandidate(t *testing.T) {
	withFakeVoteRepo(t, &fakeVoteRepo{candidate: nil})

	r := newTestRouter()
	r.POST("/sms/inbound", newSMSWebhookHandler().InboundVote)

	w := performRequest(r, "POST", "/sms/inbound", []byte(`{"from":"08012345678","text":"VOTE A1"}`))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body)
	}
}

func TestInboundVote_AlreadyVoted(t *testing.T) {
	withFakeVoteRepo(t, &fakeVoteRepo{
		candidate: &models.Candidate{ID: uuid.New(), Code: "A1", Name: "Ada"},
		hasVoted:  true,
	})

	r := newTestRouter()
	r.POST("/sms/inbound", newSMSWebhookHandler().InboundVote)

	w := performRequest(r, "POST", "/sms/inbound", []byte(`{"from":"08012345678","text":"VOTE A1"}`))
	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d: %s", w.Code, w.Body)
	}
}

func TestInboundVote_Success(t *testing.T) {
	withFakeVoteRepo(t, &fakeVoteRepo{
		candidate: &models.Candidate{ID: uuid.New(), Code: "A1", Name: "Ada Obi"},
	})

	r := newTestRouter()
	r.POST("/sms/inbound", newSMSWebhookHandler().InboundVote)

	w := performRequest(r, "POST", "/sms/inbound", []byte(`{"from":"08012345678","text":"VOTE A1"}`))
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body)
	}

	var payload struct {
		Data struct {
			CandidateCode  string `json:"candidate_code"`
			CandidateName  string `json:"candidate_name"`
			ConfirmationID string `json:"confirmation_id"`
			Phone          string `json:"phone"`
			Channel        string `json:"channel"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if payload.Data.CandidateCode != "A1" || payload.Data.CandidateName != "Ada Obi" {
		t.Fatalf("unexpected candidate payload: %+v", payload.Data)
	}
	if payload.Data.Channel != string(models.SMSVoteChannel) {
		t.Fatalf("expected sms channel, got %q", payload.Data.Channel)
	}
	if payload.Data.Phone != "+234******5678" {
		t.Fatalf("expected masked phone +234******5678, got %q", payload.Data.Phone)
	}
	if payload.Data.ConfirmationID == "" {
		t.Fatal("expected confirmation id")
	}
}

// Verify that the SMS vote webhook accepts bare code format ("A1") in addition to "VOTE A1".
func TestInboundVote_BareCodeSuccess(t *testing.T) {
	withFakeVoteRepo(t, &fakeVoteRepo{
		candidate: &models.Candidate{ID: uuid.New(), Code: "A1", Name: "Ada Obi"},
	})

	r := newTestRouter()
	r.POST("/sms/inbound", newSMSWebhookHandler().InboundVote)

	w := performRequest(r, "POST", "/sms/inbound", []byte(`{"from":"08012345678","text":"A1"}`))
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for bare code, got %d: %s", w.Code, w.Body)
	}
}
