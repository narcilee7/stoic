package main

import (
	"fmt"
	"os"

	"github.com/stoic/internal/config"
	"github.com/stoic/internal/philosopher"
	"github.com/stoic/internal/ui"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	philosopherManager, err := philosopher.NewManager(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing philosopher manager: %v\n", err)
		os.Exit(1)
	}

	app := ui.NewApp(cfg, philosopherManager)
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running app: %v\n", err)
		os.Exit(1)
	}
}
