package services

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/emmanuella-codes/olu/models"
	smsrepo "github.com/emmanuella-codes/olu/repositories/sms"
	smspkg "github.com/emmanuella-codes/olu/sms"
	"github.com/google/uuid"
)

type fakeSMSProvider struct {
	sent        []smspkg.Message
	sendErrByTo map[string]error
}

func (p *fakeSMSProvider) Send(_ context.Context, msg smspkg.Message) (*smspkg.Result, error) {
	p.sent = append(p.sent, msg)
	if err := p.sendErrByTo[msg.To]; err != nil {
		return nil, err
	}
	return &smspkg.Result{MessageID: "msg_001", Status: "sent"}, nil
}

func (p *fakeSMSProvider) Name() string {
	return "fake"
}

type fakeSMSRepository struct {
	enqueuedMSISDN string
	enqueuedBody   string
	batch          []models.SMSQueueItem
	batchErr       error
	lastBatchSize  int
	processed      []uuid.UUID
	processErr     error
	rejected       map[uuid.UUID]string
	rejectErr      error
}

func (r *fakeSMSRepository) EnqueueSMS(_ context.Context, msisdn, rawMessage string) error {
	r.enqueuedMSISDN = msisdn
	r.enqueuedBody = rawMessage
	return nil
}

func (r *fakeSMSRepository) GetPendingSMSBatch(_ context.Context, size int) ([]models.SMSQueueItem, error) {
	r.lastBatchSize = size
	if r.batchErr != nil {
		return nil, r.batchErr
	}
	return r.batch, nil
}

func (r *fakeSMSRepository) MarkSMSAsProcessed(_ context.Context, id uuid.UUID) error {
	r.processed = append(r.processed, id)
	return r.processErr
}

func (r *fakeSMSRepository) MarkSMSRejected(_ context.Context, id uuid.UUID, reason string) error {
	if r.rejected == nil {
		r.rejected = map[uuid.UUID]string{}
	}
	r.rejected[id] = reason
	return r.rejectErr
}

var _ smsrepo.SMSRepository = (*fakeSMSRepository)(nil)

func TestBuildOTPMessageIncludesCode(t *testing.T) {
	msg := buildOTPMessage("123456")

	if !strings.Contains(msg, "123456") {
		t.Fatalf("expected OTP message to include code, got %q", msg)
	}
}

func TestBuildVoteConfirmationMessageTrimsAndTruncatesID(t *testing.T) {
	msg := buildVoteConfirmationMessage(" Ada ", "123e4567-e89b-12d3-a456-426614174000")
	want := "Vote confirmed for Ada. Confirmation ID: 123E4567."

	if msg != want {
		t.Fatalf("expected %q, got %q", want, msg)
	}
}

func TestBuildVoteConfirmationMessageShortID(t *testing.T) {
	msg := buildVoteConfirmationMessage("Ben", "abc")
	want := "Vote confirmed for Ben. Confirmation ID: ABC."

	if msg != want {
		t.Fatalf("expected %q, got %q", want, msg)
	}
}

func TestBuildVoteRejectionMessageFallsBackWhenReasonIsEmpty(t *testing.T) {
	msg := buildVoteRejectionMessage("   ")
	want := "We could not process your vote."

	if msg != want {
		t.Fatalf("expected %q, got %q", want, msg)
	}
}

func TestBuildVoteRejectionMessageReturnsNonEmptyReason(t *testing.T) {
	reason := "already voted"
	msg := buildVoteRejectionMessage(reason)

	if msg != reason {
		t.Fatalf("expected reason %q returned as-is, got %q", reason, msg)
	}
}

