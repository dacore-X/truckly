package main

import (
	"log"

	"github.com/dacore-x/truckly/config"
	"github.com/dacore-x/truckly/internal/app"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(cfg)
}
