package model

type BranchOffice struct {
	ID        uint    `json:"id" gorm:"primary_key"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"langitude"`
	Address   string  `json:"address"`
}
