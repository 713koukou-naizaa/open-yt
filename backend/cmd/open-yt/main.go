package main

import (
	"log"
	"os"

	"open-yt/internal/cli"
	"open-yt/internal/config"
)

func main() {
	cfg := config.Default()
	app := cli.NewApp(cfg)
	if err := app.Run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}