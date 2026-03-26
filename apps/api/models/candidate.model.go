package models

import (
	"time"

	"github.com/google/uuid"
)

type PoliticalParty string

const (
	ActionPeoplesParty           PoliticalParty = "app"
	Accord                       PoliticalParty = "a"
	AlliedPeoplesMovement        PoliticalParty = "apm"
	AllProgressivesGrandAlliance PoliticalParty = "apga"
	AllProgressivesCongress      PoliticalParty = "apc"
	ActionDemocraticParty        PoliticalParty = "adp"
	AfricanDemocraticCongress    PoliticalParty = "adc"
	AfricanActionCongress        PoliticalParty = "aac"
	ActionAlliance               PoliticalParty = "aa"
	BootParty                    PoliticalParty = "bp"
	DemocraticLeadershipAlliance PoliticalParty = "dla"
	LabourParty                  PoliticalParty = "lp"
	NationalRescueMovement       PoliticalParty = "nrm"
	NewNigeriaPeoplesParty       PoliticalParty = "nnpp"
	PeoplesDemocraticParty       PoliticalParty = "pdp"
	PeoplesRedemptionParty       PoliticalParty = "prp"
	SocialDemocraticParty        PoliticalParty = "sdp"
	YoungProgressivesParty       PoliticalParty = "ypp"
	YouthParty                   PoliticalParty = "yp"
	ZenithLabourParty            PoliticalParty = "zlp"
)

type Candidate struct {
	ID           uuid.UUID      `json:"id"`
	Code         string         `json:"code"`
	Name         string         `json:"name"`
	Party        PoliticalParty `json:"party"`
	Bio          string         `json:"bio"`
	Achievements string         `json:"achievements"`
	PhotoURL     *string        `json:"photo_url"`
	IsActive     bool
	CreatedAt    time.Time
}
