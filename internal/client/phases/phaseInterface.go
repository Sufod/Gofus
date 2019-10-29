package phases

//PhaseInterface is the interface that permit to each phase to handle packets
type PhaseInterface interface {
	HandlePackets()
}

//Phase is a struct defining what a player does depending on the received packets
