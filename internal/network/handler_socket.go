package network

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
