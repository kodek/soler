package soler

import (
	"bitbucket.org/kodek64/soler/sense"
	"github.com/golang/glog"
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
