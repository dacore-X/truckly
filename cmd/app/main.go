package main

import (
	"log"

	"github.com/dacore-x/truckly/config"
	"github.com/dacore-x/truckly/internal/app"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	app.Run(cfg)
}
