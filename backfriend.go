package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/nawa/back-friend/config"
	"github.com/nawa/back-friend/rest"
	"github.com/nawa/back-friend/storage/postgres"
)

func main() {
	configFile := flag.String("config", "config/env/dev.yml", "configuration file")
	flag.Parse()

	cfg, err := config.FromFile(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	storage, err := postgres.NewPostgresStorage(*cfg.SQLDb)
	if err != nil {
		log.Fatal(err)
	}

	exitCh := make(chan bool, 1)

	restService := rest.NewService(cfg.Port, storage)
	go func() {
		err := restService.Start()
		if err != nil {
			log.Error(err)
			exitCh <- true
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	select {
	case <-c:
	case <-exitCh:
	}
}
