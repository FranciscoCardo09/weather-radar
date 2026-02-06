package main

import "fmt"

func main() {

	//Prueaba de GetCityByID
	r := GetCityByID("cordoba")
	fmt.Printf("Ciudad: %s, Latitud: %.2f°C, Longitud: %.2f%%\n", r.Name, r.Latitude, r.Longitude)

	//Prueba de FetchWeatherForCity
	w, err := FetchWeatherForCity(*r)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Ciudad: %s, Temperatura: %.2f°C, Humedad: %.2f%%, Viento: %.2f km/h, Condición: %s\n",
		w.CityName, w.Temperature, w.Humidity, w.WindSpeed, w.Condition)
}
