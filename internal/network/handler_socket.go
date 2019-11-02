package network

import "fmt"

type HandlerSocket struct {
	dofusSocket *DofusSocket
}

func NewHandlerSocket(dofusSocket *DofusSocket) *HandlerSocket {
	h := &HandlerSocket{
		dofusSocket: dofusSocket,
	}
	return h
}

//Send a message in socket
func (socket HandlerSocket) Send(message string) {
	socket.dofusSocket.Send(message)
}

//WaitForPacket blocks until a message is available to read in the channel
func (socket HandlerSocket) WaitForPacket() (string, error) {
	message, err := socket.dofusSocket.WaitForPacket()
	if err != nil || message == "" {
		return "", err
	}
	return message, nil
}

//HandleEmptyPacket reads the next packet and ignore it
func (socket HandlerSocket) HandleEmptyPacket() {
	packet, err := socket.WaitForPacket()
	if err != nil {
		//TODO better error handling
		fmt.Println(err)
	}
	if len(packet) > 0 {
		return
	}
}
