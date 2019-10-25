package main

import (
	"fmt"
	"os"
)

func main() {
	os.Setenv("DOFUS_ACCOUNT", "gofusnextrie")
	os.Setenv("DOFUS_PASSWORD", "nextriegofus2019")
	var client DofusClient
	err := client.start()
	if err != nil {
		fmt.Println(err)
	}
}
