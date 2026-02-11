package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// FIX: Cliente HTTP global con connection pooling para mejor performance
// Reutilizar el cliente evita crear nuevas conexiones TCP en cada request
var weatherClient *http.Client

// InitWeatherClient inicializa el cliente HTTP global con el timeout configurado
func InitWeatherClient(timeout time.Duration) {
	weatherClient = &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		},
	}
}

// WeatherCodeToCondition convierte un código de clima de Open-Meteo a descripción en español
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

// FetchWeatherForCity obtiene los datos meteorológicos actuales de la API Open-Meteo
// para la ciudad especificada. Respeta el contexto para cancelación y timeout.
//
// Retorna error si:
//   - El contexto es cancelado
//   - La API responde con status != 200
//   - La respuesta JSON no puede ser parseada
func FetchWeatherForCity(ctx context.Context, city Cities) (*WeatherData, error) {
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current=temperature_2m,relative_humidity_2m,wind_speed_10m,weather_code",
		city.Latitude, city.Longitude)

	// FIX: Crear request con context para respetar cancelación
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// FIX: Usar cliente HTTP global con connection pooling
	resp, err := weatherClient.Do(req)
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
