package network

import (
	"fmt"
	"net"
)

type AuthHandlerSocket struct {
	dofusSocket *DofusSocket
}

func NewAuthHandlerSocket(dofusSocket *DofusSocket) *AuthHandlerSocket {
	h := &AuthHandlerSocket{
		dofusSocket: dofusSocket,
	}
	return h
}

//Send a message in socket
func (socket AuthHandlerSocket) Send(message string) {
	socket.dofusSocket.Send(message)
}

//WaitForPacket blocks until a message is available to read in the channel
func (socket AuthHandlerSocket) WaitForPacket() (string, error) {
	message, err := socket.dofusSocket.WaitForPacket()
	if err != nil || message == "" {
		return "", err
	}
	return message, nil
}

//CloseSocketAndConnectTo closes the current socket and initialize a new connection to specified ip address
//This method is to be used during the auth phase gameserver connexion step
func (socket AuthHandlerSocket) CloseSocketAndConnectTo(gameServerIp string) error {
	socket.dofusSocket.Close()                       // Closing current socket
	serverConn, err := net.Dial("tcp", gameServerIp) // Connecting to game server
	if err != nil {
		return err
	}
	socket.dofusSocket.Initialize(serverConn) //Initializing server socket conn for game server
	go socket.dofusSocket.Listen()            //Starting server listen loop in a goroutine
	return nil
}

func (socket AuthHandlerSocket) GetHandlerSocket() *HandlerSocket {
	return NewHandlerSocket(socket.dofusSocket)
}

//HandleEmptyPacket reads the next packet and ignore it
func (socket AuthHandlerSocket) HandleEmptyPacket() {
	packet, err := socket.WaitForPacket()
	if err != nil {
		//TODO better error handling
		fmt.Println(err)
	}
	if len(packet) > 0 {
		return
	}
}
