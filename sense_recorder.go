package soler

import (
	"errors"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/golang/glog"
	sense "github.com/kodek/sense-api"
)

type SenseRecorder struct {
	Db *Database
}

func (rec *SenseRecorder) StartAndLoop(config Sense) {
	c, err := sense.NewClient(config.Email, config.Password)
	if err != nil {
		glog.Fatal("Cannot authenticate to Sense. ", err)
	}

	operation := func() error {
		rec.recordIndefinitely(c)
		glog.Error("Reconnecting...")
		return errors.New("Retryable error")
	}
	err = backoff.Retry(operation, backoff.NewExponentialBackOff())
	if err != nil {
		glog.Fatal(err)
	}
}

func (rec *SenseRecorder) recordIndefinitely(client sense.Client) {
	recv := rec.connectOrDie(client)

	throttler := time.NewTicker(2 * time.Second)
	for response := range recv {
		err := rec.Db.AddSenseRealtimePoint(response)
		if err != nil {
			glog.Fatal("Cannot write to InfluxDb. ", err)
		}

		<-throttler.C
	}
	glog.Error("Lost connection to Sense.")
}

func (rec *SenseRecorder) connectOrDie(c sense.Client) <-chan sense.RealtimeResponse {
	recv, err := c.Realtime(make(chan struct{}))
	if err != nil {
		glog.Fatal("Cannot connect to Sense realtime service. ", err)
	}
	return recv
}
