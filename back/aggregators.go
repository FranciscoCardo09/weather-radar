package main

func CompareWeather(w1, w2 WeatherData) bool {
	temperature := w1.Temperature - w2.Temperature
	humidity := w1.Humidity - w2.Humidity
	windspeed := w1.WindSpeed - w2.WindSpeed
	condition := w1.Condition == w2.Condition

	return temperature == 0 && humidity == 0 && windspeed == 0 && condition
}

func ComputeSummary(weatherData []WeatherData) WeatherSummary {
	var summary WeatherSummary
	summary.ByCondition = make(map[string][]string)

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

	var hotterCity, colderCity, windyCity string
	maxTemp, minTemp, maxWind := weatherData[0].Temperature, weatherData[0].Temperature, weatherData[0].WindSpeed

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

	return summary
}
