package crypto

import (
	"testing"

	"gotest.tools/assert"
)

func TestCryptPassword(t *testing.T) {
	assert.Equal(t, EncryptPassword("MonSUperp4ssword", "zzybokxyrtkpjvxmmoxbnwiynojxdbqn"), "#1-haj_hNL00YR-9a75Y34YU3ZXX8f_6ZX")
}

func TestDofusCypher(t *testing.T) {
	cypher := NewDofusCypher()
	//ticket := "891cd1"
	ipPort := "34.243.173.66:443"
	encodedIpPort := "22?3:=42ag7"
	assert.Equal(t, cypher.DecodeIp(encodedIpPort), ipPort)
	assert.Equal(t, cypher.EncodeIp(ipPort), encodedIpPort)
}
