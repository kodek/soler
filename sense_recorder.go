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
	recv, _, err := c.Realtime()
	if err != nil {
		glog.Fatal("Cannot connect to Sense realtime service. ", err)
	}

	for {
		response := <-recv
		err := rec.Db.AddSenseRealtimePoint(response)
		if err != nil {
			glog.Fatal("Cannot write to InfluxDb. ", err)
		}
	}

}
