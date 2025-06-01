package main

import (
	"log"

	"github.com/Piyuuussshhh/weather-api/api"
	"github.com/Piyuuussshhh/weather-api/cache"
)

func main() {
	if err := cache.Init(); err != nil {
		log.Fatal(err)
	}

	log.Fatal(api.Route())
}