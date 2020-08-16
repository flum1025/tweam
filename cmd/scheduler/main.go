package main

import (
	"flag"
	"log"

	"github.com/flum1025/tweam/internal/config"
	"github.com/flum1025/tweam/internal/scheduler"
)

func main() {
	configPath := flag.String("config", "", "config path")
	flag.Parse()

	config, err := config.NewConfig(*configPath)
	if err != nil {
		log.Fatalf("failed to get config: %v", err)

		return
	}

	scheduler, err := scheduler.NewScheduler(config)
	if err != nil {
		log.Fatalf("failed to get scheduler: %v", err)

		return
	}

	scheduler.Run()
}
