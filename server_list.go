package main

import (
	"errors"
	"strconv"
	"strings"
)

//ServerList is a struct that will contain an array of all available game servers and functions like get server name by id
type ServerList struct {
	packetID    string
	servers     []Server
	serverCount int
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
		serverList.servers = servers
		serverList.serverCount = len(servers)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("Invalid paquet prefix")
	}
	return serverList, nil
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
	default:
		serverName = "Unknown"
	}
	return serverName
}
