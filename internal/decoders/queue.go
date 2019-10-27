package decoders

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

//Queue is a struct to handle player position in the queue
type Queue struct {
	packetID       string
	currentPos     int
	totalSub       int
	totalNotSubbed int
	IsSub          bool
}

//NewQueue creates a Queue struct to follow the position of the player in the queue
func NewQueue(packet string) (q *Queue, err error) {
	q = new(Queue)
	if strings.Contains(packet, "|") && strings.HasPrefix(packet, "Af") {
		q.packetID = "Af"
		return q, q.UpdateQueuePosition(packet)
	}

	return nil, errors.New("Invalid paquet prefix / content")
}

//UpdateQueuePosition updates an existing queue to get current player queue position
func (q *Queue) UpdateQueuePosition(packet string) (err error) {
	if q == nil || q.packetID != "Af" {
		return errors.New("Queue is not initialized / Does not exist")
	}
	if strings.Contains(packet, "|") && strings.HasPrefix(packet, "Af") {
		packetContent := strings.TrimPrefix(packet, "Af")
		packetInfo := strings.Split(packetContent, "|")
		if len(packetInfo) == 5 {
			q.currentPos, err = strconv.Atoi(packetInfo[0])
			q.totalSub, err = strconv.Atoi(packetInfo[1])
			q.totalNotSubbed, err = strconv.Atoi(packetInfo[2])
			if err != nil {
				return errors.New("queue Atoi Failed")
			}
			if len(packetInfo[3]) == 0 {
				q.IsSub = false
				fmt.Println("[WARN] - Un compte non abonné ne peux pas jouer sur Dofus retro")
			} else {
				q.IsSub = true
			}
		} else { // invalid packet, serverInfo length != 4
			return errors.New("Invalid packet: queue packetContent length = " + string(len(packetInfo)) + ", expected = 5")
		}
	} else {
		return errors.New("Invalid paquet prefix")
	}
	return nil
}

//LogQueuePosition returns a ready to log string containing the player position in the queue
func (q *Queue) LogQueuePosition() {
	fmt.Println("Position dans la file d'attente : " + strconv.Itoa(q.currentPos) + "/" + strconv.Itoa(q.totalSub))
}
