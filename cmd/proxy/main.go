package main

import "github.com/Sufod/Gofus/internal/proxy"

func main() {
	proxy := proxy.NewDofusProxy()
	proxy.Start()
}
