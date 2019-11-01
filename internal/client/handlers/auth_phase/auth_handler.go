package handlers

import (
	"fmt"
	"strings"

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

func (authHandler authHandler) handleAuthenticationResult() bool {
	packet, err := authHandler.WaitForPacket()
	if err != nil {
		//TODO better error handling
		fmt.Println(err)
	}
	isConnected := strings.HasPrefix(packet, "Af")

	if isConnected == false {
		mapping := map[string]string{
			"AlEf": "[AUTHPHASE] [ERR] - Nom de compte ou mot de passe incorrect !",
			"AlEb": "[AUTHPHASE] [ERR] - Ce compte est banni !",
			"AlEn": "[AUTHPHASE] [ERR] - La connexion a été interrompue !",
			"AlEa": "[AUTHPHASE] [ERR] - Compte déjà en cours de connexion !",
			"AlEc": "[AUTHPHASE] [ERR] - Ce compte est déjà connecté a un serveur de jeu",
			"AlEv": "[AUTHPHASE] [ERR] - Nouvelle version ! (" + strings.TrimPrefix(packet, "AlEv") + ")",
			"AlEp": "[AUTHPHASE] [ERR] - Compte invalide !",
			"AlEk": "[AUTHPHASE] [ERR] - Ce compte est banni temporairement",
			"AlEm": "[AUTHPHASE] [ERR] - Compte en maintenance !",
		}
		fmt.Println(mapping[packet[:4]])
	}
	return isConnected
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
	if len(packet) > 0 {
		return
	}
	//TODO Check for XXX packet
}

func (authHandler authHandler) handleUsername() {
	packet, err := authHandler.WaitForPacket()
	if err != nil {
		//TODO better error handling
		fmt.Println(err)
	}
	fmt.Println("Connecté à " + strings.TrimPrefix(packet, "Ad"))
}

//Handle handles packets for the auth phase
func (authHandler authHandler) Handle() {
	fmt.Println("========= ENTERING AUTH PHASE =========")
	authHandler.handleAuthentication() //HC + key

	authHandler.handleQueue() // Af + pos | total | useless data

	isConnected := authHandler.handleAuthenticationResult()

	if isConnected == true {

		authHandler.handleUsername() //Ad + username

		authHandler.handleEmptyPacket() //Ac2

		authHandler.handleServerList() //AH + servers

		authHandler.handleEmptyPacket() //AlK

		authHandler.handleEmptyPacket() //AQ

		authHandler.selectServer()

		authHandler.connectToGameServer()
	}

}
