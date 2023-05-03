package main

import (
	"log"

	webapp "github.com/begenov/Forum-Go/internal/web-app"

	"github.com/begenov/Forum-Go/config"
)

const path = "./config/config.json"

func main() {
	cfg, err := config.NewConfig(path)
	if err != nil {
		log.Fatalln(err)
	}

	app := webapp.NewApp(*cfg)

	if err := app.Run(); err != nil {
		log.Fatalln(err)
		return
	}
}
