package soler

import (
	"github.com/golang/glog"
	sense "github.com/kodek/sense-api"
)

type SenseRecorder struct {
	Db *Database
}

func (rec *SenseRecorder) StartAndLoop(conf Sense) {
	c, err := sense.NewClient(conf.Email, conf.Password)
	if err != nil {
		panic(err)
	}
	recv := rec.connectOrDie(c)

	for response := range recv {
		err := rec.Db.AddSenseRealtimePoint(response)
		if err != nil {
			glog.Error("Cannot write to InfluxDb. ", err)
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
