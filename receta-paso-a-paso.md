# Receta Paso a Paso: Implementación de Sistema Cliente-Servidor con Procesamiento Concurrente

---

## Parte 0: Antes de Empezar — Estructura del Proyecto

### Cómo debería organizarse un proyecto de este tipo

Un proyecto que integra un frontend con un backend que consulta servicios externos requiere **separación clara entre responsabilidades**. La idea central es que cada capa conoce solo la interfaz de la siguiente, no cómo funciona internamente.

### Carpetas y archivos típicos (genéricos)

Para un backend (en Go):

```
backend/
├── cmd/
│   └── main.go              # Punto de entrada: inicia servidor, configura rutas
├── internal/
│   ├── models/              # Definiciones de tipos (structs, constantes)
│   ├── handlers/            # Funciones que reciben requests HTTP
│   ├── services/            # Lógica de orquestación y procesamiento
│   ├── integrations/        # Comunicación con servicios externos
│   └── utils/               # Funciones auxiliares sin dependencias
├── go.mod
└── go.sum
```

Para un frontend (en SvelteKit):

```
frontend/
├── src/
│   ├── routes/              # Páginas y rutas
│   ├── components/          # Componentes reutilizables
│   ├── lib/                 # Utilidades, clientes HTTP, lógica
│   └── stores/              # Estado centralizado reactivo
├── svelte.config.js
├── package.json
└── vite.config.js
```

### Responsabilidad de cada capa o módulo

| Capa | Responsabilidad Principal | Qué recibe | Qué devuelve | Qué NO debería hacer |
|------|--------------------------|-----------|--------------|---------------------|
| **main.go** | Iniciar servidor, configurar rutas | - | - | Lógica de negocio, procesamiento |
| **Handlers** | Recibir HTTP, validar entrada, formatear respuesta | Solicitud HTTP | Respuesta JSON | Cálculos complejos, consultas directas |
| **Services** | Orquestar múltiples operaciones, aplicar reglas | Datos validados | Datos procesados | Detalles de HTTP, parsing JSON crudo |
| **Integrations** | Comunicar con APIs externas, parsear respuestas | Parámetros específicos | Datos tipados | Decisiones de negocio, mezclar datos |
| **Models** | Definir tipos de datos | - | - | Lógica, cálculos |
| **Frontend Routes** | Renderizar páginas | - | HTML + CSS + JS | Lógica de negocio, estado global sin store |
| **Frontend Components** | Mostrar UI, capturar interacciones | Props | Eventos | Llamadas API directas, estado no sincronizado |
| **Frontend Stores** | Centralizar estado reactivo | Acciones | Estado subscriptible | Lógica de presentación, validaciones de entrada |

### La regla de oro: Unidireccionalidad

```
Usuario
  ↓
Frontend (captura input)
  ↓
Backend - Handler (valida, delega)
  ↓
Backend - Service (orquesta)
  ↓
Backend - Integration (consulta)
  ↓
API Externa (retorna datos)
  ↓
Integration (parsea)
  ↓
Service (procesa)
  ↓
Handler (formatea)
  ↓
Frontend (muestra)
```

Cada capa **solo invoca la siguiente**, no salta ni regresa (salvo errores).

---

## Parte 1: Infraestructura Base — Hacer que Todo Compile

### Qué problema resuelve

Tener un proyecto que **compile sin errores**, con estructura inicial lista para agregar funcionalidad. Esto evita problemas de imports, dependencias rotas, o configuración faltante más adelante.

### Conceptos usados

- Inicialización de proyectos (módulos, dependencias).
- Estructura de carpetas.
- Configuración básica de servidor/cliente.

### Qué debería pensar antes de implementar

1. **Backend**: ¿Qué librería de router HTTP voy a usar? (stdlib, Gin, Echo, etc.). ¿Cómo voy a estructurar las rutas?
2. **Frontend**: ¿Voy a usar SSR o SPA? ¿Qué librerías de styling?
3. **Comunicación**: ¿En qué puerto corre el backend? ¿Cómo accede el frontend?

### Tareas de esta etapa

#### Tarea 1.1: Inicializar Backend

**¿Qué hacer?**

Crear un proyecto Go ejecutable que:
- Inicie un servidor HTTP en un puerto específico.
- Tenga rutas definidas (aunque todas retornen respuestas hardcoded por ahora).
- Tenga carpetas organizadas según la estructura descripta.

**¿Cómo verificar que funciona?**

```bash
# Desde la carpeta del backend
go run cmd/main.go
# Debería imprimir algo como: "Servidor en puerto 8080"
```

Luego, desde otra terminal:

```bash
curl http://localhost:8080/
# Debería recibir alguna respuesta (ej: 404, 200 con JSON vacío, etc)
```

**Errores comunes a evitar**:
- Olvidar inicializar el módulo Go (`go mod init`).
- Importar librerías sin instalarlas (`go get`).
- Portos en uso (cambiar puerto si necesario).
- No cerrar listeners correctamente.

---

#### Tarea 1.2: Definir Tipos de Datos

**¿Qué hacer?**

En la carpeta `models/`, crear structs Go que representen:
- Datos individuales que recibirás de la API externa (antes y después de parsear).
- Datos agregados que devolverás al frontend.

**¿Cómo pensar antes?**

Pregúntate:
- ¿Qué campos necesito del servicio externo?
- ¿En qué tipo Go se mapea cada campo? (string, int, float64, time.Time, etc).
- ¿Algunos campos son opcionales? (usar punteros o valores por defecto).
- ¿Cómo se llaman los campos en JSON? (necesitarás tags `json`).

**Ejemplo de lo que debes evitar**:

```go
// ❌ Mal: nombres genéricos, sin tags JSON
type Data struct {
    X string
    Y int
}

// ✔ Bien: nombres claros, con tags para JSON
type IndividualInfo struct {
    Attribute1 string  `json:"attribute_1"`
    Attribute2 float64 `json:"attribute_2"`
    Attribute3 int     `json:"attribute_3,omitempty"`
}
```

**Errores comunes a evitar**:
- Confundir nombres JSON con nombres Go.
- Olvidar que Go es case-sensitive y campos públicos deben comenzar con mayúscula.
- No considerar campos opcionales en la respuesta.

---

#### Tarea 1.3: Inicializar Frontend

**¿Qué hacer?**

Crear un proyecto SvelteKit que:
- Compila sin errores.
- Tiene una página inicial básica.
- Está configurado en modo SPA (sin SSR si lo prefieres).

**¿Cómo verificar que funciona?**

```bash
# Desde la carpeta del frontend
npm run dev
# Acceder a http://localhost:5173 (o el puerto que use)
# Debería ver la página inicial
```

**Errores comunes a evitar**:
- Usar paquetes incompatibles con tu versión de Node.
- Olvidar instalar dependencias (`npm install`).
- Mezclar Vite config con SvelteKit config.

---

### Checklist de Finalización (Parte 1)

- [ ] Backend compila y arranca sin panics.
- [ ] Frontend compila sin warnings.
- [ ] Backend responde a solicitudes HTTP básicas.
- [ ] Frontend muestra una página inicial.
- [ ] Las carpetas están organizadas según se describió.
- [ ] Los tipos de datos están definidos (aunque no se usen aún).

---

## Parte 2: Integración con Servicio Externo — Consultar Datos

### Qué problema resuelve

Poder **obtener datos reales** de un servicio externo. Esto es el corazón del sistema: sin datos externos, no hay nada que procesar.

### Conceptos usados

- Construcción de URLs dinámicas.
- HTTP client (request/response).
- Parsing de JSON.
- Manejo de errores básico.
- Timeouts de red.

### Qué debería pensar antes de implementar

