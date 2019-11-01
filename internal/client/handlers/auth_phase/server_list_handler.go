package handlers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Sufod/Gofus/internal/network"
)

type serverHandler struct {
	*network.HandlerSocket
	serverList         serverList
	selectedServerName string
}

func newServerHandler(socket *network.HandlerSocket, selectedServerName string) serverHandler {
	serverHandler := serverHandler{
		HandlerSocket:      socket,
		selectedServerName: selectedServerName,
	}

	return serverHandler
}

//serverList is a struct that will contain an array of all available game servers and functions like get server name by id
type serverList struct {
	packetID    string
	Servers     []Server
	ServerCount int
}

//Server is a struct that will be used to define a server in the list of game servers
type Server struct {
	serverID int
}

//newServerList is a function that creates a serverList struct containing an array of servers from the AH packet
func newServerList(packet string) (*serverList, error) {
	//serverList = new(serverList)
	if strings.HasPrefix(packet, "AH") {
		serverList := serverList{
			packetID: "AH",
		}
		servers, err := getServersFromPacket(packet)
		if err != nil {
			return nil, err
		}
		serverList.Servers = servers
		serverList.ServerCount = len(servers)
		return &serverList, nil
	} else {
		return nil, errors.New("Invalid paquet prefix")
	}
}

//HandleServerList directly handles the serverlist from the packet and anwser to it
func (serverHandler serverHandler) handleServerList() {
	packet, err := serverHandler.WaitForPacket()
	if err != nil {
		//TODO better error handling
		fmt.Println(err)
	}
	serverList, err := newServerList(packet)
	if err != nil {
		fmt.Println(err)
	} else {
		if serverList.ServerCount > 0 && serverExists(serverList, serverHandler.selectedServerName) == 1 {
			serverHandler.Send("Ax")
		} else {
			fmt.Println("[AUTHPHASE] [ERR] - Serveur " + serverHandler.selectedServerName + " indisponible / non existant")
		}
	}

}

//SelectServer sends the packet to select the game server
func (serverHandler serverHandler) selectServer() {
	serverHandler.Send("Ax")
	packet, err := serverHandler.WaitForPacket() //AxK
	if err != nil {
		//TODO better error handling
		fmt.Println(err)
	}
	//TODO Check for AxK packet
	//First checks if has characters on the selected server
	splittedPacket := strings.Split(packet, "|")
	hasCharacters := false
	//Checks if the selected server exists
	if getServerIdByName(serverHandler.selectedServerName) == 0 {
		fmt.Println("[AUTHPHASE] [ERR] - Serveur " + serverHandler.selectedServerName + " indisponible / non existant")
		return
	}
	for index := 1; index < len(splittedPacket)-1; index++ {
		server := splittedPacket[index]
		serverInfos := strings.Split(server, ",")
		characterCount, err := strconv.ParseInt(serverInfos[1], 10, 0)
		if string(serverInfos[0]) == strconv.Itoa(getServerIdByName(serverHandler.selectedServerName)) && characterCount != 0 {
			hasCharacters = true
		}
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println(hasCharacters)

	if hasCharacters == true {
		fmt.Println("Serveur choisis : " + serverHandler.selectedServerName)
		serverHandler.Send("Ax" + strconv.Itoa(getServerIdByName(serverHandler.selectedServerName)))
		return
	}
	fmt.Println("[AUTHPHASE] [ERR] - Vous n'avez pas de personnage sur le serveur " + serverHandler.selectedServerName)
	return
}

//getServersFromPacket is a function that creates an array of server ID from the AH packet
func getServersFromPacket(packet string) (servers []Server, err error) {
	if strings.Contains(packet, "|") && strings.Contains(packet, "AH") {
		packetContent := strings.TrimPrefix(packet, "AH")   // removes the AH from the packet
		fullServerData := strings.Split(packetContent, "|") // 601;1;110;0|605;1;110;0|609;1;110;0 -> {601;1;110;0, 605;1;110;0, 609;1;110;0}
		for i := 0; i < len(fullServerData); i++ {
			serversInfo := strings.Split(fullServerData[i], ";") // 601;1;110;0 -> {601, 1, 110, 0}
			if len(serversInfo) == 4 {
				serverID, err := strconv.Atoi(serversInfo[0]) // 601 ...
				if err != nil {
					return nil, errors.New("server ID Atoi Failed")
				}
				server := Server{serverID: serverID}
				servers = append(servers, server)
			} else { // invalid packet, serverInfo length != 4
				return nil, errors.New("Invalid packet: serverInfo length = " + string(len(serversInfo)) + ", expected = 4")
			}
		}
	} else {
		return nil, errors.New("Invalid or empty paquet")
	}
	return servers, nil
}

//serverExists checks if the choosen server exists in the serverlist
func serverExists(serverList *serverList, serverName string) int {
	for index := 0; index < serverList.ServerCount; index++ {
		if serverList.Servers[index].serverID == getServerIdByName(serverName) {
			return 1
		}
	}
	return 0
}

//getServerNameById is a function that returns a server name from its ID
func getServerNameById(id int) (serverName string) {
	switch id { // Source -> https://cadernis.fr/index.php?threads/aide-trouver-les-id-serveur-retro-dofus.2351/
	case 601:
		serverName = "Eratz"
	case 602:
		serverName = "Henual"
	case 603:
		serverName = "Nabur"
	case 604:
		serverName = "Arty"
	case 605:
		serverName = "Algathe"
	case 606:
		serverName = "Hogmeiser"
	case 607:
		serverName = "Droupik"
	case 608:
		serverName = "Ayuto"
	case 609:
		serverName = "Bilby"
	case 610:
		serverName = "Clustus"
	case 611:
		serverName = "Issering"
	case 612:
		serverName = "Boune"
	case 39:
		serverName = "Debug"
	default:
		serverName = "Unknown"
	}
	return serverName
}

//getServerIdByName is a function that returns a server ID from its name
func getServerIdByName(name string) (serverID int) {
	switch name { // Source -> https://cadernis.fr/index.php?threads/aide-trouver-les-id-serveur-retro-dofus.2351/
	case "Eratz":
		serverID = 601
	case "Henual":
		serverID = 602
	case "Nabur":
		serverID = 603
	case "Arty":
		serverID = 604
	case "Algathe":
		serverID = 605
	case "Hogmeiser":
		serverID = 606
	case "Droupik":
		serverID = 607
	case "Ayuto":
		serverID = 608
	case "Bilby":
		serverID = 609
	case "Clustus":
		serverID = 610
	case "Issering":
		serverID = 611
	case "Boune":
		serverID = 612
	case "Debug":
		serverID = 39
	default:
		serverID = 0
	}
	return
}
