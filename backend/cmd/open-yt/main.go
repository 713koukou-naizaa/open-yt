package main

import (
	"fmt"
	"os"

	"open-yt/internal/cli"
	"open-yt/internal/config"
)

func main() {
	configuration, configurationLoadError := config.Load()
	if configurationLoadError != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", configurationLoadError)
		os.Exit(1)
	}

	app := cli.NewApplication(configuration)

	appRunError := app.Run(os.Args[1:])
	if appRunError != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", appRunError)
		os.Exit(1)
	}
}
