package models

import (
	"time"

	"github.com/google/uuid"
)

type VoteChannel string

const (
	SMSVoteChannel VoteChannel = "sms"
	WebVoteChannel VoteChannel = "web"
)

type Vote struct {
	ID          uuid.UUID
	CandidateID uuid.UUID
	VoterHash   string
	Channel     VoteChannel
	CreatedAt   time.Time
}
