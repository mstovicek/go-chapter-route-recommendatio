package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/mstovicek/go-chapter-route-recommendation/api"
)

func main() {
	apiHttpServer := api.NewApiHttpServer(
		log.New(),
	)
	apiHttpServer.Run()
}
