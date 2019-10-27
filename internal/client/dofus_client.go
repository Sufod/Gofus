package client

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/Sufod/Gofus/configs"
	"github.com/Sufod/Gofus/internal/client/phases"
	"github.com/Sufod/Gofus/internal/network"
)

var cfg configs.ConfigHolder = configs.Config()

//DofusClient is a struct that will contain the communication between the client and the server
type DofusClient struct {
	serverSocket *network.DofusSocket
}

//Start is a function to init connection to the authServer with a DofusCLient
func (client *DofusClient) Start() error {
	fmt.Println("Establishing connexion with server")
	serverConn, err := net.Dial("tcp", cfg.DofusAuthServer) // Connecting to auth servers
	if err != nil {
		log.Panic(err)
	}
	client.serverSocket = network.NewDofusSocket(serverConn) //Creating and Initializing server socket conn
	defer client.serverSocket.Close()                        //Delaying server socket conn graceful close
	go client.serverSocket.Listen()                          //Starting server listen loop in a goroutine

	fmt.Println("Connected, starting logging packets")
	fmt.Println("=======================================")

	client.listenAndForward() //Starting proxy blocking loop
	fmt.Println("stopped client")
	return nil
}

//PhaseName is an enum of the differents phases
type PhaseName int

const (
	//AUTH is the authphase
	AUTH PhaseName = iota
	//ANOTHER is an exemple phase
	ANOTHER
)

//Blocks forever and forward + print received messages from client to server and vice-versa
//Gracefully close if client disconnect
func (client *DofusClient) listenAndForward() {

	phasesHandlers := make(map[PhaseName]phases.PhaseInterface)
	phasesHandlers[AUTH] = phases.NewAuthPhase()
	//currentPhase := AUTH
	for {
		phasesHandlers[AUTH].HandlePackets(client.serverSocket) // Appel bloquant

		// Executé a la fin de [Auth] HandlePackets
		fmt.Println("Ending auth phase")
		time.Sleep(100)
		break
	}
}
