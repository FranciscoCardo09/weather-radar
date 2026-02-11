package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

// GetCitiesHandler maneja GET /api/cities
// Retorna la lista completa de ciudades argentinas disponibles para consulta
func GetCitiesHandler(c *gin.Context) {
	cities := GetCities() // Obtiene las 15 ciudades argentinas
	c.JSON(200, cities)   // Devuelve JSON: [{"id":"cordoba","name":"Córdoba"...}]
}

// GetWeatherHandler maneja GET /api/weather/:city_id
// Retorna los datos meteorológicos actuales para una ciudad específica.
// Responde con 404 si la ciudad no existe, 500 si hay error al obtener el clima.
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

// CompareWeatherHandler maneja POST /api/compare
// Compara el clima de múltiples ciudades y retorna estadísticas agregadas.
// Requiere JSON body con formato: {"city_ids": ["cordoba", "buenosaires", ...]}
//
// Validaciones:
//   - Mínimo 2 ciudades (después de eliminar duplicados)
//   - Máximo 50 ciudades
//   - Todas las ciudades deben existir
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

	// FIX: Eliminar duplicados antes de procesar
	seen := make(map[string]bool)
	var cities []Cities
	for _, id := range req.CityIDs {
		if seen[id] {
			continue // Skip duplicados
		}
		seen[id] = true

		city := GetCityByID(id)
		if city == nil {
			c.JSON(400, gin.H{"error": "Ciudad con ID " + id + " no encontrada"})
			return
		}
		cities = append(cities, *city)
	}

	// FIX: Validar que después de eliminar duplicados quedan al menos 2 ciudades
	if len(cities) < 2 {
		c.JSON(400, gin.H{"error": "Debes proporcionar al menos dos ciudades distintas"})
		return
	}

	weatherData := FetchWeatherForCities(c.Request.Context(), cities)

	// FIX: Validar que se obtuvieron datos antes de computar el resumen
	// En Go, len() de un slice nil retorna 0, por lo que no necesitamos verificar nil explícitamente
	if len(weatherData) == 0 {
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
