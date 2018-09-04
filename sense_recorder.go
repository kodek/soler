package soler

import (
	"time"

	"github.com/golang/glog"
	sense "github.com/kodek/sense-api"
)

type SenseRecorder struct {
	Db *Database
}

func (rec *SenseRecorder) StartAndLoop(config Sense) {
	c, err := sense.NewClient(config.Email, config.Password)
	if err != nil {
		panic(err)
	}
	recv := rec.connectOrDie(c)

	throttler := time.NewTicker(2 * time.Second)
	for range throttler.C {
		response := <-recv

		err := rec.Db.AddSenseRealtimePoint(response)
		if err != nil {
			glog.Fatal("Cannot write to InfluxDb. ", err)
		}
	}
	glog.Fatal("Lost connection to Sense")

}

func (rec *SenseRecorder) connectOrDie(c sense.Client) <-chan sense.RealtimeResponse {
	recv, err := c.Realtime(make(chan struct{}))
	if err != nil {
		glog.Fatal("Cannot connect to Sense realtime service. ", err)
	}
	return recv
}
