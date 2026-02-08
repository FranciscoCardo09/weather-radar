package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	//Prueaba de GetCityByID
	// r := GetCityByID("cordoba")
	// fmt.Printf("Ciudad: %s, Latitud: %.2f°C, Longitud: %.2f%%\n", r.Name, r.Latitude, r.Longitude)

	//Prueba de FetchWeatherForCity
	// w, err := FetchWeatherForCity(*r)
	// if err != nil {
	//	fmt.Println("Error:", err)
	//	return
	// }
	// fmt.Printf("Ciudad: %s, Temperatura: %.2f°C, Humedad: %.2f%%, Viento: %.2f km/h, Condición: %s\n",
	//	w.CityName, w.Temperature, w.Humidity, w.WindSpeed, w.Condition)

	//Endpoint de prueba
	router := gin.Default()

	//Rutas para el frontend
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))

	router.GET("/api/cities", GetCitiesHandler)
	router.GET("/api/weather/:city_id", GetWeatherHandler)
	router.GET("/api/compare/:city_id1/:city_id2", CompareWeatherHandler)

	router.Run(":8080")
}
