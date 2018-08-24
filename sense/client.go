package sense

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"io/ioutil"

	"encoding/json"
	"time"

	"github.com/golang/glog"
)

type Client interface {
	Realtime() (<-chan RealtimeResponse, chan<- struct{}, error)
}

type ClientImpl struct {
	monitorId   int
	accessToken string
}

func NewClient(email string, password string) (Client, error) {
	resp, err := authenticate(email, password)
	if err != nil {
		return nil, err
	}
	if !resp.Authorized {
		glog.Error("Not authorized. Response:", resp)
		return nil, err
	}

	return &ClientImpl{
		monitorId:   resp.Monitors[0].ID,
		accessToken: resp.AccessToken,
	}, nil
}

type AuthResponse struct {
	Authorized  bool   `json:"authorized"`
	AccountID   int    `json:"account_id"`
	UserID      int    `json:"user_id"`
	AccessToken string `json:"access_token"`
	Monitors    []struct {
		ID              int    `json:"id"`
		TimeZone        string `json:"time_zone"`
		SolarConnected  bool   `json:"solar_connected"`
		SolarConfigured bool   `json:"solar_configured"`
		Online          bool   `json:"online"`
		Attributes      struct {
			ID          int         `json:"id"`
			Name        string      `json:"name"`
			State       string      `json:"state"`
			Cost        float64     `json:"cost"`
			UserSetCost bool        `json:"user_set_cost"`
			CycleStart  interface{} `json:"cycle_start"`
		} `json:"attributes"`
	} `json:"monitors"`
	BridgeServer string      `json:"bridge_server"`
	PartnerID    interface{} `json:"partner_id"`
	DateCreated  time.Time   `json:"date_created"`
}

func authenticate(email string, password string) (*AuthResponse, error) {
	apiUrl := "https://api.sense.com/apiservice/api/v1/authenticate"
	data := url.Values{}
	data.Set("email", email)
	data.Add("password", password)

	u, _ := url.ParseRequestURI(apiUrl)
	urlStr := u.String() // 'https://api.com/user/'

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, _ := client.Do(r)
	glog.Info("Auth response code ", resp.Status)

	if resp.StatusCode != 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		glog.Error("Cannot authenticate: ", string(body))
		return nil, errors.New("Cannot authenticate. Error " + string(body))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var authResp AuthResponse
	err = json.Unmarshal(body, &authResp)
	if err != nil {
		return nil, err
	}

	return &authResp, nil
}
