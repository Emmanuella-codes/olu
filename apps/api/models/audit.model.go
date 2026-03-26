package models

type AuditEntry struct {
	VoterHash     string
	CandidateCode string
	Channel       string
	Status        string
	IPAddress     string
	UserAgent     string
}
