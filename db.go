package soler

import (
	"fmt"
	"time"

	"github.com/golang/glog"
	influxdb "github.com/influxdata/influxdb/client/v2"
)

type Database struct {
	conn     influxdb.Client
	database string
}

type SolarDatapoint struct {
	Energy int
	Power  int
}

func NewDatabaseConnection(address string, username string, password string, database string) (*Database, error) {
	// Create a new HTTPClient
	c, err := influxdb.NewHTTPClient(influxdb.HTTPConfig{
		Addr:     address,
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	return &Database{
		conn:     c,
		database: database,
	}, nil
}

func (d *Database) AddPoints(siteId int, points map[time.Time]SolarDatapoint) error {
	glog.Infof("Writing %d data points to InfluxDB", len(points))

	bp, err := influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
		Database: d.database, Precision: "s",
	})
	if err != nil {
		return err
	}

	tags := map[string]string{
		"site_id": fmt.Sprintf("%d", siteId),
	}

	for t, p := range points {
		fields := map[string]interface{}{
			"energy": p.Energy,
			"power":  p.Power,
		}
		dbPoint, err := influxdb.NewPoint("production", tags, fields, t)
		if err != nil {
			return err
		}
		bp.AddPoint(dbPoint)
	}

	err = d.conn.Write(bp)
	if err != nil {
		return err
	}
	glog.Info("Successfully wrote points to database")
	return nil
}