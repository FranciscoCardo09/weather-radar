package main

// FIX: Usar map para búsquedas O(1) en lugar de O(n)
var citiesMap = map[string]Cities{
	"cordoba":               {ID: "cordoba", Name: "Córdoba", Latitude: -31.42, Longitude: -64.18},
	"buenosaires":           {ID: "buenosaires", Name: "Buenos Aires", Latitude: -34.61, Longitude: -58.38},
	"rosario":               {ID: "rosario", Name: "Rosario", Latitude: -32.95, Longitude: -60.65},
	"mendoza":               {ID: "mendoza", Name: "Mendoza", Latitude: -32.89, Longitude: -68.85},
	"la_plata":              {ID: "la_plata", Name: "La Plata", Latitude: -34.92, Longitude: -57.95},
	"mar_del_plata":         {ID: "mar_del_plata", Name: "Mar del Plata", Latitude: -38.00, Longitude: -57.55},
	"san_miguel_de_tucuman": {ID: "san_miguel_de_tucuman", Name: "San Miguel de Tucumán", Latitude: -26.82, Longitude: -65.22},
	"salta":                 {ID: "salta", Name: "Salta", Latitude: -24.78, Longitude: -65.41},
	"santa_fe":              {ID: "santa_fe", Name: "Santa Fe", Latitude: -31.63, Longitude: -60.70},
	"corrientes":            {ID: "corrientes", Name: "Corrientes", Latitude: -27.48, Longitude: -58.83},
	"neuquen":               {ID: "neuquen", Name: "Neuquén", Latitude: -38.95, Longitude: -68.06},
	"resistencia":           {ID: "resistencia", Name: "Resistencia", Latitude: -27.45, Longitude: -58.99},
	"posadas":               {ID: "posadas", Name: "Posadas", Latitude: -27.37, Longitude: -55.90},
	"bariloche":             {ID: "bariloche", Name: "Bariloche", Latitude: -41.13, Longitude: -71.31},
	"ushuaia":               {ID: "ushuaia", Name: "Ushuaia", Latitude: -54.80, Longitude: -68.30},
}

var citiesList = []Cities{
	citiesMap["cordoba"],
	citiesMap["buenosaires"],
	citiesMap["rosario"],
	citiesMap["mendoza"],
	citiesMap["la_plata"],
	citiesMap["mar_del_plata"],
	citiesMap["san_miguel_de_tucuman"],
	citiesMap["salta"],
	citiesMap["santa_fe"],
	citiesMap["corrientes"],
	citiesMap["neuquen"],
	citiesMap["resistencia"],
	citiesMap["posadas"],
	citiesMap["bariloche"],
	citiesMap["ushuaia"],
}

// GetCities retorna la lista completa de ciudades argentinas soportadas.
// Esta lista es estática y contiene 15 ciudades principales.
func GetCities() []Cities {
	return citiesList
}

// GetCityByID busca una ciudad por su ID.
// Retorna nil si la ciudad no existe.
// FIX: Usa map lookup O(1) en lugar de búsqueda lineal O(n)
func GetCityByID(id string) *Cities {
	if city, ok := citiesMap[id]; ok {
		return &city
	}
	return nil
}