1. **API Externa**: ¿Cuál es el endpoint exacto? ¿Qué parámetros requiere? ¿Qué formato tiene la respuesta?
2. **Errores**: ¿Qué pasa si la API no responde? ¿Si retorna un status code de error? ¿Si el JSON es malformado?
3. **Timeouts**: ¿Cuánto tiempo razonable esperar una respuesta? (típicamente 5-10 segundos).

### Tareas de esta etapa

#### Tarea 2.1: Función para Construir la Solicitud

**¿Qué hacer?**

Crear una función que:
- Reciba parámetros específicos (identificadores, coordenadas, etc).
- Construya una URL válida hacia el servicio externo.
- Retorne la URL o un error si los parámetros son inválidos.

**¿Cómo pensar antes?**

- ¿Los parámetros necesitan estar URL-encoded? (ej: espacios → %20).
- ¿Hay parámetros opcionales?
- ¿Cómo valido que los parámetros sean razonables? (ej: si espero un número, ¿rechazar negativos?).

**Ejemplo de forma, no de contenido**:

```go
// Función que construye una solicitud
func buildRequestToExternalService(param1 string, param2 int) (string, error) {
    // Validar parámetros
    // Construir URL
    // Retornar URL o error
    return "", nil
}

// Uso
url, err := buildRequestToExternalService("example", 42)
if err != nil {
    // Manejar error
}
```

**Errores comunes a evitar**:
- No escapar parámetros en la URL.
- Hardcodear valores que deberían venir como parámetros.
- No validar parámetros (ej: aceptar valores negativos cuando solo tiene sentido positivos).

---

#### Tarea 2.2: Ejecutar la Solicitud HTTP

**¿Qué hacer?**

Crear una función que:
- Tome una URL construida.
- Ejecute un request HTTP (GET, POST, etc).
- Maneje timeouts.
- Chequee el status code (200 es éxito, otros son errores).
- Retorne la respuesta en crudo (como string o []byte).

**¿Cómo pensar antes?**

- ¿Qué method HTTP debo usar? (GET, POST, etc).
- ¿Necesito headers especiales? (User-Agent, Authorization, etc).
- ¿Cómo sé que la respuesta es válida más allá del status code?

**Ejemplo de forma**:

```go
// Función que ejecuta HTTP request
func executeHTTPRequest(ctx context.Context, url string) ([]byte, error) {
    // Crear request con contexto
    // Ejecutar request
    // Chequear status code
    // Si no es éxito, retornar error
    // Leer body
    // Retornar body o error
    return nil, nil
}

// Uso
body, err := executeHTTPRequest(ctx, "https://example.com/data")
if err != nil {
    // Manejar error
}
```

**Errores comunes a evitar**:
- No usar context para timeout.
- Olvidar cerrar el body del response (`response.Body.Close()`).
- Asumir que status 200 significa que el JSON es válido.
- No manejar redirects (¿debería seguirlos?).

---

#### Tarea 2.3: Parsear la Respuesta

**¿Qué hacer?**

Crear una función que:
- Reciba los datos en crudo ([]byte o string).
- Use el paquete `encoding/json` para convertir a struct.
- Valide que los tipos sean correctos.
- Retorne un struct tipado o error si el JSON es inválido.

**¿Cómo pensar antes?**

- ¿Todos los campos en la respuesta son necesarios o algunos son opcionales?
- ¿Los nombres en JSON coinciden con lo que defini en Models?
- ¿Hay anidamiento en el JSON? (ej: `{ "data": { "value": 42 } }`).

**Ejemplo de forma**:

```go
// Función que parsea JSON
func parseExternalResponse(rawData []byte, targetStruct interface{}) error {
    // Usar json.Unmarshal para convertir []byte a struct
    // Si hay error, retornar error descriptivo
    // Si éxito, retornar nil
    return nil
}

// Uso
var result SomeDataType
err := parseExternalResponse(body, &result)
if err != nil {
    // Manejar error de parsing
}
```

**Errores comunes a evitar**:
- Los tags `json` en el struct no coinciden con los nombres en JSON.
- Confundir punteros con valores (algunos campos deben ser `*Type` si son opcionales).
- No validar valores después de parsear (ej: valores nulos inesperados).

---

#### Tarea 2.4: Prueba Manual (sin Backend Aún)

**¿Qué hacer?**

Crear un pequeño programa o función `main` que:
- Llame directamente a la función de consulta.
- Imprima los resultados.
- Pueda ser ejecutado desde línea de comandos.

**¿Cómo verificar que funciona?**

```bash
go run test_external_call.go
# Debería imprimir datos reales del servicio externo
```

**Errores comunes a evitar**:
- Olvidar usar `context.Background()` si no tienes timeout aún.
- No manejar el error de forma legible (imprimir `fmt.Println(err)` no siempre ayuda).

---

### Checklist de Finalización (Parte 2)

- [ ] Puedo construir una URL válida.
- [ ] Puedo ejecutar un request HTTP y recibir respuesta.
- [ ] Puedo parsear la respuesta a un struct.
- [ ] He probado manualmente que los datos son correctos.
- [ ] He identificado qué errores pueden ocurrir y cómo reportarlos.

---

## Parte 3: Lógica de Dominio — Procesar un Solo Dato

### Qué problema resuelve

Antes de paralelizar, es crucial que **la lógica sea correcta**. Esta etapa construye la lógica de procesamiento de forma secuencial, lo que permite medir un baseline de tiempo y verificar que los cálculos y agregaciones sean exactos.

### Conceptos usados

- Iteración sobre colecciones.
- Agregación de datos (cálculos, extremos, rankings).
- Transformación de tipos.

### Qué debería pensar antes de implementar

1. **Datos de entrada**: ¿Qué estructura tienen los datos crudos? ¿Cómo se convierten a lo que necesito?
2. **Cálculos**: ¿Qué operaciones debo hacer? (promedios, máximos, mínimos, conteos).
3. **Estructura de salida**: ¿Cómo debo devolver los resultados para que sean útiles al handler HTTP?

### Tareas de esta etapa

#### Tarea 3.1: Función para Procesar un Dato Individual

**¿Qué hacer?**

Crear una función que:
- Reciba un dato individual de la API externa.
- Aplique transformaciones o cálculos específicos.
- Retorne un dato procesado en una estructura interna.

**Ejemplo de forma**:

```go
// Función que procesa UN dato
func processSingleItem(rawData ExternalDataType) (ProcessedDataType, error) {
    // Validar que los datos tienen sentido
    // Transformar/calcular lo necesario
    // Retornar dato procesado
    return ProcessedDataType{}, nil
}
```

**Errores comunes a evitar**:
- Asumir que todos los campos del dato externo son válidos.
- Olvidar validar rangos razonables (ej: temperatura < -100°C probablemente es error).
- Mezclar procesamiento con agregación (esta función debe ser "pura").

---

#### Tarea 3.2: Función de Orquestación Secuencial

**¿Qué do?**

Crear una función que:
- Reciba una lista de identificadores/parámetros.
- Para cada uno, llame a la función de consulta de la Parte 2.
- Para cada respuesta, llame a la función de procesamiento (3.1).
- Acumule todos los resultados.
- Retorne todos los datos procesados (y errores si los hay).

**¿Cómo pensar antes?**

- ¿Qué pasa si una consulta falla? ¿Abortar todo o continuar con las otras?
- ¿Cómo reporto qué falló y qué no?
- ¿En qué orden proceso? ¿Importa?

**Ejemplo de forma**:

```go
// Función que procesa MÚLTIPLES datos secuencialmente
func orchestrateSequential(ctx context.Context, ids []string) ([]ProcessedDataType, []error, error) {
    results := []ProcessedDataType{}
    errors := []error{}
    
    for _, id := range ids {
        // Consultar dato externo
        // Procesar
        // Si error, agregar a errors (pero continuar)
        // Si éxito, agregar a results
    }
    
    return results, errors, nil
}
```

