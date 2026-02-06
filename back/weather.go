package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func FetchWeatherForCity(city Cities) (*WeatherData, error) {
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current_weather=true",
		city.Latitude, city.Longitude)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error: %s", resp.Status)
	}

	var apiResponse OpenMeteoResponse
	json.NewDecoder(resp.Body).Decode(&apiResponse)

	return &WeatherData{
		CityID:      city.ID,
		CityName:    city.Name,
		Temperature: apiResponse.Current.Temperature,
		Humidity:    apiResponse.Current.Relativehumidity,
		WindSpeed:   apiResponse.Current.Windspeed10,
		Condition:   WeatherCodeToCondition(apiResponse.Current.WeatherCode),
	}, nil
}
