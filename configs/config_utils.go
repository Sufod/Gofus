package configs

import (
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func readFile(cfg interface{}) {
	var f *os.File
	var err error
	for _, configPath := range []string{"config.yml", "config.yaml", "configs/config.yml", "configs/config.yaml"} {
		//fmt.Println("Searching "+configPath)
		f, err = os.Open(configPath)
		if err == nil {
			break
		}
		//fmt.Println("Can't found "+configPath)
	}

	if err != nil {
		processError(err)
	}

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		processError(err)
	}
}

func readEnv(cfg interface{}) {
	err := envconfig.Process("", cfg)
	if err != nil {
		processError(err)
	}
}

/* To be used with pointer to config struct */
func initialize(cfg interface{}) {
	readFile(cfg)
	readEnv(cfg)
}
