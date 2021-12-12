package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"location/pkg/cmd"
	"location/pkg/config"
	"location/pkg/db"
	"location/pkg/routes"
	"os"
)

func main() {

	var env string
	if env = os.Getenv("LOCATION_ENVIRONMENT"); env == "" {
		fmt.Println("LOCATION_ENVIRONMENT is not set, assuming non locally.")
	}

	conf := config.New(env)

	g := gin.Default()

	d := db.New(conf)

	c := cmd.Build(d)

	routes.AddHandler(g, c)

	if err := g.Run(); err != nil {
		fmt.Printf("there was an error executing command: %s.\n", err.Error())
		os.Exit(1)
	}
}
