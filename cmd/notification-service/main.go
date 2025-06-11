package main

import (
	"flag"
	"fmt"
	config "github.com/nogavadu/notification-service/internal/config"
	"github.com/nogavadu/notification-service/internal/consumer"
	emailService "github.com/nogavadu/notification-service/internal/service/email"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	log.Info("Initializing application")

	var configPath string
	flag.StringVar(&configPath, "config", "", "Path to config file")
	flag.Parse()

	if err := os.Setenv("CONFIG_PATH", configPath); err != nil {
		log.Error("failed to set CONFIG_PATH env var", "err", err)
		os.Exit(1)
	}
	log.Info(fmt.Sprintf("Config path: %s", configPath))

	cfg, err := config.New()
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	log.Info("Starting notification service consumer")

	keepRunning := true

	emailServ := emailService.New(log)

	c, err := consumer.New(cfg.Brokers, cfg.Topics, cfg.Group, emailServ)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	wg := &sync.WaitGroup{}
	go func() {
		defer wg.Done()
		wg.Add(1)

		if err = c.Start(); err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		select {
		case sig := <-signals:
			log.Info("terminating: signal", sig)
			c.Stop()
			keepRunning = false
		}
	}

	wg.Wait()
}
