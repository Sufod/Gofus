package phases

import (
	"fmt"
	"strings"
	"time"

	"github.com/Sufod/Gofus/configs"
	"github.com/Sufod/Gofus/internal/network"
	"github.com/Sufod/Gofus/tools/crypto"
)

var cfg configs.ConfigHolder = configs.Config()

//AuthPhase is the phase where the user logins and choose its game server / character
type AuthPhase struct {
	GenericPhase
}

//AuthSendAuthentication sends the username, encryptedpass and version to the server
func AuthSendAuthentication(message string, socket *network.DofusSocket) {
	socket.Send(cfg.DofusVersion)
	key := strings.TrimPrefix(message, "HC")
	cryptedPassword := crypto.EncryptPassword(cfg.Credentials.Password, key)
	socket.Send(cfg.Credentials.Username + "\n" + cryptedPassword)
	socket.Send("Af")
}

//AuthConnectToGameServer disconnects from the authserver to finally connect to the gameServer and init GamePhase
func AuthConnectToGameServer(message string, socket *network.DofusSocket) {

}

//HandlePackets handles packets for the auth phase
func (phase *AuthPhase) HandlePackets(socket *network.DofusSocket) {
	// TODO
	// Handle auth pĥase
	for {
		select {
		case message, ok := <-socket.Channel:
			if ok == false || message == "" {
				fmt.Println("Server closed connection, stopping...")
				return
			}
			fmt.Println("[AUTHPHASE] [RCV] - " + message)
			time.Sleep(100)
			switch {

			case strings.HasPrefix(message, phase.startingPackedID):
				phase.onStartAction(message, socket)
			case strings.HasPrefix(message, phase.endingPacketID):
				phase.onEndAction(message, socket)
			default:
				fmt.Println("[AUTHPHASE] [ERR] - Cant handle packet")
			}
		}
	}
}

//NewAuthPhase Initialize the AuthPhase's
func NewAuthPhase() PhaseInterface {
	fmt.Println("========= ENTERING AUTH PHASE =========")
	authPhase := &AuthPhase{
		GenericPhase: GenericPhase{},
	}
	authPhase.initAuthPhase()
	return authPhase
}

func (authPhase *AuthPhase) initAuthPhase() {
	authPhase.startingPackedID = "HC"
	authPhase.endingPacketID = "AXK"
	authPhase.onStartAction = AuthSendAuthentication
	authPhase.onEndAction = AuthConnectToGameServer
	//authPhase.packetHandler = phases.AuthHandlePacket
}
