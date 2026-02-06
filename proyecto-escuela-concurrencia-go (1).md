# Proyecto Escuela: Weather Radar 🌦️

## Dashboard de clima multi-ciudad con procesamiento concurrente

---

## Descripción

Aplicación web que permite seleccionar múltiples ciudades y obtener un dashboard comparativo con datos meteorológicos en tiempo real. El backend en Go consulta una API de clima de forma concurrente (una goroutine por ciudad), agrega los resultados usando un patrón fan-out/fan-in, y expone los datos procesados vía REST. El frontend en SvelteKit permite seleccionar ciudades y visualizar los resultados.

---

## API Pública

**Open-Meteo** — `https://open-meteo.com/`

API meteorológica gratuita, sin registro, sin API key, sin límites restrictivos. Documentación clara y respuestas rápidas.

Ejemplo de consulta:

```
GET https://api.open-meteo.com/v1/forecast?latitude=-31.42&longitude=-64.18&current=temperature_2m,wind_speed_10m,relative_humidity_2m,weather_code&timezone=auto
```

---

## Arquitectura

```
┌──────────────────────────────────────────────────────────┐
│                    SvelteKit (SPA)                        │
│                                                          │
│  Selector de ciudades  →  Cards por ciudad  →  Resumen   │
│     (checkboxes)         (temp/hum/viento)    (extremos,  │
│                                                promedios, │
│                                                ranking)   │
└────────────────────────────┬─────────────────────────────┘
                             │ HTTP/JSON
                             ▼
┌──────────────────────────────────────────────────────────┐
│                    Go Backend (:8080)                     │
│                                                          │
│  GET  /api/cities              Lista de ciudades          │
│  POST /api/weather/compare     Comparación concurrente    │
│  GET  /api/weather/:cityId     Consulta individual        │
│                                                          │
│  El endpoint /compare:                                   │
│    1. Recibe lista de ciudades                           │
│    2. Lanza una goroutine por ciudad (fan-out)           │
│    3. Cada goroutine consulta Open-Meteo                 │
│    4. Los resultados se envían a un channel (fan-in)     │
│    5. Un collector agrega: ranking, promedios, extremos  │
│    6. Devuelve respuesta JSON unificada                  │
└────────────────────────────┬─────────────────────────────┘
                             │
                             ▼
                       Open-Meteo API
```

---

## Endpoints

| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET` | `/api/cities` | Devuelve la lista de ciudades disponibles (hardcodeadas con nombre, latitud y longitud) |
| `POST` | `/api/weather/compare` | Recibe lista de IDs de ciudades, consulta Open-Meteo concurrentemente, devuelve datos individuales y agregados |
| `GET` | `/api/weather/:cityId` | Consulta individual del clima para una sola ciudad |

---

## Datos a procesar y mostrar

### Por ciudad

- Temperatura actual
- Humedad relativa
- Velocidad del viento
- Condición climática (derivada del `weather_code`)

### Resumen agregado

- Ranking de ciudades por temperatura (más caliente a más fría)
- Promedios de temperatura, humedad y viento del grupo seleccionado
- Extremos: ciudad más caliente, más fría, más ventosa
- Agrupación por condición climática (soleado, nublado, lluvioso)

---

## Funcionalidad del Frontend

- Selector con ~15 ciudades predefinidas (checkboxes)
- Botón "Comparar" que dispara la consulta al backend
- Cards individuales por ciudad con sus datos de clima
- Panel de resumen con los datos agregados
- Indicador de loading mientras se procesan las consultas
- Manejo de estados vacío y de error

---

## Plan de trabajo (2 días, cada estudiante realiza el proyecto completo)

### Día 1 — Bases

| Bloque | Tarea |
|--------|-------|
| Mañana | Inicializar proyecto Go. Definir los structs de datos. Implementar la función que consulta Open-Meteo para una sola ciudad y parsea la respuesta. Testear con `curl`. |
| Tarde | Implementar `GET /api/cities` y `GET /api/weather/:cityId`. Inicializar proyecto SvelteKit en modo SPA. Crear layout básico con el selector de ciudades y hacer `fetch` al endpoint de ciudades. |

### Día 2 — Concurrencia e integración

| Bloque | Tarea |
|--------|-------|
| Mañana | Implementar la lógica concurrente: fan-out con goroutines, fan-in con channel, agregación de resultados. Exponer vía `POST /api/weather/compare`. Agregar manejo de errores parciales y timeout con `context`. |
| Tarde | Conectar el frontend al endpoint de comparación. Renderizar cards por ciudad y panel de resumen. Agregar loading, manejo de errores y pulir la presentación. |

### Progresión recomendada

Antes de implementar la versión concurrente, implementar primero la comparación de forma secuencial (un loop que consulta una ciudad a la vez) y medir cuánto tarda. Después refactorizar a concurrente y comparar tiempos.

---

## Conceptos de concurrencia que practica el proyecto

| Concepto | Dónde se aplica |
|----------|----------------|
| **Goroutines** | Una goroutine por ciudad para consultar Open-Meteo en paralelo |
| **Channels** | Channel tipado para recolectar los resultados de cada goroutine |
| **Fan-out** | Lanzar N goroutines simultáneas para N ciudades |
| **Fan-in** | Un solo channel donde convergen todos los resultados |
| **sync.WaitGroup** | Esperar a que todas las goroutines terminen antes de cerrar el channel |
| **context.Context con timeout** | Cancelar consultas si la API externa no responde a tiempo |
| **Manejo de errores parciales** | Si una ciudad falla, las demás siguen y los errores se reportan por separado |
| **Patrón result wrapper** | Struct que encapsula dato exitoso o error en un solo tipo para enviar por el channel |

---

## Stack y dependencias

- **Backend**: Go 1.22+, stdlib `net/http` + un router liviano: usar Gin
- **Frontend**: SvelteKit con SSR deshabilitado (modo SPA), Node.js 18+
- **Sin base de datos, sin autenticación, sin deploy a producción**
