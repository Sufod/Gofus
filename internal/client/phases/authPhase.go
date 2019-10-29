package phases

import (
	"fmt"
	"strings"

	"github.com/Sufod/Gofus/configs"
	"github.com/Sufod/Gofus/internal/decoders"
	"github.com/Sufod/Gofus/internal/network"
	"github.com/Sufod/Gofus/tools/crypto"
)

var cfg configs.ConfigHolder = configs.Config()

//AuthPhase is the phase where the user logins and choose its game server / character
type AuthPhase struct {
	GenericPhase
}

//NewAuthPhase Initialize the AuthPhase's
func NewAuthPhase(socket *network.DofusSocket) PhaseInterface {
	authPhase := &AuthPhase{
		GenericPhase: GenericPhase{
			DofusSocket:      socket,
			startingPackedID: "HC",
			endingPacketID:   "AXK",
		},
	}

	return authPhase
}

//SendAuthentication sends the username, encryptedpass and version to the server
func (authPhase *AuthPhase) SendAuthentication(message string) {
	authPhase.Send(cfg.DofusVersion)
	key := message[2:]
	cryptedPassword := crypto.EncryptPassword(cfg.Credentials.Password, key)
	authPhase.Send(cfg.Credentials.Username + "\n" + cryptedPassword)
	authPhase.Send("Af")
}

//ConnectToGameServer disconnects from the authserver to finally connect to the gameServer and init GamePhase
func (authPhase *AuthPhase) ConnectToGameServer(message string) {

}

//HandleServerList directly handles the serverlist from the packet and anwser to it
func (authPhase *AuthPhase) HandleServerList(packet string) {
	serverList, err := decoders.NewServerList(packet)
	if err != nil {
		fmt.Println(err)
	} else {
		if serverList.ServerCount > 0 && decoders.ServerExists(serverList, cfg.DofusServerName) == 1 {
			authPhase.Send("Ax")
		} else {
			fmt.Println("[AUTHPHASE] [ERR] - Serveur " + cfg.DofusServerName + " indisponible / non existant")
		}
	}

}

//HandlePackets handles packets for the auth phase
func (authPhase *AuthPhase) HandlePackets() {
	fmt.Println("========= ENTERING AUTH PHASE =========")
	for {
		select {
		case message, ok := <-authPhase.DofusSocket.Channel:
			if ok == false || message == "" {
				fmt.Println("Server closed connection, stopping...")
				return
			}
			//	fmt.Println("[AUTHPHASE] [RCV] - " + message)
			switch {

			case strings.HasPrefix(message, authPhase.startingPackedID):
				authPhase.SendAuthentication(message)
				//message, ok := <-authPhase.DofusSocket.Channel
				//Ici on s'attends a recevoir un paquet particulier
				// On peut aussi deleguer le travail a un sous-handler si necessaire

			case strings.HasPrefix(message, authPhase.endingPacketID):
				authPhase.ConnectToGameServer(message)
			case strings.HasPrefix(message, "Af"):
				if decoders.HandleQueue(message).IsSub == false {
					fmt.Println("[AUTHPHASE] [ERR] - Un compte non abonné ne peux pas jouer sur Dofus retro")
					return
				}
			case strings.HasPrefix(message, "Ad"):
				fmt.Println("Connecté a " + strings.TrimPrefix(message, "Ad"))
			case strings.HasPrefix(message, "Ac"):
				//Empty packet
			case strings.HasPrefix(message, "AlK0"):
				//Empty packet
			case strings.HasPrefix(message, "AQ"):
				//Empty packet
			case strings.HasPrefix(message, "AH"):
				authPhase.HandleServerList(message)
			case strings.HasPrefix(message, "AxK"):
				decoders.SelectServer(message, authPhase.DofusSocket)
			default:
				fmt.Println("[AUTHPHASE] [RCV] - " + message)
				fmt.Println("[AUTHPHASE] [ERR] - Cant handle packet")
				break
			}
		}
	}
}
