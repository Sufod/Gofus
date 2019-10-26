package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

type DofusProxy struct {
	clientSocket *DofusSocket
	serverSocket *DofusSocket
}

func NewDofusProxy() *DofusProxy {
	proxy := &DofusProxy{
		clientSocket: NewDofusSocket(),
		serverSocket: NewDofusSocket(),
	}
	return proxy
}

func (proxy *DofusProxy) start() {
	fmt.Println("Waiting for client to be connected")
	ln, _ := net.Listen("tcp", "127.0.0.1:8081") //Starting listening on local interface
	clientConn, _ := ln.Accept()                 //Blocking until a client connect
	fmt.Println("Client connected")

	proxy.clientSocket.init(clientConn) //Initializing client socket conn
	go proxy.clientSocket.listen()      //Starting client listen loop in a goroutine

	fmt.Println("Establishing connexion with auth server")
	serverConn, err := net.Dial("tcp", "34.251.172.139:443") // Connecting to auth servers
	if err != nil {
		log.Panic(err)
	}
	proxy.serverSocket.init(serverConn) //Initializing server socket conn
	go proxy.serverSocket.listen()      //Starting server listen loop in a goroutine

	fmt.Println("Connected, starting logging packets until game server choice")
	fmt.Println("=======================================")

	authMessage, err := proxy.listenAndForwardAuthSequence() //Starting proxy blocking loop
	if err != nil {
		proxy.clientSocket.close()
		proxy.serverSocket.close()
		log.Panicln(err)
	}

	fmt.Println("=======================================")
	fmt.Println("Handling connexion to game server")
	ticket := authMessage[11:]
	encodedIp := authMessage[:11]
	cypher := NewDofusCypher()
	ip := cypher.decodeIp(encodedIp)
	proxy.serverSocket.close()
	fmt.Println("Waiting for client to be connected again")
	go func() {
		clientConn, _ = ln.Accept() //Blocking until a client connect
		fmt.Println("Client connected")
		proxy.clientSocket.init(clientConn) //Initializing client socket conn
		proxy.clientSocket.listen()         //Starting client listen
	}()

	proxy.clientSocket.send("AXK" + cypher.encodeIp("127.0.0.1:8081") + ticket)
	time.Sleep(1 * time.Second)
	fmt.Println("Establishing connexion with game server " + ip)
	serverConn, err = net.Dial("tcp", ip) // Connecting to game server
	if err != nil {
		log.Panic(err)
	}
	proxy.serverSocket.init(serverConn) //Initializing server socket conn for game server
	go proxy.serverSocket.listen()      //Starting server listen loop in a goroutine

	err = proxy.listenAndForward() //Starting proxy blocking loop
	if err != nil {
		log.Println(err)
	}
	proxy.serverSocket.close()
	proxy.clientSocket.close()
	fmt.Println("stopped proxy")
}

//Blocks forever and forward + print received messages from client to server and vice-versa
//Gracefully close if client disconnect
func (proxy *DofusProxy) listenAndForward() error {
	for {
		select {
		case message, ok := <-proxy.clientSocket.channel:
			fmt.Println("Message from client: " + message)
			if ok == false || message == "" {
				return errors.New("Client closed connection, stopping...")
			}
			proxy.serverSocket.send(message)
		case message, ok := <-proxy.serverSocket.channel:
			fmt.Println("Message from server: " + message)
			if ok == false || message == "" {
				return errors.New("Server closed connection, stopping...")
			}
			proxy.clientSocket.send(message)
		}
	}
}

//Blocks forever and forward + print received messages from client to server and vice-versa
//Gracefully close if client disconnect
func (proxy *DofusProxy) listenAndForwardAuthSequence() (string, error) {
	for {
		select {
		case message, ok := <-proxy.clientSocket.channel:
			if ok == false || message == "" {
				return "", errors.New("Client closed connection, stopping...")
			}
			fmt.Println("Message from client: " + message)
			proxy.serverSocket.send(message)
		case message, ok := <-proxy.serverSocket.channel:
			if ok == false || message == "" {
				return "", errors.New("Server closed connection, stopping...")
			}
			fmt.Println("Message from server: " + message)
			if strings.HasPrefix(message, "AXK") {
				return message[3:], nil
			}
			proxy.clientSocket.send(message)
		}
	}
}
