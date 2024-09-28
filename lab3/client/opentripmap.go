package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type OpenTripClient struct {
	apiKey string
}

type FeatureCollection struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

type Feature struct {
	Type       string     `json:"type"`
	ID         string     `json:"id"`
	Geometry   Geometry   `json:"geometry"`
	Properties Properties `json:"properties"`
}

type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type Properties struct {
	XID      string  `json:"xid"`
	Name     string  `json:"name"`
	Dist     float64 `json:"dist"`
	Rate     int     `json:"rate"`
	OSM      string  `json:"osm"`
	Wikidata string  `json:"wikidata,omitempty"`
	Kinds    string  `json:"kinds"`
}

type PropertiesResponse struct {
	Kinds     string  `json:"kinds"`
	Sources   Sources `json:"sources"`
	Bbox      Bbox    `json:"bbox"`
	Point     PPoint  `json:"point"`
	OSM       string  `json:"osm"`
	OTM       string  `json:"otm"`
	XID       string  `json:"xid"`
	Name      string  `json:"name"`
	Wikipedia string  `json:"wikipedia"`
	Image     string  `json:"image"`
	Wikidata  string  `json:"wikidata"`
	Rate      string  `json:"rate"`
	Info      Info    `json:"info"`
}

// Структура для источников
type Sources struct {
	Geometry   string   `json:"geometry"`
	Attributes []string `json:"attributes"`
}

// Структура для боксов
type Bbox struct {
	LatMax float64 `json:"lat_max"`
	LatMin float64 `json:"lat_min"`
	LonMax float64 `json:"lon_max"`
	LonMin float64 `json:"lon_min"`
}

// Структура для точки
type PPoint struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

// Структура для информации
type Info struct {
	Descr     string `json:"descr"`
	Image     string `json:"image"`
	ImgWidth  int    `json:"img_width"`
	Src       string `json:"src"`
	SrcID     int    `json:"src_id"`
	ImgHeight int    `json:"img_height"`
}

func NewOpenTripClient(apiKey string) *OpenTripClient {
	return &OpenTripClient{apiKey: apiKey}
}

func (c *OpenTripClient) FetchObjects(lon, lat float64) (*FeatureCollection, error) {
	reqUrl := "https://api.opentripmap.com/0.1/ru/places/radius"
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Add("lon", fmt.Sprintf("%f", lon))
	query.Add("lat", fmt.Sprintf("%f", lat))
	query.Add("radius", "10000")
	query.Add("limit", "5")
	query.Add("apikey", c.apiKey)
	req.URL.RawQuery = query.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.New(res.Status)
	}

	var response FeatureCollection
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *OpenTripClient) FetchProperties(xid string) (*PropertiesResponse, error) {
	reqUrl := "https://api.opentripmap.com/0.1/ru/places/xid/" + xid
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Add("apikey", c.apiKey)
	req.URL.RawQuery = query.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.New(res.Status)
	}

	var response PropertiesResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
