package main

import (
	"os"
	"io/ioutil"

	"SIPConfig"
	"SIPCore"
	"SIPServer"
)

func main() {
	configPath := os.Args[1]
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	config, err := SIPConfig.LoadConfig(configData)
	if err != nil {
		panic(err)
	}

	ctx := SIPCore.NewContext(config)

	server := SIPServer.NewAPIServer(ctx)
	server.Run()
}
