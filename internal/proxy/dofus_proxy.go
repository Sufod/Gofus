package proxy

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/Sufod/Gofus/internal/network"
	"github.com/Sufod/Gofus/tools/crypto"
)

type dofusProxy struct {
	clientSocket *network.ProxySocket
	serverSocket *network.ProxySocket
	authServerIp string
}

func NewDofusProxy(authServerIp string) *dofusProxy {
	proxy := &dofusProxy{
		authServerIp: authServerIp,
	}
	return proxy
}

func (proxy *dofusProxy) Start() {
	fmt.Println("Waiting for client to be connected")
	ln, _ := net.Listen("tcp", "127.0.0.1:8081") //Starting listening on local interface
	clientConn, _ := ln.Accept()                 //Blocking until a client connect
	fmt.Println("Client connected")
	proxy.clientSocket = network.NewProxySocket(clientConn) //Creating and initializing client socket conn
	go proxy.clientSocket.Listen()                          //Starting client listen loop in a goroutine

	fmt.Println("Establishing connexion with auth server")
	serverConn, err := net.Dial("tcp", proxy.authServerIp) // Connecting to auth servers
	if err != nil {
		log.Panic(err)
	}
	proxy.serverSocket = network.NewProxySocket(serverConn) //Creating and initializing server socket conn
	go proxy.serverSocket.Listen()                          //Starting server listen loop in a goroutine

	fmt.Println("Connected, starting logging packets until game server choice")
	fmt.Println("=======================================")

	authMessage, err := proxy.listenAndForwardAuthSequence() //Starting proxy blocking loop
	if err != nil {
		proxy.clientSocket.Close()
		proxy.serverSocket.Close()
		log.Panicln(err)
	}

	fmt.Println("=======================================")
	fmt.Println("Handling connexion to game server")
	ticket := authMessage[11:]
	encodedIp := authMessage[:11]
	cypher := crypto.NewDofusCypher()
	ip := cypher.DecodeIp(encodedIp)
	proxy.serverSocket.Close()
	fmt.Println("Waiting for client to be connected again")
	go func() {
		clientConn, _ = ln.Accept() //Blocking until a client connect
		fmt.Println("Client connected")
		proxy.clientSocket.Initialize(clientConn) //Initializing client socket conn
		proxy.clientSocket.Listen()               //Starting client listen
	}()

	proxy.clientSocket.Send("AXK" + cypher.EncodeIp("127.0.0.1:8081") + ticket)
	time.Sleep(1 * time.Second)
	fmt.Println("Establishing connexion with game server " + ip)
	serverConn, err = net.Dial("tcp", ip) // Connecting to game server
	if err != nil {
		log.Panic(err)
	}
	proxy.serverSocket.Initialize(serverConn) //Initializing server socket conn for game server
	go proxy.serverSocket.Listen()            //Starting server listen loop in a goroutine

	err = proxy.listenAndForward() //Starting proxy blocking loop
	if err != nil {
		log.Println(err)
	}
	proxy.serverSocket.Close()
	proxy.clientSocket.Close()
	fmt.Println("stopped proxy")
}

//Blocks forever and forward + print received messages from client to server and vice-versa
//Gracefully close if client disconnect
func (proxy *dofusProxy) listenAndForward() error {
	for {
		select {
		case message, ok := <-proxy.clientSocket.Channel:
			fmt.Println("Message from client: " + message)
			if ok == false || message == "" {
				return errors.New("Client closed connection, stopping...")
			}
			proxy.serverSocket.Send(message)
		case message, ok := <-proxy.serverSocket.Channel:
			fmt.Println("Message from server: " + message)
			if ok == false || message == "" {
				return errors.New("Server closed connection, stopping...")
			}
			proxy.clientSocket.Send(message)
		}
	}
}

//Blocks forever and forward + print received messages from client to server and vice-versa
//Gracefully close if client disconnect
func (proxy *dofusProxy) listenAndForwardAuthSequence() (string, error) {
	for {
		select {
		case message, ok := <-proxy.clientSocket.Channel:
			if ok == false || message == "" {
				return "", errors.New("Client closed connection, stopping...")
			}
			fmt.Println("Message from client: " + message)
			proxy.serverSocket.Send(message)
		case message, ok := <-proxy.serverSocket.Channel:
			if ok == false || message == "" {
				return "", errors.New("Server closed connection, stopping...")
			}
			fmt.Println("Message from server: " + message)
			if strings.HasPrefix(message, "AXK") {
				return message[3:], nil
			}
			proxy.clientSocket.Send(message)
		}
	}
}
