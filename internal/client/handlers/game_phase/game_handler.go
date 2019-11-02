package handlers

import (
	"fmt"

	"github.com/Sufod/Gofus/internal/network"
)

//gameHandler is the phase where the user logins and choose its game server / character
type gameHandler struct {
	*network.HandlerSocket
}

//NewGameHandler Instantiate and initialize a new gameHandler
func NewGameHandler(socket *network.DofusSocket) gameHandler {
	gameHandler := gameHandler{
		HandlerSocket: network.NewHandlerSocket(socket),
	}
	return gameHandler
}

//Handle handles packets for the game phase
func (gameHandler gameHandler) Handle() {
	fmt.Println("========= ENTERING GAME PHASE =========")
	gameHandler.HandleEmptyPacket() // ATK0
	gameHandler.Send("Ak0")
	gameHandler.Send("AV")
	gameHandler.HandleEmptyPacket() // BN
	gameHandler.HandleEmptyPacket() // AV0
	gameHandler.Send("Agfr")

	//To be continued...

}
