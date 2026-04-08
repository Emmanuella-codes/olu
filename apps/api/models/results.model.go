package models

import "time"

type Results struct {
	Tally      []TallyRow `json:"tally"`
	TotalVotes int64      `json:"total_votes"`
	IsTie      bool       `json:"is_tie"`
	Leaders    []TallyRow `json:"leaders"`
	CachedAt   time.Time  `json:"cached_at"`
}
