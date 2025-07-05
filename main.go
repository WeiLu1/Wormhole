package main

import (
	"log"
	"os"

	"github.com/WeiLu1/wormhole/config"
	"github.com/WeiLu1/wormhole/server"
)

func main() {

	configurationFilePath := os.Args[1:]
	if len(configurationFilePath) != 1 {
		log.Fatal("Provide path to configuration file")
	}
	configuration := config.GetConfig(configurationFilePath[0])

	svr := server.NewServer(configuration)

	log.Fatal(svr.Run())
}
