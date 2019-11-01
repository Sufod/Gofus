package network

import "net"

type ProxySocket struct {
	*DofusSocket
}

func NewProxySocket(conn net.Conn) *ProxySocket {
	h := &ProxySocket{
		DofusSocket: NewDofusSocket(conn),
	}
	return h
}

//Send a message in socket
func (socket ProxySocket) Send(message string) {
	socket.conn.Write(append([]byte(message), '\x00'))
}
