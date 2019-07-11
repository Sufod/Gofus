package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	proxy := NewDofusProxy()
	proxy.start()
}

type DofusSocket struct {
	conn    net.Conn
	reader  *bufio.Reader
	channel chan (string)
}

//Method to initialize socket with reader and channel
func (socket *DofusSocket) initConn(conn net.Conn) {
	socket.channel = make(chan (string), 20)
	socket.conn = conn
	socket.reader = bufio.NewReader(conn)
}

//Blocks forever and forward received messages from socket to channel
func (socket *DofusSocket) listen() {
	for {
		message, _ := socket.reader.ReadString('\x00')
		socket.channel <- message
	}
}

//Send a message in socket
func (socket *DofusSocket) send(message string) {
	socket.conn.Write(append([]byte(message), '\x00'))
}

type DofusProxy struct {
	clientSocket *DofusSocket
	serverSocket *DofusSocket
}

func NewDofusProxy() *DofusProxy {
	proxy := &DofusProxy{
		clientSocket: &DofusSocket{},
		serverSocket: &DofusSocket{},
	}
	return proxy
}

func (proxy *DofusProxy) start() {
	fmt.Println("Waiting for client to be connected")
	ln, _ := net.Listen("tcp", "127.0.0.1:8081") //Starting listening on local interface
	clientConn, _ := ln.Accept()                 //Blocking until a client connect
	fmt.Println("Client connected")

	proxy.clientSocket.initConn(clientConn) //Initializing client socket conn
	defer proxy.clientSocket.conn.Close()   //Delaying client socket conn graceful close
	go proxy.clientSocket.listen()          //Starting client listen loop in a goroutine

	fmt.Println("Establishing connexion with server")
	serverConn, err := net.Dial("tcp", "34.251.172.139:443") // Connecting to auth servers
	if err != nil {
		log.Panic(err)
	}
	proxy.serverSocket.initConn(serverConn) //Initializing server socket conn
	defer proxy.serverSocket.conn.Close()   //Delaying server socket conn graceful close
	go proxy.serverSocket.listen()          //Starting server listen loop in a goroutine

	fmt.Println("Connected, starting logging packets")
	fmt.Println("=======================================")

	proxy.listenAndForward() //Starting proxy blocking loop
	fmt.Println("stopped proxy")
}

//Blocks forever and forward + print received messages from client to server and vice-versa
//Gracefully close if client disconnect
func (proxy *DofusProxy) listenAndForward() {
	for {
		select {
		case message, ok := <-proxy.clientSocket.channel:
			if ok == false || message == "" {
				fmt.Println("Client closed connection, stopping...")
				return
			}
			fmt.Println("Message from client: " + message)
			proxy.serverSocket.send(message)
		case message := <-proxy.serverSocket.channel:
			fmt.Println("Message from server: " + message)
			proxy.clientSocket.send(message)
		}
	}
}

func cryptPassword(password string, key string) string {
	chArray := []rune{
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p',
		'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D', 'E', 'F',
		'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V',
		'W', 'X', 'Y', 'Z', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-', '_'}
	str := []rune("#1")
	passwd := []rune(password)
	k := []rune(key)
	for i := 0; i < len(passwd); i++ {
		ch := passwd[i]
		ch2 := k[i]
		num2 := int(ch / '\u0010')
		num3 := int(ch % '\u0010')
		index := (num2 + int(ch2)) % len(chArray)
		num5 := (num3 + int(ch2)) % len(chArray)
		str = append(str, chArray[index], chArray[num5])
	}
	return string(str)
}
