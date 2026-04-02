package workers

import (
	"context"
	"errors"
	"time"

	"github.com/emmanuella-codes/olu/services"
	"github.com/rs/zerolog/log"
)

type SMSWorker struct {
	smsService *services.SMSService
}

func NewSMSWorker(smsService *services.SMSService) *SMSWorker {
	return &SMSWorker{smsService: smsService}
}

func (w *SMSWorker) RunQueueWorker(ctx context.Context, batchSize int, pollInterval time.Duration) {
	if w.smsService == nil {
		return
	}
	if batchSize <= 0 {
		batchSize = 20
	}
	if pollInterval <= 0 {
		pollInterval = 2 * time.Second
	}

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.deliverPendingBatch(ctx, batchSize)
		}
	}
}

func (w *SMSWorker) deliverPendingBatch(ctx context.Context, batchSize int) {
	err := w.smsService.DeliverPendingBatch(ctx, batchSize)
	if err != nil {
		if ctx.Err() != nil || errors.Is(err, context.Canceled) {
			return
		}
		log.Warn().Err(err).Msg("sms worker batch failed")
	}
}
