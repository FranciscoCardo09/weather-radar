package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// Config contiene toda la configuración validada de la aplicación
type Config struct {
	Port           string
	FrontendURL    string
	APITimeout     time.Duration
	RateLimit      float64
	RateLimitBurst int
}

// LoadConfig carga y valida las variables de entorno
func LoadConfig() (*Config, error) {
	portStr := getEnv("PORT", "8080")
	port, err := strconv.Atoi(portStr)
	if err != nil || port < 1 || port > 65535 {
		return nil, fmt.Errorf("invalid PORT: %s (must be 1-65535)", portStr)
	}

	frontendURL := getEnv("FRONTEND_URL", "http://localhost:5173")

	// FIX: Validar que CORS no use wildcard con credentials
	if frontendURL == "*" {
		return nil, fmt.Errorf("FRONTEND_URL cannot be '*' when AllowCredentials is true")
	}

	if _, err := url.Parse(frontendURL); err != nil {
		return nil, fmt.Errorf("invalid FRONTEND_URL: %s", frontendURL)
	}

	// Timeout configurable para requests HTTP a API externa
	timeoutStr := getEnv("API_TIMEOUT", "10")
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil || timeout < 1 || timeout > 60 {
		return nil, fmt.Errorf("invalid API_TIMEOUT: %s (must be 1-60 seconds)", timeoutStr)
	}

	// Rate limiting configurable
	rateLimitStr := getEnv("RATE_LIMIT", "10")
	rateLimit, err := strconv.ParseFloat(rateLimitStr, 64)
	if err != nil || rateLimit < 1 {
		return nil, fmt.Errorf("invalid RATE_LIMIT: %s (must be >= 1)", rateLimitStr)
	}

	rateLimitBurstStr := getEnv("RATE_LIMIT_BURST", "20")
	rateLimitBurst, err := strconv.Atoi(rateLimitBurstStr)
	if err != nil || rateLimitBurst < 1 {
		return nil, fmt.Errorf("invalid RATE_LIMIT_BURST: %s (must be >= 1)", rateLimitBurstStr)
	}

	return &Config{
		Port:           portStr,
		FrontendURL:    frontendURL,
		APITimeout:     time.Duration(timeout) * time.Second,
		RateLimit:      rateLimit,
		RateLimitBurst: rateLimitBurst,
	}, nil
}

// getEnv obtiene una variable de entorno o retorna un valor por defecto
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func main() {
	// FIX: Cargar y validar configuración al inicio
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("[FATAL] Configuration error: %v", err)
	}

	// Inicializar cliente HTTP global con la configuración
	InitWeatherClient(cfg.APITimeout)

	router := gin.Default()

	// FIX: Rate limiting para prevenir abuso
	limiter := rate.NewLimiter(rate.Limit(cfg.RateLimit), cfg.RateLimitBurst)
	router.Use(func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(429, gin.H{"error": "Too many requests. Please try again later."})
			c.Abort()
			return
		}
		c.Next()
	})

	// CORS para el frontend con validación
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.FrontendURL},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))

	router.GET("/api/cities", GetCitiesHandler)
	router.GET("/api/weather/:city_id", GetWeatherHandler)
	router.POST("/api/compare", CompareWeatherHandler)

	// FIX: Graceful shutdown para terminar requests en curso correctamente
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Iniciar servidor en goroutine
	go func() {
		log.Printf("[INFO] Starting server on port %s", cfg.Port)
		log.Printf("[INFO] Rate limit: %.0f req/s, burst: %d", cfg.RateLimit, cfg.RateLimitBurst)
		log.Printf("[INFO] API timeout: %v", cfg.APITimeout)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[FATAL] Server failed: %v", err)
		}
	}()

	// Esperar señal de terminación (Ctrl+C o SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[INFO] Shutting down server...")

	// Dar 5 segundos para terminar requests en curso
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("[ERROR] Server forced to shutdown: %v", err)
	}

	log.Println("[INFO] Server exited gracefully")
}
