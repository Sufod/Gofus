package main

import (
	"github.com/Sufod/Gofus/configs"
	"github.com/Sufod/Gofus/internal/client"
)

var cfg configs.ConfigHolder = configs.Config()

func main() {
	client := client.NewDofusClient(cfg)
	client.Start()
}
