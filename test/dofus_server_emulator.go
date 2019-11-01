package test

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/Sufod/Gofus/configs"
	"github.com/Sufod/Gofus/internal/client"
	"github.com/Sufod/Gofus/internal/network"
	"gotest.tools/assert"
)

type DofusServerEmulator struct {
	*network.DofusSocket
}

func (emulator *DofusServerEmulator) Start(t *testing.T) {
	fmt.Println("Waiting for client to be connected")
	ln, _ := net.Listen("tcp", "127.0.0.1:8081") //Starting listening on local interface

	go emulator.startClient()
	clientConn, _ := ln.Accept() //Blocking until a client connect
	fmt.Println("Client connected")
	emulator.DofusSocket = network.NewDofusSocket(clientConn) //Creating and initializing client socket conn
	go emulator.Listen()
	defer emulator.Close()
	emulator.handleClient(t)
	fmt.Println("stopped DofusServerEmulator")
}

func (emulator *DofusServerEmulator) handleClient(t *testing.T) {
	emulator.Send("HCzzybokxyrtkpjvxmmoxbnwiynojxdbqn")
	emulator.WaitForPacketAndAssertEqual(t, "1.30.1")
	emulator.WaitForPacketAndAssertEqual(t, "testUser\n#1-haj_hNL00YR-9a75Y34YU3ZXX8f_6ZX")
	emulator.WaitForPacketAndAssertEqual(t, "Af")
	emulator.Send("Af2|3|0||-1")
	emulator.Send("AdtestUser")
	emulator.Send("Ac2")
	emulator.Send("AH601;1;110;0|605;1;110;0|609;1;110;0|604;1;110;0|608;1;110;0|603;1;110;0|607;1;110;0|611;1;110;0|602;1;110;0|606;1;110;0|610;1;110;0")
	emulator.Send("AlK0")
	emulator.Send("AQ")
	emulator.WaitForPacketAndAssertEqual(t, "Ax")
	emulator.Send("AxK0")
}

//WaitForPacket blocks until a message is available to read in the channel
func (emulator *DofusServerEmulator) WaitForPacketAndAssertEqual(t *testing.T, expectedPacket string) {
	packet, err := emulator.WaitForPacket()
	assert.NilError(t, err, "Expected paquet "+expectedPacket)
	assert.Equal(t, packet[:len(packet)-2], expectedPacket)
}

func (emulator *DofusServerEmulator) startClient() {
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