**Errores comunes a evitar**:
- Abortar la iteración cuando un elemento falla (manejo de errores parciales).
- No distinguir entre "error en un elemento" y "error total en la orquestación".
- Olvidar medir el tiempo de ejecución (necesario para comparar con versión concurrente).

---

#### Tarea 3.3: Agregación de Resultados

**¿Qué hacer?**

Crear una función que:
- Reciba todos los datos procesados.
- Calcule agregaciones (promedios, máximos, mínimos, rankings, conteos, agrupaciones).
- Retorne una estructura con los datos agregados.

**¿Cómo pensar antes?**

- ¿Qué comparaciones debo hacer? (ej: cuál es mayor).
- ¿Cómo agrupo datos? (ej: por categoría).
- ¿Cómo rankeo? (ascendente, descendente).
- ¿Qué pasa si la lista está vacía?

**Ejemplo de forma**:

```go
// Función que agrega resultados
func aggregateResults(processedData []ProcessedDataType) (AggregateResult, error) {
    if len(processedData) == 0 {
        return AggregateResult{}, errors.New("no data to aggregate")
    }
    
    aggregate := AggregateResult{}
    
    // Iterar sobre processedData para calcular:
    // - Máximos/mínimos
    // - Promedios
    // - Rankings
    // - Agrupaciones
    
    return aggregate, nil
}

// Uso
agg, err := aggregateResults(allProcessedData)
```

**Errores comunes a evitar**:
- División por cero (lista vacía).
- Comparar tipos incompatibles.
- No manejar valores nulos/inválidos en los datos.

---

#### Tarea 3.4: Función de Composición

**¿Qué hacer?**

Crear una función que combine 3.2 y 3.3:
- Orquestar la consulta y procesamiento de múltiples datos.
- Agregar los resultados.
- Retornar tanto datos individuales como agregados.

**Ejemplo de forma**:

```go
// Función que une todo
func processAndAggregate(ctx context.Context, ids []string) ([]ProcessedDataType, AggregateResult, error) {
    // Llamar a orchestrateSequential
    // Llamar a aggregateResults
    // Retornar ambos conjuntos de datos
    return nil, AggregateResult{}, nil
}
```

**Errores comunes a evitar**:
- Perder información sobre qué falló.
- No propagar errores correctamente.

---

#### Tarea 3.5: Prueba Manual de Lógica Completa

**¿Qué hacer?**

Crear un programa de prueba que:
- Llame a `processAndAggregate` con IDs hardcoded.
- Imprima los resultados individuales y agregados.
- Mida el tiempo de ejecución.

**¿Cómo verificar que funciona?**

```bash
go run test_logic.go
# Debería imprimir:
# - Datos individuales procesados
# - Datos agregados
# - Tiempo total
```

**Lo importante**: Anotar el tiempo. Será tu baseline para comparar con la versión concurrente.

---

### Checklist de Finalización (Parte 3)

- [ ] Puedo procesar un dato individual.
- [ ] Puedo procesar múltiples datos secuencialmente.
- [ ] Puedo agregar datos en una estructura de salida.
- [ ] He medido el tiempo de ejecución.
- [ ] Manejo errores parciales (si una falla, las otras continúan).
- [ ] La lógica es correcta con datos de prueba.

---

## Parte 4: Concurrencia — Paralelizar Consultas

### Qué problema resuelve

Si la Parte 3 tardó N segundos, y hiciste M consultas, entonces tardó ~M * (tiempo por consulta). Con concurrencia, debería tardar ~máximo(tiempo por consulta). Este es el salto de performance importante.

### Conceptos usados

- Goroutines (crear operaciones paralelas).
- Channels (comunicación entre goroutines).
- Patrón fan-out (lanzar N operaciones).
- Patrón fan-in (recolectar N resultados).
- sync.WaitGroup (sincronización).
- context.Context (cancelación, timeout).

### Qué debería pensar antes de implementar

1. **Paralelismo**: ¿Cuántas goroutines debo lanzar? (típicamente, una por elemento).
2. **Channel design**: ¿Qué tipo de dato envío por el channel? (¿resultado, error, o ambos?).
3. **Sincronización**: ¿Cómo sé que todas las goroutines terminaron?
4. **Errores**: ¿Si una goroutine falla, aborto las otras o continúo?

### Tareas de esta etapa

#### Tarea 4.1: Struct para Encapsular Resultado + Error

**¿Qué hacer?**

Crear un struct que contenga tanto el resultado exitoso como el error de una operación. Esto permite enviar ambos por el channel sin complicaciones.

**¿Cómo pensar antes?**

- ¿Puede haber tanto resultado como error? (típicamente, uno u otro).
- ¿Cómo serializo esto si necesito enviarlo?

**Ejemplo de forma**:

```go
// Struct que encapsula resultado O error
type ResultWrapper struct {
    Data  ProcessedDataType
    Error error
}
```

**Errores comunes a evitar**:
- Enviar puntero nil sin validar en el receptor.
- No inicializar los campos correctamente.

---

#### Tarea 4.2: Función que Consulta + Procesa en Goroutine

**¿Qué hacer?**

Refactorizar la lógica de la Parte 3 para que:
- Reciba un ID y un channel.
- Consulte el dato externo.
- Lo procese.
- Envíe el resultado (exitoso o con error) al channel.
- Retorne (la goroutine finaliza).

**¿Cómo pensar antes?**

- ¿Necesito context aquí para respetar timeouts?
- ¿Qué información envío al channel si falla?

**Ejemplo de forma**:

```go
// Función que ejecuta en una goroutine
func fetchAndProcessInGoroutine(ctx context.Context, id string, resultChan chan ResultWrapper) {
    // Consultar
    // Procesar
    // Enviar resultado al channel (éxito o error)
}

// Uso (desde otra función)
go fetchAndProcessInGoroutine(ctx, "id1", resultChan)
```

**Errores comunes a evitar**:
- Olvidar cerrar la goroutine (no debe quedar bloqueada).
- No respetar el context (cancelación de timeout).
- Enviar al channel sin chequear si está cerrado.

---

#### Tarea 4.3: Patrón Fan-Out (Lanzar Múltiples Goroutines)

**¿Qué hacer?**

Crear una función que:
- Cree un channel.
- Para cada ID, lance una goroutine (4.2) que escriba en el channel.
- Retorne el channel y algún mecanismo para saber cuándo esperar.

**¿Cómo pensar antes?**

- ¿Unbuffered o buffered channel? (buffered es más eficiente si tienes muchas goroutines).
- ¿Cómo sé cuántas goroutines lancé? (necesito contar para luego esperar).

**Ejemplo de forma**:

```go
// Función que lanza múltiples goroutines (fan-out)
func fanOut(ctx context.Context, ids []string) chan ResultWrapper {
    resultChan := make(chan ResultWrapper, len(ids)) // buffered
    
    for _, id := range ids {
        go fetchAndProcessInGoroutine(ctx, id, resultChan)
    }
    
    return resultChan
}
```

**Errores comunes a evitar**:
- Cerrar el channel antes de que todas las goroutines escriban (panic).
- Unbuffered channel con muchas goroutines (puede causar deadlock).

---

#### Tarea 4.4: Patrón Fan-In (Recolectar Múltiples Resultados)

**¿Qué hacer?**

Crear una función que:
- Reciba un channel.
- Sepa cuántos resultados esperar.
- Lea del channel hasta recibir todos.
- Agrupe éxitos y errores.
- Retorne listas separadas.

**¿Cómo pensar antes?**

- ¿De dónde sé cuántos resultados esperar?
- ¿Qué pasa si una goroutine no envía nada? (deadlock).
- ¿Cómo evito leer indefinidamente?

