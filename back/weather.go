package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func WeatherCodeToCondition(code int) string {
	switch {
	case code == 0:
		return "Despejado"
	case code > 0 && code <= 3:
		return "Parcialmente nublado"
	case code > 3 && code <= 48:
		return "Nublado"
	case code > 48 && code <= 57:
		return "Lluvia ligera"
	case code > 57 && code <= 67:
		return "Lluvia moderada"
	case code > 67 && code <= 77:
		return "Lluvia intensa"
	case code > 77 && code <= 86:
		return "Nieve ligera"
	case code > 86 && code <= 95:
		return "Nieve moderada"
	case code > 95 && code <= 99:
		return "Nieve intensa"
	default:
		return "Desconocido"
	}
}

// FIX: Añadido context.Context para permitir cancelación y timeout
// Esto permite que las requests HTTP respeten el timeout y la cancelación del usuario
func FetchWeatherForCity(ctx context.Context, city Cities) (*WeatherData, error) {
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current=temperature_2m,relative_humidity_2m,wind_speed_10m,weather_code",
		city.Latitude, city.Longitude)

	// FIX: Crear request con context para respetar cancelación
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// FIX: Usar cliente HTTP con timeout de 10 segundos
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error: %s", resp.Status)
	}

	var apiResponse OpenMeteoResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, err
	}

	return &WeatherData{
		CityID:      city.ID,
		CityName:    city.Name,
		Temperature: apiResponse.Current.Temperature,
		Humidity:    apiResponse.Current.Relativehumidity,
		WindSpeed:   apiResponse.Current.Windspeed10,
		Condition:   WeatherCodeToCondition(apiResponse.Current.WeatherCode),
	}, nil
}
