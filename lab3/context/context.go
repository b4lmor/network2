package context

import (
	"lab3/client"
	"os"
	"sync"
)

type Context struct {
	PlacesChan        chan *client.HitsResponse
	WeatherChan       chan *client.WeatherResponse
	LocationChan      chan *client.FeatureCollection
	PropertiesChan    chan *client.PropertiesResponse
	GraphHopperClient *client.GraphHopperClient
	OpenWeatherClient *client.OpenWeatherClient
	OpenTripMapClient *client.OpenTripClient
	Wg                sync.WaitGroup
}

func NewContext() *Context {
	return &Context{
		make(chan *client.HitsResponse),
		make(chan *client.WeatherResponse),
		make(chan *client.FeatureCollection),
		make(chan *client.PropertiesResponse),
		client.NewGraphHopperClient(os.Getenv("API_KEY_GRAPHHOPPER")),
		client.NewOpenWeatherClient(os.Getenv("API_KEY_OPENWEATHER")),
		client.NewOpenTripClient(os.Getenv("API_KEY_OPENTRIPMAP")),
		sync.WaitGroup{},
	}
}
