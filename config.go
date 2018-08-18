package soler

import (
	"encoding/json"
	"flag"
	"os"
)

type SolarEdge struct {
	Site   int
	ApiKey string
}

type InfluxDbConfig struct {
	Address  string
	Username string
	Password string
	Database string
}

type Configuration struct {
	SolarEdge      SolarEdge
	InfluxDbConfig InfluxDbConfig
}

var configPath = flag.String("config", "", "The path to the config file")

func LoadConfig() Configuration {
	var path string
	if *configPath == "" {
		path = os.Getenv("HOME") + "/.soler_conf.json"
	} else {
		path = *configPath
	}

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(f)
	conf := Configuration{}
	err = decoder.Decode(&conf)
	if err != nil {
		panic(err)
	}
	return conf
}
