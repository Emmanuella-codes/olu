package dtos

type CreateCandidateDTO struct {
	Code         string  `json:"code"         binding:"required"`
	Name         string  `json:"name"         binding:"required"`
	Party        string  `json:"party"        binding:"required"`
	Bio          string  `json:"bio"          binding:"required"`
	Achievements string  `json:"achievements" binding:"required"`
	PhotoURL     *string `json:"photo_url"`
}
