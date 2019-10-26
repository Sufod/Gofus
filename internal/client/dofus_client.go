package client

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/Sufod/Gofus/internal/network"
	"github.com/Sufod/Gofus/tools/crypto"
)

type DofusClient struct {
	serverSocket *network.DofusSocket
}

func (client *DofusClient) start() error {
	if os.Getenv("DOFUS_ACCOUNT") == "" {
		return errors.New("Please setup environment variable DOFUS_ACCOUNT with your account name")
	}
	if os.Getenv("DOFUS_PASSWORD") == "" {
		return errors.New("Please setup environment variable DOFUS_PASSWORD with your account password")
	}
	fmt.Println("Establishing connexion with server")
	serverConn, err := net.Dial("tcp", "34.251.172.139:443") // Connecting to auth servers
	if err != nil {
		log.Panic(err)
	}
	client.serverSocket.Initialize(serverConn) //Initializing server socket conn
	defer client.serverSocket.Close()          //Delaying server socket conn graceful close
	go client.serverSocket.Listen()            //Starting server listen loop in a goroutine

	fmt.Println("Connected, starting logging packets")
	fmt.Println("=======================================")

	client.listenAndForward() //Starting proxy blocking loop
	fmt.Println("stopped client")
	return nil
}

func (d *DofusClient) Start() {
	d.listenAndForward()
}

//Blocks forever and forward + print received messages from client to server and vice-versa
//Gracefully close if client disconnect
func (client *DofusClient) listenAndForward() {
	for {
		select {
		case message, ok := <-client.serverSocket.Channel:
			if ok == false || message == "" {
				fmt.Println("Server closed connection, stopping...")
				return
			}
			fmt.Println("Message from Server: " + message)
			time.Sleep(100)
			switch {
			case strings.HasPrefix(message, "HC"):
				client.serverSocket.Send("1.29.1")
				key := message[2:]
				cryptedPassword := crypto.EncryptPassword(os.Getenv("DOFUS_PASSWORD"), key)
				client.serverSocket.Send(os.Getenv("DOFUS_ACCOUNT") + "\n" + cryptedPassword)
				client.serverSocket.Send("Af")
			case strings.HasPrefix(message, "AQ"):
				client.serverSocket.Send("Ax")
			case strings.HasPrefix(message, "AH"):
				client.serverSocket.Send("AX602")
			case strings.HasPrefix(message, "AXK"):
				// ticket := message[14:]
				// ip := decodeIp(message[3:])
			}
		}
	}
}
