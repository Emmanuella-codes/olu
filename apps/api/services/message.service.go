package services

import (
	"context"
	"fmt"
	"strings"

	smsrepo "github.com/emmanuella-codes/olu/repositories/sms"
	"github.com/emmanuella-codes/olu/sms"
	"github.com/rs/zerolog/log"
)

type SMSService struct {
	provider  sms.Provider
	queueRepo smsrepo.SMSRepository
	from      string
}

func NewSMSService(provider sms.Provider, from string) *SMSService {
	return &SMSService{
		provider:  provider,
		queueRepo: smsrepo.SMSRepo,
		from:      from,
	}
}

func (s *SMSService) Provider() sms.Provider {
	return s.provider
}

func (s *SMSService) SendOTP(ctx context.Context, phone, code string) error {
	return s.send(ctx, sms.Message{
		To:   phone,
		From: s.from,
		Body: buildOTPMessage(code),
		Type: "otp",
	})
}

func (s *SMSService) QueueVoteConfirmation(ctx context.Context, phone, candidateName, confirmationID string) error {
	return s.enqueue(ctx, phone, buildVoteConfirmationMessage(candidateName, confirmationID))
}

func (s *SMSService) QueueVoteRejection(ctx context.Context, phone, reason string) error {
	return s.enqueue(ctx, phone, buildVoteRejectionMessage(reason))
}

func (s *SMSService) DeliverPendingBatch(ctx context.Context, batchSize int) error {
	if s.provider == nil {
		return fmt.Errorf("sms: provider is not configured")
	}
	if s.queueRepo == nil {
		return fmt.Errorf("sms: queue repository is not initialized")
	}
	if batchSize <= 0 {
		batchSize = 20
	}

	items, err := s.queueRepo.GetPendingSMSBatch(ctx, batchSize)
	if err != nil {
		return fmt.Errorf("sms: fetch pending batch: %w", err)
	}

	for _, item := range items {
		_, err := s.provider.Send(ctx, sms.Message{
			To:   item.MSISDN,
			From: s.from,
			Body: item.RawMessage,
			Type: "queued",
		})
		if err != nil {
			if markErr := s.queueRepo.MarkSMSRejected(ctx, item.ID, err.Error()); markErr != nil {
				log.Warn().Err(markErr).Str("id", item.ID.String()).Msg("sms: mark rejected failed")
			}
			log.Warn().Err(err).Str("to", item.MSISDN).Msg("sms: delivery failed")
			continue
		}

		if err := s.queueRepo.MarkSMSAsProcessed(ctx, item.ID); err != nil {
			log.Warn().Err(err).Str("id", item.ID.String()).Msg("sms: mark processed failed")
			continue
		}
	}

	return nil
}

func (s *SMSService) send(ctx context.Context, msg sms.Message) error {
	if s.provider == nil {
		return fmt.Errorf("sms: provider is not configured")
	}

	_, err := s.provider.Send(ctx, msg)
	return err
}

func (s *SMSService) enqueue(ctx context.Context, phone, body string) error {
	if s.queueRepo == nil {
		return fmt.Errorf("sms: queue repository is not initialized")
	}

	return s.queueRepo.EnqueueSMS(ctx, phone, body)
}

func buildOTPMessage(code string) string {
	return fmt.Sprintf("Your Olu voting code is %s. It expires in 10 minutes.", code)
}

func buildVoteConfirmationMessage(candidateName, confirmationID string) string {
	confirmationID = strings.TrimSpace(confirmationID)
	if len(confirmationID) > 8 {
		confirmationID = confirmationID[:8]
	}

	return fmt.Sprintf(
		"Vote confirmed for %s. Confirmation ID: %s.",
		strings.TrimSpace(candidateName),
		strings.ToUpper(confirmationID),
	)
}

func buildVoteRejectionMessage(reason string) string {
	reason = strings.TrimSpace(reason)
	if reason == "" {
		reason = "We could not process your vote."
	}
	return reason
}
