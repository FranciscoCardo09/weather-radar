package main

import (
	"log"

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

	// FIX: Validar que la ciudad existe antes de usarla
	// Previene panic por dereferenciar un puntero nil
	if city == nil {
		c.JSON(404, gin.H{"error": "Ciudad no encontrada"})
		return
	}

	// FIX: Pasar context de la request para respetar timeout y cancelación
	weather, err := FetchWeatherForCity(c.Request.Context(), *city)
	if err != nil {
		log.Printf("[ERROR] Error fetching weather for city %s: %v", cityID, err)
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

	// FIX: Validar que se proporcionaron al menos una ciudad
	if len(req.CityIDs) == 0 {
		c.JSON(400, gin.H{"error": "Debes proporcionar al menos una ciudad"})
		return
	}

	// FIX: Limitar el número máximo de ciudades para evitar abuso
	if len(req.CityIDs) > 50 {
		c.JSON(400, gin.H{"error": "Máximo 50 ciudades por comparación"})
		return
	}

	var cities []Cities
	for _, id := range req.CityIDs {
		city := GetCityByID(id)
		if city == nil {
			c.JSON(400, gin.H{"error": "Ciudad con ID " + id + " no encontrada"})
			return
		}
		cities = append(cities, *city)
	}

	weatherData := FetchWeatherForCities(c.Request.Context(), cities)

	// FIX: Validar que se obtuvieron datos antes de computar el resumen
	if weatherData == nil || len(weatherData) == 0 {
		c.JSON(500, gin.H{"error": "No se pudieron obtener datos del clima"})
		return
	}

	summary := ComputeSummary(weatherData)

	response := CompareResult{
		Cities:  weatherData,
		Summary: summary,
		Error:   nil,
	}
	c.JSON(200, response)
}
