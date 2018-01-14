package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

var config Config

func init() {
	err := envconfig.Process("ultime", &config)
	if err != nil {
		log.Fatal(err.Error())
	}
}

type Config struct {
	DataPath       string `default:"/tmp/ultime/"`
	Port           uint32 `default:"8080"`
	SavingInterval uint32 `default:"300"`

	//Client config field
	ServerUrl string `default:"127.0.0.1:8080"`
}

func GetConfig() Config {
	return config
}
