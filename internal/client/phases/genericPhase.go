package phases

import (
	"github.com/Sufod/Gofus/internal/network"
)

//GenericPhase is a struct containing all the common points of phases
type GenericPhase struct {
	*network.DofusSocket
	startingPackedID string
	endingPacketID   string
	onStartAction    func(message string)
	onEndAction      func(message string)
}