**Ejemplo de forma**:

```go
// Función que recolecta resultados (fan-in)
func fanIn(resultChan chan ResultWrapper, expectedCount int) ([]ProcessedDataType, []error) {
    results := []ProcessedDataType{}
    errors := []error{}
    
    for i := 0; i < expectedCount; i++ {
        result := <-resultChan  // Bloquea hasta recibir
        if result.Error != nil {
            errors = append(errors, result.Error)
        } else {
            results = append(results, result.Data)
        }
    }
    
    return results, errors
}
```

**Errores comunes a evitar**:
- Deadlock si `expectedCount` no coincide con goroutines lanzadas.
- Leer del channel después de cerrarlo (panic).

---

#### Tarea 4.5: Sincronización con WaitGroup (Alternativa)

**¿Qué hacer?**

Alternativa a contar resultados: usar `sync.WaitGroup` para sincronizar.

**¿Cómo pensar antes?**

- ¿Prefiero contar o usar WaitGroup?
- WaitGroup es más explícito si tengo múltiples operaciones.

**Ejemplo de forma**:

```go
// Función que usa WaitGroup para sincronizar
func orchestrateConcurrentWithWaitGroup(ctx context.Context, ids []string) ([]ProcessedDataType, []error) {
    resultChan := make(chan ResultWrapper, len(ids))
    var wg sync.WaitGroup
    
    // Lanzar goroutines
    for _, id := range ids {
        wg.Add(1)
        go func(id string) {
            defer wg.Done()
            // Consultar, procesar, enviar
        }(id)
    }
    
    // En otra goroutine, esperar a que terminen todas y cerrar channel
    go func() {
        wg.Wait()
        close(resultChan)
    }()
    
    // Leer resultados hasta que channel se cierre
    results := []ProcessedDataType{}
    errors := []error{}
    for result := range resultChan {
        // Procesar
    }
    
    return results, errors
}
```

**Errores comunes a evitar**:
- Olvidar `defer wg.Done()`.
- Cerrar channel antes de que las goroutines terminen.
- No leer del channel correctamente.

---

#### Tarea 4.6: Context y Timeout

**¿Qué hacer?**

Agregar timeout a nivel global:
- Si alguna goroutine tarda demasiado, cancelar todas.
- Reportar qué se canceló.

**¿Cómo pensar antes?**

- ¿Cuál es un timeout razonable? (ej: 30 segundos para un conjunto de consultas).
- ¿Las goroutines chequean `ctx.Done()` para salir rápido?

**Ejemplo de forma**:

```go
// Función que envuelve con timeout
func orchestrateConcurrentWithTimeout(baseCtx context.Context, ids []string, timeout time.Duration) ([]ProcessedDataType, []error) {
    ctx, cancel := context.WithTimeout(baseCtx, timeout)
    defer cancel()
    
    // Usar ctx en lugar de baseCtx en todas las goroutines
    // Las goroutines deben chequear ctx.Done() o ctx.Err()
    
    return orchestrateConcurrent(ctx, ids)
}
```

**Errores comunes a evitar**:
- Olvidar `defer cancel()`.
- Las goroutines no chequean el contexto.
- Timeout muy bajo (muchos falsos positivos).

---

#### Tarea 4.7: Prueba Comparativa

**¿Qué hacer?**

Crear un programa de prueba que:
- Ejecute la versión secuencial (Parte 3) y mida tiempo.
- Ejecute la versión concurrente (esta parte) y mida tiempo.
- Compare y imprima el resultado.

**¿Cómo verificar que funciona?**

```bash
go run test_concurrent.go
# Debería imprimir algo como:
# Secuencial: 10s
# Concurrente: 2s
# Speedup: 5x
```

**Lo importante**: Los **resultados deben ser idénticos** (solo debe cambiar el tiempo). Si los datos son diferentes, hay un bug.

---

### Checklist de Finalización (Parte 4)

- [ ] Puedo lanzar múltiples goroutines sin panics.
- [ ] Puedo recolectar todos los resultados correctamente.
- [ ] No hay deadlock ni race conditions.
- [ ] El tiempo de ejecución es significativamente menor que secuencial.
- [ ] Los resultados son idénticos a la versión secuencial.
- [ ] Manejo timeouts sin colgar.
- [ ] Manejo errores parciales (si una goroutine falla, las otras continúan).

---

## Parte 5: Exposición HTTP — Crear Endpoints

### Qué problema resuelve

La lógica funciona, pero está aislada. Ahora necesita ser **accesible vía HTTP** para que el frontend (y herramientas como `curl`) puedan usarla.

### Conceptos usados

- HTTP handlers (recibir solicitudes, escribir respuestas).
- Parsing de JSON en entrada.
- Serialización de JSON en salida.
- Códigos de estado HTTP.
- Validación de entrada.

### Qué debería pensar antes de implementar

1. **Diseño de rutas**: ¿Qué endpoints necesito? ¿Qué información reciben?
2. **Métodos HTTP**: ¿GET, POST, PUT?
3. **Validación**: ¿Qué valido en el handler?
4. **Respuestas de error**: ¿Qué codes devuelvo para cada caso?

### Tareas de esta etapa

#### Tarea 5.1: Endpoint para Datos Individuales (GET)

**¿Qué hacer?**

Crear un handler HTTP que:
- Reciba un ID como parámetro (en la ruta o query).
- Valide que el ID es válido.
- Llame a la función de consulta y procesamiento (Parte 3).
- Retorne el resultado como JSON.

**¿Cómo pensar antes?**

- ¿Dónde va el ID? (`:id` en la ruta o `?id=...` en query).
- ¿Cómo extraigo el parámetro del request?
- ¿Qué status code devuelvo si no hay error? (200 OK).
- ¿Qué si el ID no existe? (404 Not Found? 400 Bad Request?).

**Ejemplo de forma**:

```go
// Handler para un dato individual
func handlerForIndividualData(w http.ResponseWriter, r *http.Request) {
    // Extraer parámetro de la ruta/query
    // Validar parámetro
    // Llamar a función de consulta
    // Si error, escribir status 400 o 500 y mensaje de error
    // Si éxito, escribir status 200 y JSON del resultado
}

// En main.go:
// router.HandleFunc("/api/path/{param}", handlerForIndividualData).Methods("GET")
```

**Errores comunes a evitar**:
- No validar entrada (ej: ID vacío).
- Confundir 400 (mala solicitud del cliente) con 500 (error del servidor).
- Olvidar escribir header `Content-Type: application/json`.

---

#### Tarea 5.2: Endpoint para Comparación (POST)

**¿Qué hacer?**

Crear un handler HTTP que:
- Reciba una lista de IDs en el body (JSON).
- Valide la lista (no vacía, IDs válidos).
- Llame a la función de orquestación concurrente (Parte 4).
- Retorne datos individuales y agregados como JSON.

**¿Cómo pensar antes?**

- ¿Cómo parseo el JSON del body?
- ¿Qué estructura espero en el JSON?
- ¿Cómo valido la lista? (vacía, duplicados, tamaño máximo).
- ¿Qué status code si la lista es inválida? (400).

**Ejemplo de forma**:

```go
// Struct para parsear el body del request
type RequestData struct {
    IDs []string `json:"ids"`
}

// Handler para comparación
func handlerForComparison(w http.ResponseWriter, r *http.Request) {
    // Parsear body JSON
    // Validar que se parseó correctamente
    // Validar que los IDs son válidos
    // Llamar a función de orquestación
    // Construir respuesta (datos individuales + agregados)
    // Escribir JSON de respuesta
}
```

**Errores comunes a evitar**:
- No usar `json.NewDecoder` (leer body correctamente).
- Olvidar validar después de parsear.
- No cerrar el body del request.
- Enviar datos parciales si hay error.

