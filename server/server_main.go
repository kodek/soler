package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/kodek/soler"
	"github.com/kodek/soler/solaredge"
)

func main() {
	flag.Set("logtostderr", "true")
	flag.Parse()

	glog.Info("Loading config...")
	config := soler.LoadConfig()

	client, err := solaredge.NewClient(config.SolarEdge.Site, config.SolarEdge.ApiKey)
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

	go recordSolarEdge(s)
	go recordSenseRealtime(s)

	// Start HTTP server.
	glog.Info("Starting HTTP server on port 10000 (/healthz)")
	http.Handle("/healthz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ok")
	}))
	http.Handle("/force", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := s.GetDataForToday()
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err.Error())
		} else {
			fmt.Fprintf(w, "Done")
		}
	}))
	http.Handle("/startsense", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	http.Handle("/upload", &soler.GreenButtonHandler{Db: database})
	http.ListenAndServe(":10000", nil)

}

func recordSolarEdge(s soler.Soler) {
	ticker := time.NewTicker(1 * time.Hour)
	// Don't run the first tick instantly, because if there's a crash loop, we'll keep making requests.
	glog.Info("Starting SolarEdge polling. First tick will be in 1 hour...")
	for range ticker.C {
		s.GetDataForToday()
		glog.Info("Waiting for 1 hour...")
	}
}

func recordSenseRealtime(s soler.Soler) {
	glog.Info("Starting Sense WSS connection...")
	rec := soler.SenseRecorder{
		Db: s.DbClient,
	}
	rec.StartAndLoop(s.Config.Sense)
}
