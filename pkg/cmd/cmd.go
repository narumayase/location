package cmd

import (
	"location/api"
	"location/pkg/db"
	"math"
)

//go:generate mockgen -destination=mocks/mock_Cmd.go -package=cmd -source=cmd.go

type Command interface {
	Nearest(latitude float64, longitude float64) (*api.BranchOffice, error)
	Get(id int) (*api.BranchOffice, error)
	Create(branchOffice *api.BranchOffice) (*api.BranchOffice, error)
	Find(latitude float64, longitude float64) (*api.BranchOffice, error)
}

type command struct {
	db db.DB
}

func Build(db db.DB) Command {
	return &command{db}
}

func (c *command) Nearest(latitude float64, longitude float64) (nearest *api.BranchOffice, err error) {

	// primero busca si en la ubicación seleccionada no hay ya una sucursal
	if nearest, err = c.Find(latitude, longitude); err == nil {
		return nearest, nil
	}

	bos := c.all()
	mindist := 0.0

	for _, bo := range bos {

		catetoy := math.Abs(bo.Latitude - latitude)
		catetox := math.Abs(bo.Longitude - longitude) //módulo de (bo.Longitude - longitude)

		aux := (catetoy * catetoy) + (catetox * catetox)
		hypotenuse := math.Sqrt(aux) //raíz cuadrada de (catetoy al cuadrado + catetox al cuadrado)

		if mindist == 0.0 || mindist > hypotenuse {
			mindist = hypotenuse
			nearest = bo
		}
	}
	return nearest, nil
}

func (c *command) Get(id int) (*api.BranchOffice, error) {
	if bo, err := c.db.Get(id); err != nil {
		return nil, err
	} else {
		return toApi(bo), nil
	}
}

func (c *command) Create(bo *api.BranchOffice) (*api.BranchOffice, error) {
	var err error
	model := toModel(bo)
	if model, err = c.db.Create(model); err != nil {
		return nil, err
	}
	bo.ID = model.ID
	return bo, nil
}

func (c *command) all() []*api.BranchOffice {
	return toApis(c.db.All())
}

func (c *command) Find(latitude float64, longitude float64) (*api.BranchOffice, error) {
	if bo, err := c.db.ByLatLong(latitude, longitude); err != nil {
		return nil, err
	} else {
		return toApi(bo), nil
	}
}