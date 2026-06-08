package main

import (
	"github.com/farshidboroomand/jobs-collector/configs"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfg, err := configs.Load()
	if err != nil {
		log.WithError(err).Fatal("failed to load application configurations.")
	}

	log.Infof("application configs %+v \n", cfg)
}
