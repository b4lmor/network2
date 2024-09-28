package consumer

import (
	"fmt"
	"lab3/context"
	"log"
	"strconv"
	"strings"
	"utils"
)

func ConsumeGraphHopperResponse(ctx *context.Context) {
	for {
		response := <-ctx.PlacesChan
		for i, hit := range response.Hits {
			fmt.Printf(
				"%d. %s %s %s (%.5f, %.5f)\n",
				i+1, hit.Country, hit.City, hit.Postcode, hit.Point.Lng, hit.Point.Lat,
			)
		}

		fmt.Print("Choose place (enter a number): ")

		index, err := strconv.Atoi(utils.ReadLine())
		if err != nil {
			log.Println(err)
			continue
		}

		first := response.Hits[index-1]
		go requestWeather(first.Point.Lng, first.Point.Lat, ctx)
		go requestLocations(first.Point.Lng, first.Point.Lat, ctx)
		ctx.Wg.Done()
	}
}

func requestWeather(lon, lat float64, ctx *context.Context) {
	r, err := ctx.OpenWeatherClient.FetchWeather(lon, lat)
	if err != nil {
		log.Println(err)
		return
	}
	ctx.WeatherChan <- r
}

func requestLocations(lon, lat float64, ctx *context.Context) {
	r, err := ctx.OpenTripMapClient.FetchObjects(lon, lat)
	if err != nil {
		log.Println(err)
		return
	}
	ctx.LocationChan <- r
}

func ConsumeOpenWeatherResponse(ctx *context.Context) {
	for {
		response := <-ctx.WeatherChan
		fmt.Printf("Weather: %.1f°C\n", response.Main.Temp)
		ctx.Wg.Done()
	}
}

func ConsumeOpenTripMapLocationResponse(ctx *context.Context) {
	for {
		response := <-ctx.LocationChan
		var builder strings.Builder
		builder.WriteString("The closest attractions:\n")
		for _, r := range response.Features {
			name := r.Properties.Name
			//log.Println(r.Properties.XID)
			k, err := ctx.OpenTripMapClient.FetchProperties(r.Properties.XID)
			if err != nil {
				continue
			}
			if len(name) == 0 {
				continue
			}
			desc := utils.RemoveHTMLTags(k.Info.Descr)
			desc = utils.DecodeHTMLEntities(desc)
			desc = strings.TrimSpace(desc)
			desc = strings.ReplaceAll(desc, "\t", " ")
			desc = strings.ReplaceAll(desc, "\n", " ")
			if len(desc) > 0 {
				desc = ": " + desc
			}
			runes := []rune(desc)
			var result string
			if len(runes) > 40 {
				runes = runes[:80]
				result = string(runes) + "..."
			} else {
				result = string(runes)
			}
			builder.WriteString(fmt.Sprintf("\t- %s%s\n", name, result))
		}
		fmt.Println(builder.String())
		ctx.Wg.Done()
	}
}
