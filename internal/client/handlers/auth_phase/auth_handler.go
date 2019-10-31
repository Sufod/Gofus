package handlers

import (
	"fmt"

	"github.com/Sufod/Gofus/configs"
	"github.com/Sufod/Gofus/internal/network"
	"github.com/Sufod/Gofus/tools/crypto"
)

//authHandler is the phase where the user logins and choose its game server / character
type authHandler struct {
	*network.HandlerSocket
	cfg configs.ConfigHolder
	serverHandler
}

//NewAuthHandler Instantiate and initialize a new authHandler
func NewAuthHandler(socket *network.DofusSocket, cfg configs.ConfigHolder) authHandler {
	authHandler := authHandler{
		HandlerSocket: network.NewHandlerSocket(socket),
		cfg:           cfg,
	}
	authHandler.serverHandler = newServerHandler(authHandler.HandlerSocket, cfg.DofusServerName)

	return authHandler
}

//HandleAuthentication sends the username, encryptedpass and version to the server
func (authHandler authHandler) handleAuthentication() {
	packet, err := authHandler.WaitForPacket()
	if err != nil {
		//TODO better error handling
		fmt.Println(err)
	}
	//TODO Check for HC packet
	authHandler.Send(authHandler.cfg.DofusVersion)
	key := packet[2:]
	cryptedPassword := crypto.EncryptPassword(authHandler.cfg.Credentials.Password, key)
	authHandler.Send(authHandler.cfg.Credentials.Username + "\n" + cryptedPassword)
	authHandler.Send("Af")
}

//ConnectToGameServer disconnects from the authserver to finally connect to the gameServer and init GamePhase
func (authHandler authHandler) connectToGameServer() {
	packet, err := authHandler.WaitForPacket()
	if err != nil {
		//TODO better error handling
		fmt.Println(err)
	}
	fmt.Println(packet)
	//TODO Check for XXX packet
}

//ConnectToGameServer disconnects from the authserver to finally connect to the gameServer and init GamePhase
func (authHandler authHandler) handleEmptyPacket() {
	packet, err := authHandler.WaitForPacket()
	if err != nil {
		//TODO better error handling
		fmt.Println(err)
	}
	fmt.Println(packet)
	//TODO Check for XXX packet
}

//Handle handles packets for the auth phase
func (authHandler authHandler) Handle() {
	fmt.Println("========= ENTERING AUTH PHASE =========")
	authHandler.handleAuthentication()

	authHandler.handleServerList()

	//TODO: placer au bons endroits
	//authHandler.HandleEmptyPacket()

	authHandler.selectServer()

	authHandler.connectToGameServer()

	// for {
	// 	select {
	// 	case message, ok := <-authHandler.DofusSocket.Channel:
	// 		if ok == false || message == "" {
	// 			fmt.Println("Server closed connection, stopping...")
	// 			return
	// 		}
	// 		//	fmt.Println("[AUTHPHASE] [RCV] - " + message)
	// 		switch {

	// 		case strings.HasPrefix(message, authHandler.startingPackedID):
	// 			authHandler.SendAuthentication(message)
	// 			//message, ok := <-authHandler.DofusSocket.Channel
	// 			//Ici on s'attends a recevoir un paquet particulier
	// 			// On peut aussi deleguer le travail a un sous-handler si necessaire

	// 		case strings.HasPrefix(message, authHandler.endingPacketID):
	// 			authHandler.ConnectToGameServer(message)
	// 		case strings.HasPrefix(message, "Af"):
	// 			if decoders.HandleQueue(message).IsSub == false {
	// 				fmt.Println("[AUTHPHASE] [ERR] - Un compte non abonné ne peux pas jouer sur Dofus retro")
	// 				return
	// 			}
	// 		case strings.HasPrefix(message, "Ad"):
	// 			fmt.Println("Connecté a " + strings.TrimPrefix(message, "Ad"))
	// 		case strings.HasPrefix(message, "Ac"):
	// 			//Empty packet
	// 		case strings.HasPrefix(message, "AlK0"):
	// 			//Empty packet
	// 		case strings.HasPrefix(message, "AQ"):
	// 			//Empty packet
	// 		case strings.HasPrefix(message, "AH"):
	// 			authHandler.HandleServerList(message)
	// 		case strings.HasPrefix(message, "AxK"):
	// 			decoders.SelectServer(message, authHandler.DofusSocket, authHandler.cfg.DofusServerName)
	// 		default:
	// 			fmt.Println("[AUTHPHASE] [RCV] - " + message)
	// 			fmt.Println("[AUTHPHASE] [ERR] - Cant handle packet")
	// 			break
	// 		}
	// 	}
	// }
}
