package main

type Cities struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type OpenMeteoResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timezone  string  `json:"timezone"`
	Current   struct {
		Time             string  `json:"time"`
		Temperature      float64 `json:"temperature_2m"`
		Relativehumidity float64 `json:"relative_humidity_2m"`
		Windspeed10      float64 `json:"wind_speed_10m"`
		WeatherCode      int     `json:"weather_code"`
	} `json:"current"`
}

type WeatherData struct {
	CityID      string  `json:"city_id"`
	CityName    string  `json:"city_name"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	WindSpeed   float64 `json:"wind_speed"`
	Condition   string  `json:"condition"`
}

type WeatherResult struct {
	Datos *WeatherData `json:"datos"`
	// FIX: Cambiado de 'error' a '*string' para correcta serialización JSON
	// El tipo 'error' no se serializa correctamente y siempre retorna null
	Error *string `json:"error,omitempty"`
}

type WeatherSummary struct {
	AverageTemperature float64             `json:"average_temperature"`
	AverageHumidity    float64             `json:"average_humidity"`
	AverageWindSpeed   float64             `json:"average_wind_speed"`
	HotterCity         string              `json:"hotter_city"`
	ColderCity         string              `json:"colder_city"`
	WindyCity          string              `json:"windy_city"`
	Ranking            []string            `json:"ranking"`
	ByCondition        map[string][]string `json:"by_condition"`
}

type RankingEntry struct {
	CityName    string  `json:"city_name"`
	Temperature float64 `json:"temperature"`
}

type CompareResult struct {
	Cities  []WeatherData  `json:"cities"`
	Summary WeatherSummary `json:"summary"`
	// FIX: Cambiado de 'error' a '*string' para correcta serialización JSON
	Error *string `json:"error,omitempty"`
}

type CompareRequest struct {
	CityIDs []string `json:"city_ids"`
}
