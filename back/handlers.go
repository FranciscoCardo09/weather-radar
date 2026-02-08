package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// GetCitiesHandler maneja la solicitud para obtener la lista de ciudades
func GetCitiesHandler(c *gin.Context) {
	cities := GetCities() // Obtiene las 15 ciudades argentinas
	c.JSON(200, cities)   // Devuelve JSON: [{"id":"cordoba","name":"Córdoba"...}]
}

func GetWeatherHandler(c *gin.Context) {
	cityID := c.Param("city_id") // Obtiene el ID de la ciudad desde la URL
	city := GetCityByID(cityID)  // Busca la ciudad por su ID

	weather, err := FetchWeatherForCity(*city) // Obtiene el clima para la ciudad
	if err != nil {
		c.JSON(500, gin.H{"error": "Error al obtener el clima"})
		return
	}

	c.JSON(200, weather) // Devuelve JSON con los datos del clima
}

func CompareWeatherHandler(c *gin.Context) {
	var req CompareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "JSON inválido"})
		return
	}

	var cities []Cities
	for _, id := range req.CityIDs {
		city := GetCityByID(id)
		if city == nil {
			c.JSON(400, gin.H{"error": fmt.Sprintf("Ciudad con ID %s no encontrada", id)})
			return
		}
		cities = append(cities, *city)
	}

	weatherData := FetchWeatherForCities(c.Request.Context(), cities)
	summary := ComputeSummary(weatherData)

	response := CompareResult{
		Cities:  weatherData,
		Summary: summary,
		Error:   nil,
	}
	c.JSON(200, response)
}
