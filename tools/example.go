package main

import (
	"flag"

	"github.com/kodek/soler"
	"github.com/kodek/soler/fake"
)

func main() {
	flag.Set("logtostderr", "true")
	flag.Parse()

	config := soler.Configuration{
		SolarEdge: soler.SolarEdge{
			Site: 809417,
			//ApiKey: API KEY HERE
		},
		InfluxDbConfig: soler.InfluxDbConfig{
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
	client, _ := soler.NewClient(config)
	client.HttpClient = testServer.Client()
	client.SolarEdgeHost = testServer.URL

	dbConfig := config.InfluxDbConfig
	database, err := soler.NewDatabaseConnection(dbConfig.Address, dbConfig.Username, dbConfig.Password, dbConfig.Database)
	if err != nil {
		panic(err)
	}

	s := soler.Soler{
		Client:   client,
		Config:   config,
		DbClient: database,
	}

	s.GetDataForToday()
}
