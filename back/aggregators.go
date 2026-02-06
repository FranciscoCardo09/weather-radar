package main

func CompareWeather(w1, w2 WeatherData) bool {
	temperature := w1.Temperature - w2.Temperature
	humidity := w1.Humidity - w2.Humidity
	windspeed := w1.WindSpeed - w2.WindSpeed
	condition := w1.Condition == w2.Condition

	return temperature == 0 && humidity == 0 && windspeed == 0 && condition
}