---

#### Tarea 5.3: Estructura de Respuesta

**¿Qué hacer?**

Definir qué estructura devolverá cada endpoint:
- Datos individuales.
- Datos agregados.
- Errores parciales (si aplica).

**¿Cómo pensar antes?**

- ¿El frontend necesita saber qué falló específicamente?
- ¿Devuelvo ambos (datos y errores) o aborto en error?

**Ejemplo de forma**:

```go
// Struct para respuesta HTTP
type ResponseData struct {
    Success         bool                   `json:"success"`
    IndividualData  []ProcessedDataType    `json:"individual_data"`
    AggregateData   AggregateResult        `json:"aggregate_data"`
    PartialErrors   []string               `json:"partial_errors,omitempty"`
    ErrorMessage    string                 `json:"error_message,omitempty"`
}
```

**Errores comunes a evitar**:
- Estructuras inconsistentes (a veces error está en field, a veces en JSON raíz).
- Demasiados niveles de anidamiento.

---

#### Tarea 5.4: Códigos de Estado HTTP

**¿Qué hacer?**

Decidir y documentar qué status code devuelves en cada caso:
- 200 OK: Éxito total.
- 400 Bad Request: Entrada inválida.
- 404 Not Found: Recurso no existe.
- 500 Internal Server Error: Error del servidor.
- 503 Service Unavailable: Servicio externo no disponible.

**¿Cómo pensar antes?**

- ¿Distingo entre error del cliente vs. error del servidor?
- ¿Hay situaciones especiales? (ej: servicio externo caído → 503).

**Errores comunes a evitar**:
- Siempre devolver 200 (difícil debuggear).
- Devolver 500 cuando es 400 (confunde al cliente).

---

#### Tarea 5.5: Prueba Manual con curl

**¿Qué hacer?**

Testear cada endpoint manualmente:

```bash
# GET individual
curl -X GET http://localhost:8080/api/path/someId

# POST comparación
curl -X POST http://localhost:8080/api/path/compare \
  -H "Content-Type: application/json" \
  -d '{"ids": ["id1", "id2", "id3"]}'
```

**¿Cómo verificar que funciona?**

- Deberías recibir JSON válido.
- Los códigos de status deberían ser correctos (200 para éxito, 400 para input inválido).
- Prueba con entrada inválida y verifica que se rechace apropiadamente.

---

### Checklist de Finalización (Parte 5)

- [ ] Hay un endpoint para datos individuales.
- [ ] Hay un endpoint para comparación (concurrente).
- [ ] Ambos endpoints retornan JSON válido.
- [ ] Los códigos de status HTTP son correctos.
- [ ] La validación de entrada funciona.
- [ ] `curl` puede consumir los endpoints sin error.

---

## Parte 6: Frontend Base — Estructura y Componentes

### Qué problema resuelve

El backend está listo, pero el usuario no puede usarlo sin un navegador. Esta parte construye la interfaz: selector de opciones, visualización de resultados.

### Conceptos usados

- Componentes Svelte (props, eventos).
- Reactividad (bindings, stores).
- Renderizado condicional (if, each).
- Estilos CSS.

### Qué debería pensar antes de implementar

1. **Layout**: ¿Cómo se organizan los elementos? (selector arriba, resultados abajo).
2. **Componentes**: ¿Qué reutilizar? (lista de opciones, card de resultado).
3. **Estado**: ¿Qué es local vs. global?

### Tareas de esta etapa

#### Tarea 6.1: Componente Selector

**¿Qué hacer?**

Crear un componente que:
- Reciba una lista de opciones (como props).
- Permita seleccionar una o múltiples.
- Emita eventos cuando el usuario selecciona/deselecciona.

**¿Cómo pensar antes?**

- ¿Checkboxes o radio buttons?
- ¿Qué información tiene cada opción? (solo ID, o también nombre/descripción).
- ¿Cómo comunico la selección al padre?

**Ejemplo de forma**:

```svelte
<!-- Componente que recibe lista y emite eventos -->
<script>
  export let options = [];
  
  let selected = {};
  
  function toggle(id) {
    selected[id] = !selected[id];
    // Emitir evento hacia el padre
    // dispatch('selectionChanged', selected);
  }
</script>

<div>
  {#each options as option (option.id)}
    <label>
      <input 
        type="checkbox" 
        on:change={() => toggle(option.id)}
      />
      {option.name}
    </label>
  {/each}
</div>
```

**Errores comunes a evitar**:
- Two-way binding sin `bind:` puede causar desincronización.
- Olvidar emitir eventos (el padre no sabe qué cambió).

---

#### Tarea 6.2: Componente de Resultado Individual

**¿Qué hacer?**

Crear un componente que:
- Reciba un dato individual procesado (como props).
- Lo visualice en forma clara (ej: card con campos).

**¿Cómo pensar antes?**

- ¿Qué campos mostrar?
- ¿Cómo organizo visualmente? (grid, flexbox).
- ¿Necesito iconos o colores para mejorar UI?

**Ejemplo de forma**:

```svelte
<!-- Componente que muestra UN resultado -->
<script>
  export let data;
</script>

<div class="card">
  <h3>{data.name}</h3>
  <p>Attribute 1: {data.attribute1}</p>
  <p>Attribute 2: {data.attribute2}</p>
</div>

<style>
  .card {
    border: 1px solid #ccc;
    padding: 1rem;
    border-radius: 4px;
  }
</style>
```

**Errores comunes a evitar**:
- No manejar datos nulos/undefined.
- Suponer estructura de datos (usa props con tipos claros si es posible).

---

#### Tarea 6.3: Componente de Resumen Agregado

**¿Qué hacer?**

Crear un componente que:
- Reciba datos agregados (como props).
- Visualice extremos, promedios, rankings, etc.

**¿Cómo pensar antes?**

- ¿Qué destacar visualmente? (máximo, mínimo).
- ¿Cómo mostrar una lista ordenada (ranking)?

**Errores comunes a evitar**:
- No validar si existen datos antes de mostrarlos.

---

#### Tarea 6.4: Página Principal

**¿Qué hacer?**

Ensamblar todos los componentes en una página:
- Selector arriba.
- Botón "Comparar".
- Área para resultados individuales.
- Área para resumen.

**¿Cómo pensar antes?**

- ¿Cómo organizo espacialmente?
- ¿Cómo paso datos entre componentes? (props hacia abajo, eventos hacia arriba).

**Ejemplo de forma**:

```svelte
<script>
  import Selector from './components/Selector.svelte';
  import ResultsGrid from './components/ResultsGrid.svelte';
  import Summary from './components/Summary.svelte';
  
  let options = [];
  let selected = {};
  let results = null;
  let loading = false;
  
  function handleCompare() {
    // Preparar datos seleccionados
    // Llamar API
    // Actualizar results
  }
</script>

<div class="container">
  <Selector {options} on:selectionChanged={...} />
  <button on:click={handleCompare}>Compare</button>
  
  {#if loading}
    <p>Loading...</p>
  {/if}
  
  {#if results}
    <ResultsGrid data={results.individual} />
    <Summary data={results.aggregate} />
  {/if}
</div>
```

**Errores comunes a evitar**:
- Props no se sincronizan automáticamente (necesitas eventos o stores).
- Olvidar indicador de carga.

---

### Checklist de Finalización (Parte 6)

- [ ] Hay un componente para seleccionar opciones.
- [ ] Hay componentes para visualizar resultados.
- [ ] La página principal ensambla todo sin errores.
- [ ] Los componentes se renderizan sin datos (sin crashes).

---

## Parte 7: Integración Frontend-Backend — Conectar

### Qué problema resuelve

El frontend existe, el backend existe, pero no se hablan. Esta parte es el puente.

