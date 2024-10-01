package main

import (
	"fmt"
	"lab3/consumer"
	"lab3/context"
	"log"
	"utils"
)

func main() {

	ctx := context.NewContext()

	go consumer.ConsumeGraphHopperResponse(ctx)
	go consumer.ConsumeOpenWeatherResponse(ctx)
	go consumer.ConsumeOpenTripMapLocationResponse(ctx)

	for {
		fmt.Print("Enter location: ")
		input := utils.ReadLine()
		if input == "exit" {
			break
		}
		response, err := ctx.GraphHopperClient.FetchPlaces(input)
		if err != nil {
			log.Println(err)
			continue
		}
		ctx.Wg.Add(3)
		ctx.PlacesChan <- response
		ctx.Wg.Wait()
	}

}
