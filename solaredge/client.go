package solaredge

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/golang/glog"
)

type EnergyDetailsResponse struct {
	EnergyDetails struct {
		TimeUnit string `json:"timeUnit"`
		Unit     string `json:"unit"`
		Meters   []struct {
			Type   string `json:"type"`
			Values []struct {
				Date  string  `json:"date"`
				Value float64 `json:"value,omitempty"`
			} `json:"values"`
		} `json:"meters"`
	} `json:"energyDetails"`
}

type PowerDetailsResponse struct {
	PowerDetails struct {
		TimeUnit string `json:"timeUnit"`
		Unit     string `json:"unit"`
		Meters   []struct {
			Type   string `json:"type"`
			Values []struct {
				Date  string  `json:"date"`
				Value float64 `json:"value,omitempty"`
			} `json:"values"`
		} `json:"meters"`
	} `json:"powerDetails"`
}

type DataPeriodResponse struct {
	DataPeriod struct {
		StartDate string `json:"startDate"`
		EndDate   string `json:"endDate"`
	} `json:"dataPeriod"`
}

// https://monitoringapi.solaredge.com/site/<site>/dataPeriod.json?api_key=<key>
const urlBase = "https://monitoringapi.solaredge.com"
const version = "1.0.0"

type Client struct {
	SiteId        int
	ApiKey        string
	HttpClient    *http.Client
	SolarEdgeHost string
}

func NewClient(site int, apiKey string) (*Client, error) {
	httpClient := http.Client{
		Timeout: time.Second * 2,
	}
	return &Client{
		SiteId:        site,
		ApiKey:        apiKey,
		HttpClient:    &httpClient,
		SolarEdgeHost: urlBase,
	}, nil
}

func (c *Client) getURL(method string) (*url.URL, error) {
	u, err := url.Parse(fmt.Sprintf("%s/site/%d/%s", c.SolarEdgeHost, c.SiteId, method))
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Add("api_key", c.ApiKey)
	q.Add("version", version)
	u.RawQuery = q.Encode()
	return u, nil
}

func (c *Client) GetDatePeriod() (*DataPeriodResponse, error) {
	u, err := c.getURL("dataPeriod.json")
	if err != nil {
		return nil, err
	}
	bytes, err := c.fetch(u)
	if err != nil {
		return nil, err
	}

	dpr := new(DataPeriodResponse)
	json.Unmarshal(bytes, dpr)
	return dpr, nil
}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func (c *Client) GetEnergyDetailsForLastWeek(day time.Time) (*EnergyDetailsResponse, error) {
	minusOneWeek, _ := time.ParseDuration("-168h")
	startTime := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location()).Add(minusOneWeek)
	endTime := time.Date(day.Year(), day.Month(), day.Day(), 23, 59, 59, 0, day.Location())
	return c.GetEnergyDetails(startTime, endTime)
}

func (c *Client) GetEnergyDetails(startTime, endTime time.Time) (*EnergyDetailsResponse, error) {
	u, err := c.getURL("energyDetails.json")
	if err != nil {
		return nil, err
	}
	fmt.Printf("Requesting times from %s to %s\n", formatTime(startTime), formatTime(endTime))
	q := u.Query()
	q.Add("timeUnit", "QUARTER_OF_AN_HOUR")
	q.Add("startTime", formatTime(startTime))
	q.Add("endTime", formatTime(endTime))
	u.RawQuery = q.Encode()

	bytes, err := c.fetch(u)
	if err != nil {
		return nil, err
	}

	dpr := new(EnergyDetailsResponse)
	json.Unmarshal(bytes, dpr)
	return dpr, nil
}

func (c *Client) GetPowerDetailsForLastWeek(day time.Time) (*PowerDetailsResponse, error) {
	minusOneWeek, _ := time.ParseDuration("-168h")
	startTime := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location()).Add(minusOneWeek)
	endTime := time.Date(day.Year(), day.Month(), day.Day(), 23, 59, 59, 0, day.Location())
	return c.GetPowerDetails(startTime, endTime)
}
func (c *Client) GetPowerDetails(startTime, endTime time.Time) (*PowerDetailsResponse, error) {
	u, err := c.getURL("powerDetails.json")
	if err != nil {
		return nil, err
	}
	fmt.Printf("Requesting times from %s to %s\n", formatTime(startTime), formatTime(endTime))
	q := u.Query()
	q.Add("timeUnit", "QUARTER_OF_AN_HOUR")
	q.Add("startTime", formatTime(startTime))
	q.Add("endTime", formatTime(endTime))
	u.RawQuery = q.Encode()

	bytes, err := c.fetch(u)
	if err != nil {
		return nil, err
	}

	dpr := new(PowerDetailsResponse)
	json.Unmarshal(bytes, dpr)
	return dpr, nil
}

func (c *Client) fetch(u *url.URL) ([]byte, error) {

	glog.Info("Fetching URL ", u.String())
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
