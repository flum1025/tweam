package main

import (
	"flag"
	"log"

	"github.com/flum1025/tweam/internal/config"
	"github.com/flum1025/tweam/internal/scheduler"
)

func main() {
	configPath := flag.String("config", "", "config path")
	workerEndpoint := flag.String("worker_endpoint", "", "worker endpoint")
	flag.Parse()

	config, err := config.NewConfig(*configPath)
	if err != nil {
		log.Fatalf("failed to get config: %v", err)

		return
	}

	scheduler := scheduler.NewScheduler(config)
	if err = scheduler.Run(*workerEndpoint); err != nil {
		log.Fatalf("failed to run scheduler: %v", err)
	}
}
