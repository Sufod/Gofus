package phases

import (
	"fmt"
	"strconv"
	"strings"
	"time"

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

//AuthSendAuthentication sends the username, encryptedpass and version to the server
func AuthSendAuthentication(message string, socket *network.DofusSocket) {
	socket.Send(cfg.DofusVersion)
	key := message[2:]
	cryptedPassword := crypto.EncryptPassword(cfg.Credentials.Password, key)
	socket.Send(cfg.Credentials.Username + "\n" + cryptedPassword)
	socket.Send("Af")
}

//AuthConnectToGameServer disconnects from the authserver to finally connect to the gameServer and init GamePhase
func AuthConnectToGameServer(message string, socket *network.DofusSocket) {

}

//HandleQueue directly handles the af packet
func HandleQueue(packet string) *decoders.Queue {
	queue, err := decoders.NewQueue(packet)
	if err != nil {
		fmt.Println(err)
	} else {
		queue.LogQueuePosition()
	}
	return queue
}

//HandleServerList directly handles the serverlist from the packet and anwser to it
func HandleServerList(packet string, socket *network.DofusSocket, step int) {
	if step == 1 {
		serverList, err := decoders.NewServerList(packet)
		if err != nil {
			fmt.Println(err)
		} else {
			if serverList.ServerCount > 0 && decoders.ServerExists(serverList, cfg.DofusServerName) == 1 {
				socket.Send("Ax")
			} else {
				fmt.Println("[AUTHPHASE] [ERR] - Serveur " + cfg.DofusServerName + " indisponible / non existant")
			}
		}
	}
	if step == 2 {
		if decoders.GetServerIDByName(cfg.DofusServerName) != 0 {
			fmt.Println("Serveur choisis : " + cfg.DofusServerName)
			socket.Send("Ax" + strconv.Itoa(decoders.GetServerIDByName(cfg.DofusServerName)))
		} else {
			fmt.Println("[AUTHPHASE] [ERR] - Serveur " + cfg.DofusServerName + " indisponible / non existant")
		}
	}
}

//HandlePackets handles packets for the auth phase
func (phase *AuthPhase) HandlePackets(socket *network.DofusSocket) {
	// TODO
	// Handle auth pĥase
	serverListStep := 1
	for {
		select {
		case message, ok := <-socket.Channel:
			if ok == false || message == "" {
				fmt.Println("Server closed connection, stopping...")
				return
			}
			time.Sleep(100)
			switch {

			case strings.HasPrefix(message, phase.startingPackedID):
				phase.onStartAction(message, socket)
			case strings.HasPrefix(message, phase.endingPacketID):
				phase.onEndAction(message, socket)
			case strings.HasPrefix(message, "Af"):
				if HandleQueue(message).IsSub == false {
					fmt.Println("[AUTHPHASE] [ERR] - Un compte non abonné ne peux pas jouer sur Dofus retro")
					break
				}
			case strings.HasPrefix(message, "Ad"):
				fmt.Println("Connecté a " + strings.TrimPrefix(message, "Ad"))
			case strings.HasPrefix(message, "Ac"):
				//Empty packet
			case strings.HasPrefix(message, "AlK0"):
				//Empty packet
			case strings.HasPrefix(message, "AQ"):
				//Empty packet
			case strings.HasPrefix(message, "AH") || strings.HasPrefix(message, "AxK"):
				if serverListStep <= 2 {
					HandleServerList(message, socket, serverListStep)
				} else {
					break
				}
				serverListStep++
			default:
				fmt.Println("[AUTHPHASE] [RCV] - " + message)
				fmt.Println("[AUTHPHASE] [ERR] - Cant handle packet")
				break
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
