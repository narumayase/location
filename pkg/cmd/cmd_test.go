package cmd

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"location/api"
	mocks2 "location/pkg/db/mocks"
	"location/pkg/db/model"
	"testing"
)

func TestCreate_success(t *testing.T) {
	command, mocks := setUp(t)
	defer mocks.GoMockController.Controller.Finish()

	bsasm := newBsAsModel()
	bsasm.ID = 0

	bsasmWithId := newBsAsModel()

	// c.db.Create(toModel(bo))
	mocks.db.EXPECT().Create(bsasm).Return(bsasmWithId, nil).Times(1)

	bsas := newBsAsApi()
	bsas.ID = 0

	response, err := command.Create(bsas)

	assert.Equal(t, response, newBsAsApi(), "create branchoffice successfully - returns same branchoffice with id")
	assert.NoError(t, err, "create branchoffice successfully - doesn't return an error")
}

func TestCreate_fail(t *testing.T) {
	command, mocks := setUp(t)
	defer mocks.GoMockController.Controller.Finish()

	bsasm := newBsAsModel()
	bsasm.ID = 0

	// c.db.Create(toModel(bo))
	mocks.db.EXPECT().Create(bsasm).Return(nil, errors.New("db error")).Times(1)

	bsas := newBsAsApi()
	bsas.ID = 0

	response, err := command.Create(bsas)

	assert.Nil(t, response, "create branchoffice fail - returns nil")
	assert.EqualError(t, err, "db error", "create branchoffice fail - returns 'db error'")
}

func TestNearest_success(t *testing.T) {
	command, mocks := setUp(t)
	defer mocks.GoMockController.Controller.Finish()

	bsas := newBsAsApi()
	ushuaia := newUshuaiaApi() // coordenadas target

	//c.db.ByLatLong(latitude, longitude)
	mocks.db.EXPECT().ByLatLong(ushuaia.Latitude, ushuaia.Longitude).Return(nil, errors.New("not found")).Times(1)

	//c.db.All()
	mocks.db.EXPECT().All().Return([]*model.BranchOffice{newBsAsModel(), newCordobaModel()}).Times(1)

	response, err := command.Nearest(ushuaia.Latitude, ushuaia.Longitude)

	assert.Equal(t, response, bsas, "get nearest branchoffice beetween three positions successfully - nearest branchoffice between cordoba and bsas to ushuaia position should be bsas")
	assert.NoError(t, err, "get nearest branchoffice beetween three positions successfully - doesn't return an error")
}

func TestNearest_hasBranchOffice_success(t *testing.T) {
	command, mocks := setUp(t)
	defer mocks.GoMockController.Controller.Finish()

	ushuaia := newUshuaiaApi() // coordenadas target

	//c.db.ByLatLong(latitude, longitude)
	mocks.db.EXPECT().ByLatLong(ushuaia.Latitude, ushuaia.Longitude).Return(newUshuaiaModel(), nil).Times(1)

	response, err := command.Nearest(ushuaia.Latitude, ushuaia.Longitude)

	assert.Equal(t, response, ushuaia, "get nearest branchoffice with an existing branchoffice successfully - nearest branchoffice to ushuaia position should be ushuaia")
	assert.NoError(t, err, "get nearest branchoffice with an existing branchoffice successfully - doesn't return an error")
}

type mocks struct {
	GoMockController
	db *mocks2.MockDB
}

func (builder *mocks) build() Command {
	return Build(builder.db)
}

func setUp(t *testing.T) (Command, *mocks) {
	ctrl := gomock.NewController(t)

	mocks := &mocks{
		GoMockController: GoMockController{ctrl},
		db:               mocks2.NewMockDB(ctrl),
	}
	return mocks.build(), mocks
}

type GoMockController struct {
	Controller *gomock.Controller
}

func newBsAsModel() *model.BranchOffice {
	return &model.BranchOffice{
		ID:        1,
		Longitude: -58.45678,
		Latitude:  -34.12345,
		Address:   "bsas address",
	}
}

func newBsAsApi() *api.BranchOffice {
	return &api.BranchOffice{
		ID:        1,
		Longitude: -58.45678,
		Latitude:  -34.12345,
		Address:   "bsas address",
	}
}

func newCordobaModel() *model.BranchOffice {
	return &model.BranchOffice{
		ID:        2,
		Longitude: -64.18105,
		Latitude:  -31.4135,
		Address:   "cordoba address",
	}
}

func newUshuaiaModel() *model.BranchOffice {
	return &model.BranchOffice{
		ID:        3,
		Longitude: -68.31591,
		Latitude:  -54.81084,
		Address:   "ushuaia address",
	}
}

func newUshuaiaApi() *api.BranchOffice {
	return &api.BranchOffice{
		ID:        3,
		Longitude: -68.31591,
		Latitude:  -54.81084,
		Address:   "ushuaia address",
	}
}

func newCordobaApi() *api.BranchOffice {
	return &api.BranchOffice{
		ID:        2,
		Longitude: -64.18105,
		Latitude:  -31.4135,
		Address:   "cordoba address",
	}
}
