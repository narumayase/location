package routes

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
	"location/api"
	mocks2 "location/pkg/cmd/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreate_success(t *testing.T) {
	handler, mocks := setUp(t)
	defer mocks.GoMockController.Controller.Finish()

	bo := newBsas()
	boWithoutId := newBsasWithoutId()

	//h.cmd.Find(bo.Latitude, bo.Longitude)
	mocks.command.EXPECT().Find(bo.Latitude, bo.Longitude).Return(nil, errors.New("not found"))

	//h.cmd.Create(bo)
	mocks.command.EXPECT().Create(boWithoutId).Return(bo, nil)

	body, _ := json.Marshal(boWithoutId)
	request, response := newRequest("POST", "/branch-offices/branch-office", bytes.NewBuffer(body))

	handler.ServeHTTP(response, request)

	data, err := getApi(response)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, data, bo, "create successfully - returns api with id")
	assert.Equal(t, response.Code, 200, "create successfully - returns response code 200")
}

func TestCreate_server_fail(t *testing.T) {
	handler, mocks := setUp(t)
	defer mocks.GoMockController.Controller.Finish()

	boWithoutId := newBsasWithoutId()

	//h.cmd.Find(bo.Latitude, bo.Longitude)
	mocks.command.EXPECT().Find(boWithoutId.Latitude, boWithoutId.Longitude).Return(nil, errors.New("not found"))

	//h.cmd.Create(bo)
	mocks.command.EXPECT().Create(boWithoutId).Return(nil, errors.New("db error"))

	body, _ := json.Marshal(boWithoutId)
	request, response := newRequest("POST", "/branch-offices/branch-office", bytes.NewBuffer(body))

	handler.ServeHTTP(response, request)

	data, err := getError(response)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, data.Error, "db error", "create with db error - returns 'db error'")
	assert.Equal(t, response.Code, 500, "create with db error - returns response code 500")
}

func TestCreate_alreadyExists(t *testing.T) {
	handler, mocks := setUp(t)
	defer mocks.GoMockController.Controller.Finish()

	boWithoutId := newBsasWithoutId()

	//h.cmd.Find(bo.Latitude, bo.Longitude)
	mocks.command.EXPECT().Find(boWithoutId.Latitude, boWithoutId.Longitude).Return(nil, nil)

	body, _ := json.Marshal(boWithoutId)
	request, response := newRequest("POST", "/branch-offices/branch-office", bytes.NewBuffer(body))

	handler.ServeHTTP(response, request)

	data, err := getError(response)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, data.Error, "branchoffice already exists", "create existing branchoffice  - returns 'branch office already exists'")
	assert.Equal(t, response.Code, 400, "create existing branchoffice - returns response code 400")
}

func TestNearest_success(t *testing.T) {
	handler, mocks := setUp(t)
	defer mocks.GoMockController.Controller.Finish()

	ushuaia := newUshuaia()
	bsas := newBsas()

	//h.cmd.Nearest(bo.Latitude, bo.Longitude)
	mocks.command.EXPECT().Nearest(ushuaia.Latitude, ushuaia.Longitude).Return(bsas, nil)

	// latitud y longitud de ushuaia
	request, response := newRequest("GET", "/branch-offices/nearest?latitude=-54.81084&longitude=-68.31591", nil)

	handler.ServeHTTP(response, request)

	data, err := getApi(response)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, data, bsas, "get nearest branchoffice - nearest branchoffice to ushuaia returns bsas")
	assert.Equal(t, response.Code, 200, "get nearest branchoffice - returns response code 200")
}

func TestNearest_badLatitudFormat(t *testing.T) {
	handler, mocks := setUp(t)
	defer mocks.GoMockController.Controller.Finish()

	request, response := newRequest("GET", "/branch-offices/nearest?latitude=sarasa&longitude=-68.31591", nil)

	handler.ServeHTTP(response, request)

	data, err := getError(response)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, data.Error, "Latitude is not a float64!", "get nearest with bad latitude format - returns 'Latitude is not a float64!'")
	assert.Equal(t, response.Code, 400, "get nearest with bad latitude format - returns response code 400")
}

func TestNearest_badLongitudeFormat(t *testing.T) {
	handler, mocks := setUp(t)
	defer mocks.GoMockController.Controller.Finish()

	request, response := newRequest("GET", "/branch-offices/nearest?latitude=-54.81084&longitude=sarasa", nil)

	handler.ServeHTTP(response, request)

	data, err := getError(response)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, data.Error, "Longitude is not a float64!", "get nearest with bad longitude format - returns 'Longitude is not a float64!'")
	assert.Equal(t, response.Code, 400, "get nearest with bad longitude format - returns response code 400")
}

func TestNearest_server_error(t *testing.T) {
	handler, mocks := setUp(t)
	defer mocks.GoMockController.Controller.Finish()

	ushuaia := newUshuaia()

	//h.cmd.Nearest(bo.Latitude, bo.Longitude)
	mocks.command.EXPECT().Nearest(ushuaia.Latitude, ushuaia.Longitude).Return(nil, errors.New("server error"))

	// latitud y longitud de ushuaia
	request, response := newRequest("GET", "/branch-offices/nearest?latitude=-54.81084&longitude=-68.31591", nil)

	handler.ServeHTTP(response, request)

	data, err := getError(response)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, data.Error, "server error", "get nearest with server error - returns the same error, 'server error' in this case")
	assert.Equal(t, response.Code, 500, "get nearest with server error - returns response code 500")
}

func getApi(response *httptest.ResponseRecorder) (*api.BranchOffice, error) {
	data := &api.BranchOffice{}
	if err := json.Unmarshal(response.Body.Bytes(), data); err != nil {
		return nil, err
	}
	return data, nil
}

func getError(response *httptest.ResponseRecorder) (*api.Error, error) {
	data := &api.Error{}
	if err := json.Unmarshal(response.Body.Bytes(), data); err != nil {
		return nil, err
	}
	return data, nil
}

func newRequest(method, url string, body io.Reader) (*http.Request, *httptest.ResponseRecorder) {
	request, _ := http.NewRequest(method, url, body)
	return request, httptest.NewRecorder()
}

func newBsas() *api.BranchOffice {
	return &api.BranchOffice{
		ID:        1,
		Longitude: -58.45678,
		Latitude:  -34.12345,
		Address:   "bsas address",
	}
}

func newBsasWithoutId() *api.BranchOffice {
	return &api.BranchOffice{
		Longitude: -58.45678,
		Latitude:  -34.12345,
		Address:   "bsas address",
	}
}

func newUshuaia() *api.BranchOffice {
	return &api.BranchOffice{
		ID:        3,
		Longitude: -68.31591,
		Latitude:  -54.81084,
		Address:   "ushuaia address",
	}
}

type mocks struct {
	GoMockController
	command *mocks2.MockCommand
}

type GoMockController struct {
	Controller *gomock.Controller
}

func (builder *mocks) build() *gin.Engine {
	r := gin.Default()
	AddHandler(r, builder.command)
	return r
}

func setUp(t *testing.T) (*gin.Engine, *mocks) {
	ctrl := gomock.NewController(t)

	mocks := &mocks{
		GoMockController: GoMockController{ctrl},
		command:          mocks2.NewMockCommand(ctrl),
	}
	return mocks.build(), mocks
}
