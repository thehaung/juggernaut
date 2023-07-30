package main

import (
	"github.com/thehaung/juggernaut/bootstrap"
	"github.com/thehaung/juggernaut/config"
	"log"
)

func main() {
	conf, err := config.Parse()
	if err != nil {
		log.Fatalf("can't load config, error: %s", err)
	}

	bootstrap.Run(conf)
}
