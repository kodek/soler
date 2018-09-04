package main

import (
	"flag"

	"github.com/kodek/soler/config"
	"github.com/kodek/soler/solaredge"

	"github.com/kodek/soler"
	"github.com/kodek/soler/fake"
)

func main() {
	flag.Set("logtostderr", "true")
	flag.Parse()

	conf := config.Configuration{
		SolarEdge: config.SolarEdge{
			Site: 809417,
			//ApiKey: API KEY HERE
		},
		InfluxDbConfig: config.InfluxDbConfig{
			Address:  "http://docker.lan:8087",
			Username: "soler-dev",
			Password: "soler-dev",
			Database: "soler",
		},
	}

	//client, err := soler.NewClient(config)
	//if err != nil {
	//	panic(err)
	//}

	testServer := fake.NewServer()
	client, _ := solaredge.NewClient(conf.SolarEdge.Site, conf.SolarEdge.ApiKey)
	client.HttpClient = testServer.Client()
	client.SolarEdgeHost = testServer.URL

	dbConfig := conf.InfluxDbConfig
	database, err := soler.NewDatabaseConnection(dbConfig.Address, dbConfig.Username, dbConfig.Password, dbConfig.Database)
	if err != nil {
		panic(err)
	}

	s := soler.Soler{
		Client:   client,
		Config:   conf,
		DbClient: database,
	}

	s.GetDataForToday()
}
