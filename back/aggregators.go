package main

import "sort"

// ComputeSummary calcula estadísticas agregadas de los datos meteorológicos.
//
// Retorna:
//   - Promedios de temperatura, humedad y velocidad del viento
//   - Ciudades con valores extremos (más caliente, más fría, más ventosa)
//   - Ranking de ciudades por temperatura (descendente)
//   - Agrupación de ciudades por condición climática
//
// Si weatherData está vacío, retorna un WeatherSummary con valores por defecto.
func ComputeSummary(weatherData []WeatherData) WeatherSummary {
	var summary WeatherSummary
	summary.ByCondition = make(map[string][]string)

	// FIX: Early return si no hay datos, evita división por cero y panic
	if len(weatherData) == 0 {
		return summary
	}

	var totalTemp, totalHumidity, totalWind float64
	for _, data := range weatherData {
		totalTemp += data.Temperature
		totalHumidity += data.Humidity
		totalWind += data.WindSpeed

		summary.ByCondition[data.Condition] = append(summary.ByCondition[data.Condition], data.CityName)
	}

	summary.AverageTemperature = totalTemp / float64(len(weatherData))
	summary.AverageHumidity = totalHumidity / float64(len(weatherData))
	summary.AverageWindSpeed = totalWind / float64(len(weatherData))

	// FIX: Inicializar variables de extremos correctamente antes del loop
	// Evita panic si el slice está vacío o tiene valores incorrectos
	maxTemp := weatherData[0].Temperature
	minTemp := weatherData[0].Temperature
	maxWind := weatherData[0].WindSpeed
	hotterCity := weatherData[0].CityName
	colderCity := weatherData[0].CityName
	windyCity := weatherData[0].CityName

	for _, data := range weatherData {
		if data.Temperature > maxTemp {
			maxTemp = data.Temperature
			hotterCity = data.CityName
		}
		if data.Temperature < minTemp {
			minTemp = data.Temperature
			colderCity = data.CityName
		}
		if data.WindSpeed > maxWind {
			maxWind = data.WindSpeed
			windyCity = data.CityName
		}
	}

	summary.HotterCity = hotterCity
	summary.ColderCity = colderCity
	summary.WindyCity = windyCity

	// FIX: Implementar ranking en el backend (antes se calculaba solo en el frontend)
	// Ordenar ciudades por temperatura descendente
	sorted := make([]WeatherData, len(weatherData))
	copy(sorted, weatherData)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Temperature > sorted[j].Temperature
	})

	ranking := make([]string, len(sorted))
	for i, data := range sorted {
		ranking[i] = data.CityName
	}
	summary.Ranking = ranking

	return summary
}
