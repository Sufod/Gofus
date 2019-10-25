package main

import (
	"testing"

	"gotest.tools/assert"
)

func TestGetServersFromPacket(t *testing.T) {
	var serverList ServerList
	var err error
	regularpacket := "AH601;1;110;0|605;1;110;0|609;1;110;0|604;1;110;0|608;1;110;0|603;1;110;0|607;1;110;0|611;1;110;0|602;1;110;0|606;1;110;0|610;1;110;0"

	serverList.servers, err = getServersFromPacket(regularpacket) //expected = [{601} {605} {609} {604} {608} {603} {607} {611} {602} {606} {610}]
	assert.NilError(t, err)
	assert.Equal(t, len(serverList.servers), 11)

	emptypacket := "AH"
	serverList.servers, err = getServersFromPacket(emptypacket) //expected = []
	assert.Equal(t, len(serverList.servers), 0)
	assert.Assert(t, err != nil)
}
