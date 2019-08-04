package main

import (
	"github.com/philippe-vandermoere/docker-proxy-go/proxy"
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {
	for {
		err := proxy.Run()
		if err != nil {
			log.Error(err)
		}

		time.Sleep(1 * time.Second)
	}
}
