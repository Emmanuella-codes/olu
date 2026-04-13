package models

import "github.com/google/uuid"

type TallyRow struct {
	CandidateID uuid.UUID      `json:"candidate_id"`
	Code        string         `json:"code"`
	Name        string         `json:"name"`
	Party       PoliticalParty `json:"party"`
	VoteCount   int            `json:"vote_count"`
}
