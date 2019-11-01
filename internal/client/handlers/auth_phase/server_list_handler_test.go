package handlers

import (
	"testing"

	"gotest.tools/assert"
)

func TestGetServersFromPacket(t *testing.T) {
	var serverList serverList
	var err error

	regularpacket := "AH601;1;110;0|605;1;110;0|609;1;110;0|604;1;110;0|608;1;110;0|603;1;110;0|607;1;110;0|611;1;110;0|602;1;110;0|606;1;110;0|610;1;110;0"
	emptypacket := "AH"

	serverList.Servers, err = getServersFromPacket(regularpacket) //expected = [{601} {605} {609} {604} {608} {603} {607} {611} {602} {606} {610}]
	assert.NilError(t, err)
	assert.Equal(t, len(serverList.Servers), 11)

	serverList.Servers, err = getServersFromPacket(emptypacket) //expected = []
	assert.Equal(t, len(serverList.Servers), 0)
	assert.Assert(t, err != nil)
}

func TestNewServerList(t *testing.T) {
	regularpacket := "AH601;1;110;0|605;1;110;0|609;1;110;0|604;1;110;0|608;1;110;0|603;1;110;0|607;1;110;0|611;1;110;0|602;1;110;0|606;1;110;0|610;1;110;0"
	emptypacket := "AH"

	serverList, err := newServerList(regularpacket)
	assert.NilError(t, err)
	assert.Equal(t, len(serverList.Servers), 11)

	serverList, err = newServerList(emptypacket)
	assert.Assert(t, serverList == nil)
	assert.Assert(t, err != nil)

	serverList, err = newServerList("")
	assert.Assert(t, serverList == nil)
	assert.Assert(t, err != nil)
}
