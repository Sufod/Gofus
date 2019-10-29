package decoders

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Sufod/Gofus/configs"
	"github.com/Sufod/Gofus/internal/network"
)

var cfg configs.ConfigHolder = configs.Config()

//ServerList is a struct that will contain an array of all available game servers and functions like get server name by id
type ServerList struct {
	packetID    string
	Servers     []Server
	ServerCount int
}

//Server is a struct that will be used to define a server in the list of game servers
type Server struct {
	serverID int
}

//NewServerList is a function that creates a ServerList struct containing an array of servers from the AH packet
func NewServerList(packet string) (serverList *ServerList, err error) {
	serverList = new(ServerList)
	if strings.HasPrefix(packet, "AH") {
		serverList.packetID = "AH"
		servers, err := getServersFromPacket(packet)
		serverList.Servers = servers
		serverList.ServerCount = len(servers)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("Invalid paquet prefix")
	}
	return serverList, nil
}

//SelectServer sends the packet to select the game server
func SelectServer(packet string, socket *network.DofusSocket) {
	//First checks if has characters on the selected server
	splittedPacket := strings.Split(packet, "|")
	hasCharacters := false
	//Checks if the selected server exists
	if GetServerIDByName(cfg.DofusServerName) == 0 {
		fmt.Println("[AUTHPHASE] [ERR] - Serveur " + cfg.DofusServerName + " indisponible / non existant")
		return
	}
	for index := 1; index < len(splittedPacket)-1; index++ {
		server := splittedPacket[index]
		serverInfos := strings.Split(server, ",")
		characterCount, err := strconv.ParseInt(serverInfos[1], 10, 0)
		if string(serverInfos[0]) == strconv.Itoa(GetServerIDByName(cfg.DofusServerName)) && characterCount != 0 {
			hasCharacters = true
		}
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println(hasCharacters)

	if hasCharacters == true {
		fmt.Println("Serveur choisis : " + cfg.DofusServerName)
		socket.Send("Ax" + strconv.Itoa(GetServerIDByName(cfg.DofusServerName)))
		return
	}
	fmt.Println("[AUTHPHASE] [ERR] - Vous n'avez pas de personnage sur le serveur " + cfg.DofusServerName)
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

//ServerExists checks if the choosen server exists in the serverlist
func ServerExists(serverList *ServerList, serverName string) int {
	for index := 0; index < serverList.ServerCount; index++ {
		if serverList.Servers[index].serverID == GetServerIDByName(serverName) {
			return 1
		}
	}
	return 0
}

//GetServerNameByID is a function that returns a server name from its ID
func GetServerNameByID(id int) (serverName string) {
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

//GetServerIDByName is a function that returns a server ID from its name
func GetServerIDByName(name string) (serverID int) {
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
