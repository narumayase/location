package api

type BranchOffice struct {
	ID        uint   `json:"id" binding:"required"`
	Longitude float64    `json:"longitude" binding:"required"`
	Latitude  float64    `json:"latitude" binding:"required"`
	Address   string `json:"address" binding:"required"`
}
