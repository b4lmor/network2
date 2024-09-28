package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type WeatherResponse struct {
	Coord      Coordinates   `json:"coord"`
	Weather    []WeatherInfo `json:"weather"`
	Base       string        `json:"base"`
	Main       MainInfo      `json:"main"`
	Visibility int           `json:"visibility"`
	Wind       WindInfo      `json:"wind"`
	Rain       RainInfo      `json:"rain"`
	Clouds     CloudsInfo    `json:"clouds"`
	Dt         int64         `json:"dt"`
	Sys        SysInfo       `json:"sys"`
	Timezone   int           `json:"timezone"`
	Id         int           `json:"id"`
	Name       string        `json:"name"`
	Cod        int           `json:"cod"`
}

type Coordinates struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type WeatherInfo struct {
	Id          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type MainInfo struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	Pressure  int     `json:"pressure"`
	Humidity  int     `json:"humidity"`
	SeaLevel  int     `json:"sea_level"`
	GrndLevel int     `json:"grnd_level"`
}

type WindInfo struct {
	Speed float64 `json:"speed"`
	Deg   int     `json:"deg"`
	Gust  float64 `json:"gust"`
}

type RainInfo struct {
	OneHour float64 `json:"1h"`
}

type CloudsInfo struct {
	All int `json:"all"`
}

type SysInfo struct {
	Type    int    `json:"type"`
	Id      int    `json:"id"`
	Country string `json:"country"`
	Sunrise int64  `json:"sunrise"`
	Sunset  int64  `json:"sunset"`
}

type OpenWeatherClient struct {
	apiKey string
}

func NewOpenWeatherClient(apiKey string) *OpenWeatherClient {
	return &OpenWeatherClient{apiKey: apiKey}
}

func (c *OpenWeatherClient) FetchWeather(lon, lat float64) (*WeatherResponse, error) {
	reqUrl := "https://api.openweathermap.org/data/2.5/weather"
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Add("lon", fmt.Sprintf("%f", lon))
	query.Add("lat", fmt.Sprintf("%f", lat))
	query.Add("units", "metric")
	query.Add("appid", c.apiKey)
	req.URL.RawQuery = query.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.New(res.Status)
	}

	var response WeatherResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
