package models

import "time"

type Results struct {
	Tally      []TallyRow `json:"tally"`
	TotalVotes int64      `json:"total_votes"`
	CachedAt   time.Time  `json:"cached_at"`
}
