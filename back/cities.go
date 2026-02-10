package main

func GetCities() []Cities {
	return []Cities{
		{ID: "cordoba", Name: "Córdoba", Latitude: -31.42, Longitude: -64.18},
		{ID: "buenosaires", Name: "Buenos Aires", Latitude: -34.61, Longitude: -58.38},
		{ID: "rosario", Name: "Rosario", Latitude: -32.95, Longitude: -60.65},
		{ID: "mendoza", Name: "Mendoza", Latitude: -32.89, Longitude: -68.85},
		{ID: "la_plata", Name: "La Plata", Latitude: -34.92, Longitude: -57.95},
		{ID: "mar_del_plata", Name: "Mar del Plata", Latitude: -38.00, Longitude: -57.55},
		{ID: "san_miguel_de_tucuman", Name: "San Miguel de Tucumán", Latitude: -26.82, Longitude: -65.22},
		{ID: "salta", Name: "Salta", Latitude: -24.78, Longitude: -65.41},
		{ID: "santa_fe", Name: "Santa Fe", Latitude: -31.63, Longitude: -60.70},
		{ID: "corrientes", Name: "Corrientes", Latitude: -27.48, Longitude: -58.83},
		{ID: "neuquen", Name: "Neuquén", Latitude: -38.95, Longitude: -68.06},
		{ID: "resistencia", Name: "Resistencia", Latitude: -27.45, Longitude: -58.99},
		{ID: "posadas", Name: "Posadas", Latitude: -27.37, Longitude: -55.90},
		{ID: "bariloche", Name: "Bariloche", Latitude: -41.13, Longitude: -71.31},
		{ID: "ushuaia", Name: "Ushuaia", Latitude: -54.80, Longitude: -68.30},
	}
}

func GetCityByID(id string) *Cities {
	cities := GetCities()
	for _, city := range cities {
		if city.ID == id {
			return &city
		}
	}
	return nil
}
