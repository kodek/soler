package soler

import (
	"errors"
	"time"

	"github.com/golang/glog"
)

type Soler struct {
	Config   Configuration
	Client   *Client
	DbClient *Database
}

func (s *Soler) GetDataForToday() error {
	today := time.Now()
	energy, err := s.Client.GetEnergyDetailsForLastWeek(today)
	if err != nil {
		return err
	}
	glog.Info("Received energy data:", energy)

	production := energy.EnergyDetails.Meters[0]
	if production.Type != "Production" {
		return errors.New("Expected production as first and only meter type")
	}

	pointMap := make(map[time.Time]SolarDatapoint)

	for _, value := range production.Values {
		vt, err := time.ParseInLocation("2006-01-02 15:04:05", value.Date, today.Location())
		if err != nil {
			return err
		}
		wh := int(value.Value)

		glog.Info("found energy data point ", vt, wh)

		existingPoint, ok := pointMap[vt]
		if ok {
			existingPoint.Energy = wh
		} else {
			pointMap[vt] = SolarDatapoint{
				Energy: wh,
			}
		}
	}

	power, err := s.Client.GetPowerDetailsForLastWeek(today)
	if err != nil {
		return err
	}
	glog.Info("Received power data:", power)

	production = power.PowerDetails.Meters[0]
	if production.Type != "Production" {
		return errors.New("Expected production as first and only meter type")
	}

	for _, value := range production.Values {
		vt, err := time.ParseInLocation("2006-01-02 15:04:05", value.Date, today.Location())
		if err != nil {
			return err
		}
		w := int(value.Value)

		glog.Info("found power data point ", vt, w)

		point, ok := pointMap[vt]
		if ok {
			point.Power = w
		} else {
			point = SolarDatapoint{
				Power: w,
			}
		}
		pointMap[vt] = point
	}
	err = s.DbClient.AddPoints(s.Config.SolarEdge.Site, pointMap)
	if err != nil {
		return err
	}
	glog.Info("Successfully uploaded day of ", today)
	return nil
}
