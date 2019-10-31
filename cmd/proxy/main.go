package main

import (
	"github.com/Sufod/Gofus/configs"
	"github.com/Sufod/Gofus/internal/proxy"
)

var cfg configs.ConfigHolder = configs.Config()

func main() {
	proxy := proxy.NewDofusProxy(cfg.DofusAuthServer)
	proxy.Start()
}
