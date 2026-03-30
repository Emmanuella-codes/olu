package models

type Stats struct {
	TotalVotes      int64 `json:"total_votes"`
	WebVotes        int64 `json:"web_votes"`
	SMSVotes        int64 `json:"sms_votes"`
	PendingSMS      int64 `json:"pending_sms"`
	TotalCandidates int64 `json:"total_candidates"`
}
