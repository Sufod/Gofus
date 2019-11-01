package client

import (
	"fmt"
	"log"
	"net"

	"github.com/Sufod/Gofus/configs"
	auth_phase "github.com/Sufod/Gofus/internal/client/handlers/auth_phase"
	"github.com/Sufod/Gofus/internal/network"
)

//DofusClient is a struct that will contain the communication between the client and the server
type dofusClient struct {
	serverSocket *network.DofusSocket
	cfg          configs.ConfigHolder
}

func NewDofusClient(cfg configs.ConfigHolder) *dofusClient {
	proxy := &dofusClient{
		cfg: cfg,
	}
	return proxy
}

//Start is a function to init connection to the authServer with a DofusCLient
func (client *dofusClient) Start() error {
	fmt.Println("Establishing connexion with server")
	serverConn, err := net.Dial("tcp", client.cfg.DofusAuthServer) // Connecting to auth servers
	if err != nil {
		log.Panic(err)
	}
	client.serverSocket = network.NewDofusSocket(serverConn) //Creating and Initializing server socket conn
	defer client.serverSocket.Close()                        //Delaying server socket conn graceful close
	go client.serverSocket.Listen()                          //Starting server listen loop in a goroutine

	fmt.Println("Connected, starting logging packets")
	fmt.Println("=======================================")

	client.handle() //Starting main loop
	fmt.Println("stopped client")
	return nil
}

//handle is the main method for client, is in charge of the orchestration of the different handlers
//This method shouldn't stop until the end of the program
func (client *dofusClient) handle() {
	auth_phase.NewAuthHandler(client.serverSocket, client.cfg).Handle()
	fmt.Println("Ending auth phase")
}
