package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"location/api"
	"location/pkg/cmd"
	"net/http"
	"strconv"
)

type Handler interface {
	Get(c *gin.Context)
	Create(c *gin.Context)
	Nearest(c *gin.Context)
}

type handler struct {
	cmd cmd.Command
}

func AddHandler(e *gin.Engine, cmd cmd.Command) {

	h := &handler{cmd}

	rootPath := "branch-offices"

	e.GET(fmt.Sprintf("/%s/branch-office/:id", rootPath), h.Get)
	e.POST(fmt.Sprintf("/%s/branch-office", rootPath), h.Create)
	e.GET(fmt.Sprintf("/%s/nearest", rootPath), h.Nearest)
}

func (h *handler) Get(c *gin.Context) {
	var bo *api.BranchOffice

	if id, err := strconv.Atoi(c.Param("id")); err == nil {
		if bo, err = h.cmd.Get(id); err == nil {
			c.JSON(http.StatusOK, bo)
		} else {
			fmt.Println("there was an error: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	} else {
		fmt.Println("there was an error: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func (h *handler) Create(c *gin.Context) {
	var bo *api.BranchOffice

	if err := c.ShouldBindJSON(&bo); err != nil {
		fmt.Println("there was an error: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if h.alreadyExists(bo) {
		fmt.Println("branchoffice already exists")
		c.JSON(http.StatusBadRequest, gin.H{"error": "branchoffice already exists"})
		return
	}
	h.create(c, bo)
}

func (h *handler) create(c *gin.Context, bo *api.BranchOffice) {
	var err error
	if bo, err = h.cmd.Create(bo); err == nil {
		c.JSON(http.StatusOK, bo)
	} else {
		fmt.Println("there was an error: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func (h *handler) alreadyExists(bo *api.BranchOffice) bool {
	if _, err := h.cmd.Find(bo.Latitude, bo.Longitude); err == nil {
		return true
	}
	return false
}

func (h *handler) Nearest(c *gin.Context) {

	bo := &api.BranchOffice{}

	if longitude, err := strconv.ParseFloat(c.Query("longitude"), 64); err != nil {
		fmt.Println("there was an error: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Longitude is not a float64!"})
		return
	} else {
		bo.Longitude = longitude
	}

	if latitude, err := strconv.ParseFloat(c.Query("latitude"), 64); err != nil {
		fmt.Println("there was an error: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Latitude is not a float64!"})
		return
	} else {
		bo.Latitude = latitude
	}

	var err error
	if bo, err = h.cmd.Nearest(bo.Latitude, bo.Longitude); err != nil {
		fmt.Println("there was an error: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bo)
}
