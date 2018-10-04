package soler

import (
	"time"

	"github.com/jpillora/backoff"

	"github.com/golang/glog"
	sense "github.com/kodek/sense-api"
)

type SenseRecorder struct {
	Db             *Database
}

func (rec *SenseRecorder) StartAndLoop(config Sense) {
	c, err := sense.NewClient(config.Email, config.Password)
	if err != nil {
		glog.Fatal("Cannot authenticate to Sense. ", err)
	}

	backoffTracker := backoff.Backoff{
		Min:    1 * time.Second,
		Max:    15 * time.Minute,
		Factor: 1.5,
		Jitter: true,
	}

	// Run indefinitely. Reconnect with exponential backoff if needed.
	for {
		rec.recordIndefinitely(c, &backoffTracker)

		throttleTime := backoffTracker.Duration()
		glog.Error("Reconnecting in ", throttleTime.String())
		time.Sleep(throttleTime)
	}
}

func (rec *SenseRecorder) recordIndefinitely(client sense.Client, backoffTracker *backoff.Backoff) {
	recv := rec.connectOrDie(client)

	throttler := time.NewTicker(2 * time.Second)
	for response := range recv {
		// Successful communication, so reset the backoff tracker.
		backoffTracker.Reset()

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
