package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"location/pkg/config"
	"location/pkg/db/model"
)

//go:generate mockgen -destination=mocks/mock_Db.go -package=db -source=db.go

type DB interface {
	Get(id int) (*model.BranchOffice, error)
	Create(bo *model.BranchOffice) (*model.BranchOffice, error)
	All() []*model.BranchOffice
	ByLatLong(latitude float64, longitude float64) (*model.BranchOffice, error)
}

type BranchOffice struct {
	ID        uint    `json:"id" gorm:"primary_key"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"langitude"`
	Address   string  `json:"address"`
}

type db struct {
	db *gorm.DB
}

func New(conf *config.Config) *db {
	database, err := gorm.Open("sqlite3", conf.DBPath)

	if err != nil {
		panic("Failed to connect to database!" + err.Error())
	}

	database.AutoMigrate(&model.BranchOffice{})

	return &db{database}
}

func (d *db) Get(id int) (*model.BranchOffice, error) {
	bo := &model.BranchOffice{}
	if err := d.db.Where("id = ?", id).First(bo).Error; err != nil {
		return nil, err
	}
	return bo, nil
}

func (d *db) Create(bo *model.BranchOffice) (*model.BranchOffice, error) {
	d.db.Create(bo)
	return bo, nil
}

func (d *db) All() []*model.BranchOffice {
	var bos []*model.BranchOffice
	d.db.Find(&bos)
	return bos
}

func (d *db) ByLatLong(latitude float64, longitude float64) (*model.BranchOffice, error) {
	bo := &model.BranchOffice{}
	if err := d.db.Where("latitude = ? and longitude = ?", latitude, longitude).First(bo).Error; err != nil {
		return nil, err
	}
	return bo, nil
}