func TestSMSServiceSendOTPSendsThroughProvider(t *testing.T) {
	provider := &fakeSMSProvider{sendErrByTo: map[string]error{}}
	svc := &SMSService{provider: provider, from: "OLU"}

	if err := svc.SendOTP(context.Background(), "09090903080", "123456"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(provider.sent) != 1 {
		t.Fatalf("expected 1 sent message, got %d", len(provider.sent))
	}

	msg := provider.sent[0]
	if msg.To != "09090903080" || msg.From != "OLU" || msg.Type != "otp" {
		t.Fatalf("unexpected sent message: %+v", msg)
	}
	if !strings.Contains(msg.Body, "123456") {
		t.Fatalf("expected OTP body to include code, got %q", msg.Body)
	}
}

func TestSMSServiceSendOTPNilProviderReturnsError(t *testing.T) {
	svc := &SMSService{provider: nil, from: "OLU"}

	err := svc.SendOTP(context.Background(), "09090903080", "123456")
	if err == nil {
		t.Fatal("expected error with nil provider, got nil")
	}
}

func TestSMSServiceQueueVoteConfirmationEnqueuesMessage(t *testing.T) {
	repo := &fakeSMSRepository{}
	svc := &SMSService{queueRepo: repo}

	if err := svc.QueueVoteConfirmation(context.Background(), "09090903080", "Ada", "123e4567-e89b"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repo.enqueuedMSISDN != "09090903080" {
		t.Fatalf("expected enqueued phone 09090903080, got %q", repo.enqueuedMSISDN)
	}
	if !strings.Contains(repo.enqueuedBody, "Vote confirmed for Ada") {
		t.Fatalf("expected confirmation body, got %q", repo.enqueuedBody)
	}
}

func TestSMSServiceQueueVoteRejectionEnqueuesMessage(t *testing.T) {
	repo := &fakeSMSRepository{}
	svc := &SMSService{queueRepo: repo}

	if err := svc.QueueVoteRejection(context.Background(), "09090903080", "already voted"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repo.enqueuedMSISDN != "09090903080" {
		t.Fatalf("expected enqueued phone 09090903080, got %q", repo.enqueuedMSISDN)
	}
	if repo.enqueuedBody != "already voted" {
		t.Fatalf("expected rejection body %q, got %q", "already voted", repo.enqueuedBody)
	}
}

func TestSMSServiceDeliverPendingBatchMarksProcessedAndRejected(t *testing.T) {
	okID := uuid.New()
	failID := uuid.New()
	repo := &fakeSMSRepository{batch: []models.SMSQueueItem{
		{ID: okID, MSISDN: "09000000001", RawMessage: "ok"},
		{ID: failID, MSISDN: "09000000002", RawMessage: "fail"},
	}}
	provider := &fakeSMSProvider{sendErrByTo: map[string]error{
		"09000000002": errors.New("provider failed"),
	}}
	svc := &SMSService{provider: provider, queueRepo: repo, from: "OLU"}

	if err := svc.DeliverPendingBatch(context.Background(), 20); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(repo.processed) != 1 || repo.processed[0] != okID {
		t.Fatalf("expected processed id %s, got %+v", okID, repo.processed)
	}
	if repo.rejected[failID] != "provider failed" {
		t.Fatalf("expected rejected reason provider failed, got %q", repo.rejected[failID])
	}
}

func TestSMSServiceDeliverPendingBatchNilProviderReturnsError(t *testing.T) {
	svc := &SMSService{provider: nil, queueRepo: &fakeSMSRepository{}}

	err := svc.DeliverPendingBatch(context.Background(), 10)
	if err == nil {
		t.Fatal("expected error with nil provider, got nil")
	}
}

func TestSMSServiceDeliverPendingBatchNilRepoReturnsError(t *testing.T) {
	provider := &fakeSMSProvider{sendErrByTo: map[string]error{}}
	svc := &SMSService{provider: provider, queueRepo: nil}

	err := svc.DeliverPendingBatch(context.Background(), 10)
	if err == nil {
		t.Fatal("expected error with nil queue repo, got nil")
	}
}

func TestSMSServiceDeliverPendingBatchZeroBatchSizeDefaultsTwenty(t *testing.T) {
	repo := &fakeSMSRepository{}
	provider := &fakeSMSProvider{sendErrByTo: map[string]error{}}
	svc := &SMSService{provider: provider, queueRepo: repo, from: "OLU"}

	if err := svc.DeliverPendingBatch(context.Background(), 0); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if repo.lastBatchSize != 20 {
		t.Fatalf("expected default batch size 20, got %d", repo.lastBatchSize)
	}
}

func TestSMSServiceDeliverPendingBatchEmptyBatchIsNoop(t *testing.T) {
	repo := &fakeSMSRepository{batch: []models.SMSQueueItem{}}
	provider := &fakeSMSProvider{sendErrByTo: map[string]error{}}
	svc := &SMSService{provider: provider, queueRepo: repo, from: "OLU"}

	if err := svc.DeliverPendingBatch(context.Background(), 10); err != nil {
		t.Fatalf("expected no error for empty batch, got %v", err)
	}
	if len(repo.processed) != 0 {
		t.Fatalf("expected no processed items, got %+v", repo.processed)
	}
}

func TestSMSServiceDeliverPendingBatchFetchError(t *testing.T) {
	repo := &fakeSMSRepository{batchErr: errors.New("query failed")}
	provider := &fakeSMSProvider{sendErrByTo: map[string]error{}}
	svc := &SMSService{provider: provider, queueRepo: repo, from: "OLU"}

	err := svc.DeliverPendingBatch(context.Background(), 10)
	if err == nil || !strings.Contains(err.Error(), "sms: fetch pending batch") {
		t.Fatalf("expected fetch pending batch error, got %v", err)
	}
}

func TestSMSServiceDeliverPendingBatchMarkProcessedErrorDoesNotFailBatch(t *testing.T) {
	id := uuid.New()
	repo := &fakeSMSRepository{
		batch:      []models.SMSQueueItem{{ID: id, MSISDN: "09000000001", RawMessage: "ok"}},
		processErr: errors.New("mark processed failed"),
	}
	provider := &fakeSMSProvider{sendErrByTo: map[string]error{}}
	svc := &SMSService{provider: provider, queueRepo: repo, from: "OLU"}

	if err := svc.DeliverPendingBatch(context.Background(), 10); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(repo.processed) != 1 || repo.processed[0] != id {
		t.Fatalf("expected processed attempt for %s, got %+v", id, repo.processed)
	}
}

func TestSMSServiceDeliverPendingBatchMarkRejectedErrorDoesNotFailBatch(t *testing.T) {
	id := uuid.New()
	repo := &fakeSMSRepository{
		batch:     []models.SMSQueueItem{{ID: id, MSISDN: "09000000002", RawMessage: "fail"}},
		rejectErr: errors.New("mark rejected failed"),
	}
	provider := &fakeSMSProvider{sendErrByTo: map[string]error{
		"09000000002": errors.New("provider failed"),
	}}
	svc := &SMSService{provider: provider, queueRepo: repo, from: "OLU"}

	if err := svc.DeliverPendingBatch(context.Background(), 10); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if repo.rejected[id] != "provider failed" {
		t.Fatalf("expected rejected attempt with provider failure, got %q", repo.rejected[id])
	}
}
