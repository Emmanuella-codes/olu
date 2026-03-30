package dtos

type UpdateCandidateDTO struct {
	Code         *string `json:"code"`
	Name         *string `json:"name"`
	Party        *string `json:"party"`
	Bio          *string `json:"bio"`
	Achievements *string `json:"achievements"`
	PhotoURL     *string `json:"photo_url"`
	IsActive     *bool   `json:"is_active"`
}
