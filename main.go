package main

import (
	"github.com/philippe-vandermoere/docker-proxy-go/proxy"
	log "github.com/sirupsen/logrus"
	"time"
)

func prod() {
	for {
		err := proxy.Run()
		if err != nil {
			log.Error(err)
		}

		break
		time.Sleep(1 * time.Second)
	}
}
