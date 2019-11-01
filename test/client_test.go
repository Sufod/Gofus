package test

import (
	"testing"
	"time"

	"github.com/Sufod/Gofus/configs"
	"github.com/Sufod/Gofus/internal/client"
)

func TestClient(t *testing.T) {
	emulator := DofusServerEmulator{}
	go emulator.Start(t)

	var cfg configs.ConfigHolder = configs.ConfigHolder{
		DofusAuthServer: "127.0.0.1:8081",
		DofusServerName: "Eratz",
		DofusVersion:    "1.30.1",
		Credentials: &configs.Credentials{
			Username: "testUser",
			Password: "MonSUperp4ssword",
		},
	}

	time.Sleep(1 * time.Second)
	client := client.NewDofusClient(cfg)
	client.Start()
}
