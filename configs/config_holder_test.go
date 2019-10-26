package configs

import (
	"testing"

	"gotest.tools/assert"
)

func TestConfigHolder(t *testing.T) {
	cfg := Config()

	assert.Equal(t, cfg.DofusAuthServer, "34.251.172.139:443")
	assert.Equal(t, cfg.DofusServerName, "Eratz")
	assert.Equal(t, cfg.DofusVersion, "1.29.1")
	assert.Equal(t, cfg.Credentials.Username, "myDofusAccountName")
	assert.Equal(t, cfg.Credentials.Password, "myDofusAccountPassword")
}
