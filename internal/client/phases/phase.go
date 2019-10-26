package phases

import "github.com/Sufod/Gofus/internal/network"

type PhaseInterface interface {
	HandlePackets(socket *network.DofusSocket)
}

//Phase is a struct defining what a player does depending on the received packets
