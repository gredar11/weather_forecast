package clients

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type CityWeatherResp struct {
	Location struct {
		Name string `json:"name"`
	} `json:"location"`
	RealTimeInfo struct {
		Temperature float32 `json:"temp_c"`
		GeneralInfo struct {
			Text string `json:"text"`
		} `json:"condition"`
		MeasurementTime hourDate `json:"last_updated"`
	} `json:"current"`
}

type GetWeatherInCityReq struct {
	CityName string
}

type WeatherApiClient struct {
	httpClient *http.Client
}

func NewWeatherApiClient(client *http.Client) *WeatherApiClient {
	return &WeatherApiClient{client}
}

func (c *WeatherApiClient) GetWeatherInCity(ctx context.Context, cityName string) (*CityWeatherResp, error) {
	base, err := url.Parse("https://api.weatherapi.com/v1/current.json")
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Add("q", cityName)
	params.Add("lang", "ru")
	params.Add("key", "1234") //todo get key from env-var
	base.RawQuery = params.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", base.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var cityResp CityWeatherResp
	err = json.NewDecoder(resp.Body).Decode(&cityResp)
	if err != nil {
		return nil, err
	}
	return &cityResp, nil
}

type hourDate struct {
	time.Time
}

func (d *hourDate) UnmarshalJSON(b []byte) error {
	layout := "2006-01-02 15:04"

	s := strings.Trim(string(b), "\"") // remove quotes
	if s == "null" {
		return nil
	}
	var err error
	d.Time, err = time.Parse(layout, s)
	return err
}
