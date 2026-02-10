package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// FIX: Eliminar código comentado - Git mantiene el historial

	// FIX: Configuración mediante variables de entorno
	port := getEnv("PORT", "8080")
	frontendURL := getEnv("FRONTEND_URL", "http://localhost:5173")

	router := gin.Default()

	// CORS para el frontend
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{frontendURL},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))

	router.GET("/api/cities", GetCitiesHandler)
	router.GET("/api/weather/:city_id", GetWeatherHandler)
	router.POST("/api/compare", CompareWeatherHandler)

	log.Printf("[INFO] Starting server on port %s", port)
	router.Run(":" + port)
}

// getEnv obtiene una variable de entorno o retorna un valor por defecto
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
