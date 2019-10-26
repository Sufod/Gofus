package main

import (
	"bufio"
	"net"
)

type DofusSocket struct {
	conn    net.Conn
	reader  *bufio.Reader
	channel chan (string)
}

func NewDofusSocket() *DofusSocket {
	socket := &DofusSocket{
		channel: make(chan (string), 20),
	}
	return socket
}

//Method to initialize socket conn
func (socket *DofusSocket) init(conn net.Conn) {
	socket.channel = make(chan (string), 20)
	socket.conn = conn
	if socket.reader == nil {
		socket.reader = bufio.NewReader(conn)
	} else {
		socket.reader.Reset(conn)
	}
}

func (socket *DofusSocket) close() {
	socket.conn.Close()
}

//Blocks forever and forward received messages from socket to channel
func (socket *DofusSocket) listen() {
	for {
		message, err := socket.reader.ReadString('\x00')
		if err != nil {
			close(socket.channel)
			return
		}
		socket.channel <- message
	}
}

//Send a message in socket
func (socket *DofusSocket) send(message string) {
	socket.conn.Write(append([]byte(message), '\x00'))
}
