package network

import (
	"bufio"
	"errors"
	"net"
)

type DofusSocket struct {
	conn    net.Conn
	reader  *bufio.Reader
	Channel chan (string)
}

func NewDofusSocket(conn net.Conn) *DofusSocket {
	socket := &DofusSocket{}
	socket.Initialize(conn)
	return socket
}

//Initialize is a Method to initialize socket conn
func (socket *DofusSocket) Initialize(conn net.Conn) {
	socket.Channel = make(chan (string), 20)
	socket.conn = conn
	if socket.reader == nil {
		socket.reader = bufio.NewReader(conn)
	} else {
		socket.reader.Reset(conn)
	}
}

func (socket *DofusSocket) Close() {
	socket.conn.Close()
}

//Listen Blocks forever and forward received messages from socket to channel
func (socket *DofusSocket) Listen() {
	for {
		message, err := socket.reader.ReadString('\x00')
		if err != nil {
			close(socket.Channel)
			return
		}

		socket.Channel <- message
	}
}

//Send a message in socket
func (socket *DofusSocket) Send(message string) {
	//fmt.Println("[SENT] - " + message)
	socket.conn.Write(append([]byte(message), '\n', '\x00'))
}

//WaitForPacket blocks until a message is available to read in the channel
func (socket *DofusSocket) WaitForPacket() (string, error) {
	message, ok := <-socket.Channel
	//fmt.Println(message)
	if ok == false || message == "" {
		return "", errors.New("Remote host closed connection")
	}
	return message, nil
}
