package main

import (
	"flag"

	"time"

	"bitbucket.org/kodek64/soler"
	"github.com/golang/glog"
)

func main() {
	flag.Set("logtostderr", "true")
	flag.Parse()

	glog.Info("Loading config")
	config := soler.LoadConfig()

	client, err := soler.NewClient(config)
	if err != nil {
		panic(err)
	}

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

	ticker := time.NewTicker(1 * time.Hour)
	// Don't run the first tick instantly, because if there's a crash loop, we'll keep making requests.
	glog.Info("Starting wait loop. First tick will be in 1 hour...")
	for range ticker.C {
		s.GetDataForToday()
		glog.Info("Waiting for 1 hour...")
	}
}