### Conceptos usados

- HTTP requests desde el navegador (fetch).
- Manejo de promises/async-await.
- Actualización de estado cuando llegan respuestas.
- Manejo de errores y estados de carga.

### Qué debería pensar antes de implementar

1. **Cliente HTTP**: ¿Uso `fetch` directo o una librería?
2. **CORS**: ¿El backend permite requests desde el frontend?
3. **Estado**: ¿Dónde centralizo el estado? (store global vs. componente local).

### Tareas de esta etapa

#### Tarea 7.1: Cliente HTTP (Función Genérica)

**¿Qué hacer?**

Crear una función reutilizable que:
- Ejecute requests HTTP.
- Maneje timeouts.
- Parsee respuestas.
- Reporte errores.

**¿Cómo pensar antes?**

- ¿Qué métodos HTTP necesito? (GET, POST).
- ¿Cómo manejo headers?
- ¿Qué pasa si response no es JSON válido?

**Ejemplo de forma**:

```javascript
// Función genérica para HTTP requests
async function apiCall(method, endpoint, data = null) {
  const options = {
    method,
    headers: { 'Content-Type': 'application/json' },
  };
  
  if (data) {
    options.body = JSON.stringify(data);
  }
  
  const response = await fetch(`http://localhost:8080${endpoint}`, options);
  
  if (!response.ok) {
    throw new Error(`HTTP ${response.status}`);
  }
  
  return await response.json();
}

// Uso
const result = await apiCall('POST', '/api/compare', { ids: [...] });
```

**Errores comunes a evitar**:
- Olvidar `Content-Type` header.
- No chequear `response.ok`.
- Asumir que JSON es válido sin validar.

---

#### Tarea 7.2: Store para Estado Global

**¿Qué hacer?**

Crear un store que:
- Centralice el estado (selecciones, resultados, carga).
- Permita que múltiples componentes se suscriban.

**¿Cómo pensar antes?**

- ¿Qué datos son globales? (selecciones, resultados sí; cosas locales de un componente no).
- ¿Cómo estructura el estado para que sea reactivo?

**Ejemplo de forma**:

```javascript
// store.js
import { writable } from 'svelte/store';

export const appState = writable({
  options: [],
  selected: {},
  results: null,
  loading: false,
  error: null,
});

// En componente
import { appState } from './store.js';

$: state = $appState;
```

**Errores comunes a evitar**:
- Mutar el store sin usar `.update()` o `.set()`.
- Olvidar suscribirse con `$` (Svelte syntax).

---

#### Tarea 7.3: Función para Cargar Opciones

**¿Qué hacer?**

Crear una función que:
- Llame al endpoint que lista opciones disponibles.
- Actualice el store.
- Maneje errores.

**¿Cómo pensar antes?**

- ¿Cuándo llamo esto? (en `onMount`).
- ¿Cómo muestro indicador de carga?

**Ejemplo de forma**:

```javascript
async function loadOptions() {
  appState.update(s => ({ ...s, loading: true }));
  try {
    const data = await apiCall('GET', '/api/options');
    appState.update(s => ({ ...s, options: data, error: null }));
  } catch (err) {
    appState.update(s => ({ ...s, error: err.message }));
  } finally {
    appState.update(s => ({ ...s, loading: false }));
  }
}
```

**Errores comunes a evitar**:
- Olvidar limpiar el estado anterior.
- No manejar todas las tres fases (loading, success, error).

---

#### Tarea 7.4: Función para Enviar Comparación

**¿Qué hacer?**

Crear una función que:
- Tome las selecciones del usuario.
- Las envíe al backend.
- Reciba resultados.
- Actualice el store.

**¿Cómo pensar antes?**

- ¿Qué estructura envío? (lista de IDs, objetos con detalles).
- ¿Cómo manejo errores parciales vs. totales?

**Ejemplo de forma**:

```javascript
async function submitComparison() {
  const selectedIds = Object.keys($appState.selected).filter(id => $appState.selected[id]);
  
  if (selectedIds.length === 0) {
    appState.update(s => ({ ...s, error: 'Select at least one option' }));
    return;
  }
  
  appState.update(s => ({ ...s, loading: true, error: null }));
  
  try {
    const result = await apiCall('POST', '/api/compare', { ids: selectedIds });
    appState.update(s => ({ ...s, results: result }));
  } catch (err) {
    appState.update(s => ({ ...s, error: err.message }));
  } finally {
    appState.update(s => ({ ...s, loading: false }));
  }
}
```

**Errores comunes a evitar**:
- No validar selecciones (enviar lista vacía).
- No manejar respuesta del servidor correctamente.

---

#### Tarea 7.5: Conectar Componentes con Funciones

**¿Qué hacer?**

Modificar los componentes para:
- Llamar a `loadOptions()` en `onMount`.
- Llamar a `submitComparison()` cuando el usuario presiona botón.
- Mostrar indicadores de carga y errores.

**Ejemplo de forma**:

```svelte
<script>
  import { onMount } from 'svelte';
  import { appState } from './store.js';
  import { loadOptions, submitComparison } from './api.js';
  
  onMount(() => {
    loadOptions();
  });
</script>

<Selector options={$appState.options} ... />
<button on:click={submitComparison}>Compare</button>

