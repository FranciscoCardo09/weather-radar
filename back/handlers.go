package main

import (
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
	cityID1 := c.Param("city_id1") // Obtiene el ID de la primera ciudad desde la URL
	cityID2 := c.Param("city_id2") // Obtiene el ID de la segunda ciudad desde la URL

	city1 := GetCityByID(cityID1) // Busca la primera ciudad por su ID
	city2 := GetCityByID(cityID2) // Busca la segunda ciudad por su ID

	if city1 == nil || city2 == nil {
		c.JSON(400, gin.H{"error": "Ciudad no encontrada"})
		return
	}

	weather1, err1 := FetchWeatherForCity(*city1) // ← ctx PRIMERO
	weather2, err2 := FetchWeatherForCity(*city2)

	if err1 != nil || err2 != nil {
		c.JSON(500, gin.H{"error": "Error al obtener el clima de una o ambas ciudades"})
		return
	}

	comparison := CompareWeather(*weather1, *weather2) // Compara los datos del clima de ambas ciudades
	c.JSON(200, gin.H{"comparison": comparison})       // Devuelve JSON con el resultado de la comparación
}
