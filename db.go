package soler

import (
	"fmt"
	"time"

	"github.com/golang/glog"
	influxdb "github.com/influxdata/influxdb/client/v2"
	sense "github.com/kodek/sense-api"
	"github.com/kodek/soler/greenbutton"
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

func (d *Database) AddProductionPoints(siteId int, points map[time.Time]SolarDatapoint) error {
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

func (d *Database) AddConsumptionPoints(points []greenbutton.GBPoint) error {
	glog.Infof("Writing %d consumption data points to InfluxDB", len(points))

	bp, err := influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
		Database:  d.database,
		Precision: "s",
	})
	if err != nil {
		return err
	}

	tags := map[string]string{
		"home_id_temp": "cuesta",
	}

	for _, p := range points {
		fields := map[string]interface{}{
			"usage_kwh": p.UsageKwh,
		}
		dbPoint, err := influxdb.NewPoint("consumption", tags, fields, p.T)
		if err != nil {
			return err
		}
		bp.AddPoint(dbPoint)
	}

	err = d.conn.Write(bp)
	if err != nil {
		return err
	}
	glog.Info("Successfully wrote consumption points to database")
	return nil
}

func (d *Database) AddSenseRealtimePoint(p sense.RealtimeResponse) error {
	bp, err := influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
		Database:  d.database,
		Precision: "s",
	})
	if err != nil {
		return err
	}

	tags := map[string]string{
		"type": p.Type,
	}

	fields := map[string]interface{}{
		"solar_w": p.Payload.SolarW,
		"w":       p.Payload.W,
	}
	dbPoint, err := influxdb.NewPoint("sense_realtime_overall", tags, fields, time.Unix(p.Payload.Epoch, 0))
	if err != nil {
		return err
	}
	bp.AddPoint(dbPoint)

	// Devices
	for _, device := range p.Payload.Devices {
		tags := map[string]string{
			"id":   device.ID,
			"name": device.Name,
		}

		fields := map[string]interface{}{
			"w": device.W,
		}

		devicePoint, err := influxdb.NewPoint("sense_realtime_devices", tags, fields, time.Unix(p.Payload.Epoch, 0))
		if err != nil {
			return err
		}
		bp.AddPoint(devicePoint)
	}

	err = d.conn.Write(bp)
	if err != nil {
		return err
	}
	glog.Info("Successfully wrote Sense point to database")
	return nil
}