{#if $appState.loading}
  <p>Loading...</p>
{/if}

{#if $appState.error}
  <p class="error">{$appState.error}</p>
{/if}

{#if $appState.results}
  <ResultsGrid data={$appState.results.individual} />
{/if}
```

**Errores comunes a evitar**:
- No chequear `$appState` (el $ es syntax sugar de Svelte para suscripción).
- Olvidar `onMount` (loaded data aparece después de renderizar).

---

#### Tarea 7.6: CORS si es necesario

**¿Qué hacer?**

Si el frontend está en `http://localhost:5173` y el backend en `http://localhost:8080`, configurar CORS en el backend para permitir requests.

**¿Cómo pensar antes?**

- ¿Frontend y backend en el mismo origen? (raro en desarrollo local).
- ¿Qué headers de CORS necesito? (típicamente `Access-Control-Allow-Origin: *` para desarrollo).

**Errores comunes a evitar**:
- Olvidar CORS (browser bloquea la request silenciosamente).
- CORS demasiado permisivo en producción.

---

#### Tarea 7.7: Prueba Manual End-to-End

**¿Qué hacer?**

Usar la aplicación como usuario final:
1. Frontend carga opciones.
2. Seleccionar algunas.
3. Presionar comparar.
4. Ver resultados en pantalla.

**¿Cómo verificar que funciona?**

- Abre el navegador en el frontend.
- Inspecciona la Network tab en DevTools para ver requests.
- Verifica que datos llegan correctamente.

---

### Checklist de Finalización (Parte 7)

- [ ] El cliente HTTP funciona (fetch wrapper).
- [ ] El store se actualiza cuando cambian datos.
- [ ] Las opciones se cargan automáticamente.
- [ ] La comparación se envía y se reciben resultados.
- [ ] Se muestran indicadores de carga y errores.
- [ ] CORS está configurado (si necesario).

---

## Parte 8: Pulido y Refinamientos

### Qué problema resuelve

La aplicación funciona, pero puede ser más robusta y usable. Esta parte agrega:
- Validaciones más detalladas.
- Mensajes de error más claros.
- UI más pulida.
- Performance minor optimizations.

### Conceptos usados

- Validación avanzada.
- Manejo de edge cases.
- Estilos CSS refinados.
- Logging (para debugging).

### Tareas de esta etapa

#### Tarea 8.1: Validación Robusta

**¿Qué hacer?**

Agregar validaciones en múltiples niveles:
- Backend: rechazar IDs duplicados, listas vacías, IDs inválidos.
- Frontend: deshabilitar botón si no hay selecciones, indicar campos obligatorios.

**¿Cómo pensar antes?**

- ¿Qué casos extremos puedo encontrar? (vacío, muy grande, inválido).
- ¿Cómo evito que el usuario cometa errores?

---

#### Tarea 8.2: Mensajes de Error Descriptivos

**¿Qué hacer?**

Reemplazar errores genéricos con mensajes útiles:
- ❌ "Error" → ✔ "La API no respondió en 30 segundos, intenta nuevamente"
- ❌ "400" → ✔ "Selecciona al menos una opción antes de comparar"

**¿Cómo pensar antes?**

- ¿Quién leerá el mensaje? (usuario final, no desarrollador).
- ¿Puedo indicar qué hacer?

---

#### Tarea 8.3: Mejoras Visuales

**¿Qué hacer?**

Pulir UI:
- Responsive design (mobile-friendly).
- Colores, tipografía consistente.
- Spacing y alignment adecuados.
- Hover states en botones.

**¿Cómo pensar antes?**

- ¿Cómo se ve en mobile?
- ¿Hay suficiente contraste para accesibilidad?

---

#### Tarea 8.4: Logging para Debugging

**¿Qué hacer?**

Agregar logs útiles (sin abrumar):
- Backend: qué goroutines se lanzaron, qué errores ocurrieron, tiempos.
- Frontend: qué datos llegaron, errores de parsing.

**¿Cómo pensar antes?**

- ¿Qué información es útil para debuggear?
- ¿Cuánto logging es demasiado?

---

#### Tarea 8.5: Performance (Opcional)

**¿Qué hacer?**

Pequeñas optimizaciones:
- Backend: caching de consultas frecuentes.
- Frontend: lazy loading de componentes, memoization.

**¿Cómo pensar antes?**

- ¿Dónde está el cuello de botella?
- ¿La optimización vale el esfuerzo?

---

### Checklist de Finalización (Parte 8)

- [ ] No hay errores desconocidos.
- [ ] Mensajes de error son claros.
- [ ] Validación en frontend y backend.
- [ ] UI se ve bien en desktop y mobile.
- [ ] Logs ayudan a entender qué sucede.

---

---

## Métodos / Funciones / Componentes — Referencia Conceptual

### Backend: Función de Consulta a API Externa

**Qué hace**: Ejecuta una solicitud HTTP a un servicio externo, parsea la respuesta, maneja errores.

**Cuándo se llama**: Dentro de una goroutine, cuando necesitas datos frescos.

**Qué recibe**:
- Context (para timeouts y cancelación).
- Parámetros específicos (ID, coordenadas, etc).

**Qué devuelve**:
- Struct tipado con los datos.
- Error si algo falló.

**Responsabilidades**:
- Construir URL correctamente.
- Respetar el contexto (ctx.Done()).
- Manejar errores de red.
- Parsear JSON válido.

**Límites**:
- No toma decisiones sobre qué hacer si falla.
- No combina datos de múltiples fuentes.
- No aplica lógica de negocio.

```go
// Forma genérica
func queryExternalResource(ctx context.Context, param string) (ResultType, error) {
    // Construir solicitud
    // Ejecutar (respetar ctx)
    // Manejar errores
    // Parsear respuesta
    // Retornar resultado o error
}
```

---

### Backend: Función de Procesamiento Individual

**Qué hace**: Toma un dato crudo de la API, lo transforma o calcula algo.

**Cuándo se llama**: Inmediatamente después de obtener un dato.

**Qué recibe**:
- Dato crudo (del parser).

**Qué devuelve**:
- Dato procesado.
- Error si la transformación falla.

**Responsabilidades**:
- Validar que el dato tiene sentido.
- Transformar/enriquecer.

**Límites**:
- No hace más de una transformación.
- No combina con otros datos.
- No es concurrente.

```go
func processSingleItem(raw RawType) (ProcessedType, error) {
    // Validar
    // Transformar
    // Retornar
}
```

---

### Backend: Función de Orquestación Concurrente

**Qué hace**: Lanza múltiples goroutines, cada una consultando y procesando, y recolecta resultados.

**Cuándo se llama**: Cuando necesitas datos de múltiples fuentes en paralelo.

**Qué recibe**:
- Context (timeout global).
- Lista de identificadores.

**Qué devuelve**:
- Resultados exitosos.
- Errores parciales.
- Error total si algo grave sucede.

**Responsabilidades**:
- Coordinar goroutines.
- Sincronizar con WaitGroup o contador.
- Manejar timeouts.
- Separar éxitos de fallos.

**Límites**:
- No parsea respuestas (delega).
- No aplica agregaciones (solo recolecta).
- No formatea para HTTP.

```go
func orchestrateConcurrent(ctx context.Context, ids []string) ([]ProcessedType, []error, error) {
    // Crear channel
    // Lanzar goroutines (fan-out)
    // Recolectar resultados (fan-in)
    // Sincronizar
    // Retornar ambas listas
}
```

---

### Backend: Función de Agregación

**Qué hace**: Toma múltiples datos procesados, calcula estadísticas, rankings, extremos.

**Cuándo se llama**: Después de recolectar todos los resultados.

**Qué recibe**:
- Lista de datos procesados.

**Qué devuelve**:
- Struct con agregaciones (promedios, máximos, rankings, etc).
- Error si la lista está vacía o datos inválidos.

**Responsabilidades**:
- Iterar sobre datos.
- Realizar cálculos.
- Ordenar/agrupar según necesario.

**Límites**:
- Solo lee datos (no los modifica).
- No hace nuevas consultas.

```go
func aggregateData(data []ProcessedType) (AggregateType, error) {
    // Validar no vacío
    // Calcular promedios, máximos, etc
    // Crear ranking
    // Agrupar si necesario
    // Retornar agregaciones
}
```

---

### Backend: Handler HTTP

**Qué hace**: Recibe solicitud HTTP, valida, delega a lógica, formatea respuesta.

**Cuándo se llama**: Automáticamente cuando llega una solicitud HTTP.

**Qué recibe**:
- Solicitud HTTP (método, path, body, query).

**Qué devuelve**:
- Respuesta HTTP (status, headers, body).

**Responsabilidades**:
- Extraer parámetros.
- Validar entrada.
- Delegar a servicio.
- Formatear salida JSON.
- Escribir status code.

**Límites**:
- No contiene lógica de negocio.
- No consulta APIs externas directamente.
- No es puro (tiene side effects HTTP).

```go
func handlerForSomething(w http.ResponseWriter, r *http.Request) {
    // Extraer parámetros
    // Validar
    // Delegar
    // Formatear respuesta
    // Escribir status y body
}
```

---

### Frontend: Componente Selector

**Qué hace**: Renderiza una lista de opciones, permite seleccionar, emite eventos.

**Cuándo se renderiza**: En el montaje inicial y cuando props cambian.

**Qué recibe** (props):
- Lista de opciones.
- Callback opcional para cambios.

**Qué emite**: Eventos cuando el usuario selecciona/deselecciona.

**Responsabilidades**:
- Mostrar opciones.
- Capturar interacción.
- Emitir eventos hacia padre.

**Límites**:
- No valida opciones.
- No hace llamadas API.
- No almacena estado global.

```svelte
<script>
  export let options = [];
  // emit event on change
</script>

<div>
  {#each options as option}
    <input ... on:change={() => emitEvent()} />
  {/each}
</div>
```

---

### Frontend: Componente de Resultado

**Qué hace**: Renderiza un dato individual en forma clara.

**Cuándo se renderiza**: Cuando exista un dato para mostrar.

**Qué recibe** (props):
- Un objeto con datos procesados.

**Qué emite**: Nada (solo presentación).

**Responsabilidades**:
- Mostrar datos en formato legible.
- Indicar campos faltantes si aplica.

**Límites**:
- No modifica datos.
- No valida.
- No llama APIs.

```svelte
<script>
  export let data;
</script>

<div class="result">
  <h3>{data.name}</h3>
  <p>Attr: {data.value}</p>
</div>
```

---

### Frontend: Store (Estado Global)

**Qué hace**: Centraliza estado, permite suscripción reactiva.

**Cuándo se accede**: Desde cualquier componente que necesite compartir estado.

**Qué contiene**:
- Selecciones del usuario.
- Resultados de la API.
- Estado de carga.
- Errores.

**Cómo se actualiza**: Con `.set()` o `.update()`.

**Responsabilidades**:
- Mantener estado consistente.
- Notificar a componentes suscritos.

**Límites**:
- No contiene lógica de presentación.
- No hace validaciones complejas.
- No llama APIs directamente.

```javascript
export const appState = writable({
    selections: {},
    results: null,
    loading: false,
    error: null,
});
```

---

### Frontend: Función Cliente HTTP

**Qué hace**: Abstracta las complejidades de hacer requests HTTP.

**Cuándo se llama**: Desde funciones de "negocio" del frontend.

**Qué recibe**:
- Método HTTP.
- Endpoint.
- Datos (si aplica).

**Qué devuelve**:
- Respuesta parseada.
- Error si algo falló.

**Responsabilidades**:
- Construir request.
- Manejar response.
- Parsear JSON.
- Reportar errores claros.

**Límites**:
- No maneja estado.
- No toma decisiones de negocio.

```javascript
async function apiCall(method, endpoint, data) {
    // Construir request
    // Ejecutar
    // Chequear response
    // Parsear
    // Retornar o error
}
```

---

---

## Sintaxis y Ejemplos Genéricos

### Go: Definir Struct para Tipos de Datos

```go
// Struct con tags JSON
type DataItem struct {
    Name     string  `json:"name"`
    Value    float64 `json:"value"`
    Optional string  `json:"optional,omitempty"`
    Pointer  *int    `json:"pointer"`
}
```

**Notas**:
- Campo público = comienza con mayúscula.
- Tag `json:"campo_name"` mapea a JSON.
- `omitempty` evita incluir campos nulos.
- `*Type` es puntero (puede ser nil).

---

### Go: Función Básica que Retorna Error

```go
func DoSomething(input string) (OutputType, error) {
    if input == "" {
        return OutputType{}, errors.New("input cannot be empty")
    }
    
    // Do work
    result := OutputType{ /* ... */ }
    
    return result, nil
}

// Uso
result, err := DoSomething("test")
if err != nil {
    log.Fatal(err)
}
```

---

### Go: HTTP Request con Context

```go
func MakeRequest(ctx context.Context, url string) ([]byte, error) {
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }
    
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
    }
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    
    return body, nil
}
```

---

### Go: Parsing JSON a Struct

```go
var data DataItem
err := json.Unmarshal(jsonBytes, &data)
if err != nil {
    return fmt.Errorf("parse error: %w", err)
}
```

---

### Go: Goroutine con Channel

```go
type Result struct {
    Data  string
    Error error
}

func worker(id string, resultChan chan Result) {
    // Do work
    result := Result{
        Data: "some result",
    }
    resultChan <- result  // Send result
}

// Lanzar
resultChan := make(chan Result, 3)
go worker("id1", resultChan)
go worker("id2", resultChan)
go worker("id3", resultChan)

// Recolectar
for i := 0; i < 3; i++ {
    result := <-resultChan  // Receive result
    fmt.Println(result)
}
```

---

### Go: WaitGroup para Sincronización

```go
var wg sync.WaitGroup

for i := 0; i < 5; i++ {
    wg.Add(1)
    go func(i int) {
        defer wg.Done()
        // Do work
    }(i)
}

wg.Wait()  // Esperar a que todas terminen
```

---

### Go: Context con Timeout

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

// Usar ctx en requests, goroutines, etc
// Si pasa el tiempo, ctx.Done() se cierra
```

---

### Go: HTTP Handler

```go
func handleSomething(w http.ResponseWriter, r *http.Request) {
    // Validar método
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    // Parsear body
    var input struct {
        Name string `json:"name"`
    }
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    // Do business logic
    result := DoLogic(input.Name)
    
    // Write response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(result)
}
```

---

### Svelte: Componente Básico

```svelte
<script>
  export let title = "Default Title";
  export let onAction = () => {};
  
  let count = 0;
  
  function increment() {
    count++;
    onAction(count);
  }
</script>

<div class="card">
  <h2>{title}</h2>
  <p>Count: {count}</p>
  <button on:click={increment}>Increment</button>
</div>

<style>
  .card {
    border: 1px solid #ccc;
    padding: 1rem;
  }
</style>
```

---

### Svelte: Store Writable

```javascript
// store.js
import { writable } from 'svelte/store';

export const count = writable(0);

// En componente
<script>
  import { count } from './store.js';
  
  function increment() {
    count.update(n => n + 1);
  }
</script>

<p>Count: {$count}</p>
<button on:click={increment}>+</button>
```

---

### Svelte: Renderizado Condicional y Loops

```svelte
{#if loading}
  <p>Loading...</p>
{:else if error}
  <p class="error">{error}</p>
{:else}
  <div>
    {#each items as item (item.id)}
      <div>{item.name}</div>
    {/each}
  </div>
{/if}
```

---

### Svelte: Llamada a Función en Evento

```svelte
<script>
  async function handleSubmit() {
    try {
      const response = await fetch('/api/endpoint', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ data: "value" })
      });
      const result = await response.json();
      console.log(result);
    } catch (error) {
      console.error(error);
    }
  }
</script>

<button on:click={handleSubmit}>Submit</button>
```

---

### JavaScript: Función Async/Await Genérica

```javascript
async function fetchData(url) {
  try {
    const response = await fetch(url);
    if (!response.ok) throw new Error(`HTTP ${response.status}`);
    const data = await response.json();
    return data;
  } catch (error) {
    console.error('Error:', error);
    throw error;
  }
}

// Uso
const data = await fetchData('/api/data');
```

---

---

## Checklist Final — Antes de Considerar "Hecho"

- [ ] **Parte 1**: Backend y frontend compilan sin errores.
- [ ] **Parte 2**: Puedo consultar la API externa manualmente y parsear respuesta.
- [ ] **Parte 3**: Lógica secuencial funciona, medí tiempo base.
- [ ] **Parte 4**: Lógica concurrente funciona, tiempo es menor, resultados idénticos.
- [ ] **Parte 5**: Endpoints HTTP responden correctamente, status codes son semánticos.
- [ ] **Parte 6**: Frontend renderiza sin errores, componentes se pasan props.
- [ ] **Parte 7**: Frontend y backend se comunican, datos fluyen end-to-end.
- [ ] **Parte 8**: Validaciones robustas, mensajes claros, UI pulida.

---

---

## Conclusión

Esta receta no es un plan de ejecución lineal. Es una **guía conceptual** de qué pensar en cada etapa.

**Recuerda**:
1. **Cada parte debe funcionar independientemente** antes de avanzar.
2. **Prueba manualmente** (curl para backend, DevTools para frontend).
3. **Mide impacto** (velocidad secuencial vs. concurrente).
4. **Errores son señales** (no síntomas a ignorar).
5. **Refactoriza cuando sea necesario** (la arquitectura inicial no es perfecta).

**Las preguntas que debes hacerte en cada tarea**:
- ¿Qué responsabilidad tiene este código?
- ¿Quién lo invoca y por qué?
- ¿Qué no debería conocer?
- ¿Cómo sé que funciona sin depender de capas superiores?

Usa esta guía como referencia, pero **la solución es tuya**.
