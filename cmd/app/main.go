package main

import (
	"fourth-exam/post-service-clean-arch/config"
	"fourth-exam/post-service-clean-arch/internal/app"
	"log"
)

func main() {
	// Configuration 
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run 
	app.Run(cfg)
}