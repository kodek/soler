package sense

import (
	"log"
	"net/url"

	"fmt"

	"encoding/json"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
)

type RealtimeResponse struct {
	Type    string `json:"type"`
	Payload struct {
		Hz       float64   `json:"hz"`
		Channels []float64 `json:"channels"`
		Devices  []struct {
			W    float64 `json:"w"`
			Name string  `json:"name"`
			Icon string  `json:"icon"`
			ID   string  `json:"id"`
			Tags struct {
				DefaultUserDeviceType       string `json:"DefaultUserDeviceType"`
				TimelineAllowed             string `json:"TimelineAllowed"`
				UserDeviceType              string `json:"UserDeviceType"`
				UserEditable                string `json:"UserEditable"`
				UserDeviceTypeDisplayString string `json:"UserDeviceTypeDisplayString"`
				DeviceListAllowed           string `json:"DeviceListAllowed"`
				//UserEditable                string `json:"user_editable"`
			} `json:"tags"`
		} `json:"devices"`
		W     float64 `json:"w"`
		Stats struct {
			Mrcv float64 `json:"mrcv"`
			Brcv float64 `json:"brcv"`
			Msnd float64 `json:"msnd"`
		} `json:"_stats"`
		Epoch   int64         `json:"epoch"`
		Deltas  []interface{} `json:"deltas"`
		Voltage []float64     `json:"voltage"`
		Frame   int           `json:"frame"`
		SolarW  float64       `json:"solar_w"`
	} `json:"payload"`
}

func (c *ClientImpl) Realtime() (<-chan RealtimeResponse, chan<- struct{}, error) {
	u, err := url.Parse(fmt.Sprintf("wss://clientrt.sense.com/monitors/%d/realtimefeed?access_token=%s", c.monitorId, c.accessToken))
	if err != nil {
		return nil, nil, err
	}

	recv, done, err := connect(*u)
	if err != nil {
		return nil, nil, err
	}

	recvParsed := make(chan RealtimeResponse)
	go func() {
		for {
			select {
			case msg := <-recv:
				var r RealtimeResponse
				err := json.Unmarshal(msg, &r)
				if err != nil {
					glog.Errorf("Cannot parse to JSON: ", err, string(msg))
					return
				}
				recvParsed <- r
			}

		}
	}()
	return recvParsed, done, nil
}

func connect(u url.URL) (<-chan []byte, chan<- struct{}, error) {
	//u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	done := make(chan struct{})
	sendChan := make(chan []byte)

	go func() {
		defer c.Close()
		<-done
	}()

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				glog.Errorf("read error:", err)
				return
			}
			select {
			case sendChan <- message:
				glog.Infof("Sent ", string(message))
			default:
				glog.Infof("Nothing sent")
			}
		}
	}()

	return sendChan, done, err
}
