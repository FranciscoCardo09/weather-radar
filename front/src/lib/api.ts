const API_URL = 'http://localhost:8080/api';

export async function getCities() {
    const response = await fetch(`${API_URL}/cities`);
    if (!response.ok) {
        throw new Error('Error al obtener las ciudades');
    }
    return await response.json();
}

export async function compareWeather(cityId1: string, cityId2: string) {
  const response = await fetch(`http://localhost:8080/api/compare/${cityId1}/${cityId2}`);
  return response.json();
}

export async function getWeatherForCity(cityId: string) {
    const response = await fetch(`http://localhost:8080/api/weather/${cityId}`);
    if (!response.ok) {
        throw new Error('Error al obtener el clima de la ciudad');
    }
    return await response.json();
}