package client

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Point struct {
	Lng float64 `json:"lng"`
	Lat float64 `json:"lat"`
}

type Hit struct {
	OsmID    int    `json:"osm_id"`
	OsmType  string `json:"osm_type"`
	Country  string `json:"country"`
	OsmKey   string `json:"osm_key"`
	City     string `json:"city"`
	OsmValue string `json:"osm_value"`
	Postcode string `json:"postcode"`
	Name     string `json:"name"`
	Point    Point  `json:"point"`
}

type HitsResponse struct {
	Hits []Hit `json:"hits"`
	Took int   `json:"took"`
}

type GraphHopperClient struct {
	apiKey string
}

func NewGraphHopperClient(apiKey string) *GraphHopperClient {
	return &GraphHopperClient{apiKey: apiKey}
}

func (c *GraphHopperClient) FetchPlaces(location string) (*HitsResponse, error) {
	reqUrl := "https://graphhopper.com/api/1/geocode"
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Add("q", location)
	query.Add("key", c.apiKey)
	req.URL.RawQuery = query.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.New(res.Status)
	}

	var response HitsResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
