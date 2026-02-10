package main

import (
	"context"
	"log"
)

// FIX: Estructura para manejar resultados de manera thread-safe
// Evita race conditions al usar channels en lugar de escribir directamente al slice
type weatherResult struct {
	index int
	data  *WeatherData
	err   error
}

func FetchWeatherForCities(ctx context.Context, cities []Cities) []WeatherData {
	// FIX: Usar channel con buffer para evitar race conditions
	// Cada goroutine envía su resultado por el channel en lugar de escribir directamente
	resultsChan := make(chan weatherResult, len(cities))
	results := make([]WeatherData, 0, len(cities))

	// FAN-OUT: Lanzar una goroutine por ciudad
	for i, city := range cities {
		go func(i int, city Cities) {
			// FIX: Pasar context a FetchWeatherForCity para respetar cancelación
			data, err := FetchWeatherForCity(ctx, city)
			resultsChan <- weatherResult{i, data, err}
		}(i, city)
	}

	// FAN-IN: Recolectar resultados de manera thread-safe
	resultMap := make(map[int]WeatherData)
	for range cities {
		select {
		case res := <-resultsChan:
			if res.err != nil {
				// FIX: Usar log estructurado en lugar de fmt.Printf
				log.Printf("[ERROR] Error fetching weather for city at index %d: %v", res.index, res.err)
			} else if res.data != nil {
				resultMap[res.index] = *res.data
			}
		case <-ctx.Done():
			// FIX: Si se cancela el contexto, retornar nil inmediatamente
			log.Printf("[WARN] Request cancelled by context")
			return nil
		}
	}

	// FIX: Construir slice final manteniendo el orden original
	for i := 0; i < len(cities); i++ {
		if data, ok := resultMap[i]; ok {
			results = append(results, data)
		}
	}

	return results
}
