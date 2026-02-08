package main

import (
	"context"
	"fmt"
)

func FetchWeatherForCities(ctx context.Context, cities []Cities) []WeatherData {
	results := make([]WeatherData, len(cities))
	done := make(chan struct{})

	//FAN-OUT
	for i, city := range cities {
		go func(i int, city Cities) {
			data, err := FetchWeatherForCity(city)
			if err != nil {
				fmt.Printf("Error fetching weather for %s: %v\n", city.Name, err)
			} else if data != nil {
				results[i] = *data
			}
			done <- struct{}{}
		}(i, city)
	}

	//FAN-IN
	for range cities {
		select {
		case <-done:
		case <-ctx.Done():
			return nil
		}
	}

	return results
}
