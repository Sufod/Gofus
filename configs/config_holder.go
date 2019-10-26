package configs

type ConfigHolder struct {
	DofusAuthServer string       `yaml:"DOFUS_AUTH_SERVER" envconfig:"DOFUS_AUTH_SERVER"`
	DofusServerName string       `yaml:"DOFUS_SERVER_NAME" envconfig:"DOFUS_SERVER_NAME"`
	DofusVersion    string       `yaml:"DOFUS_VERSION" envconfig:"DOFUS_VERSION"`
	Credentials     *Credentials `yaml:"credentials"`
}

type Credentials struct {
	Username string `yaml:"DOFUS_USERNAME" envconfig:"DOFUS_USERNAME"`
	Password string `yaml:"DOFUS_PASSWORD" envconfig:"DOFUS_PASSWORD"`
}

func Config() ConfigHolder {
	cfg := ConfigHolder{}
	initialize(&cfg)
	return cfg
}
