import type { CompareResult } from './types';

// FIX: Usar variable de entorno para la URL del API
// Permite configurar diferentes URLs para desarrollo, staging y producción
const API_URL = import.meta.env.PUBLIC_API_URL || 'http://localhost:8080/api';

export async function getCities() {
    const response = await fetch(`${API_URL}/cities`);
    if (!response.ok) {
        throw new Error('Error al obtener las ciudades');
    }
    return await response.json();
}

export async function compareWeather(cityIds: string[]): Promise<CompareResult> {
    const response = await fetch(`${API_URL}/compare`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ city_ids: cityIds })
    });
    
    if (!response.ok) {
        throw new Error('Error al comparar ciudades');
    }
    
    return await response.json();
}

export async function getWeatherForCity(cityId: string) {
    // FIX: Usar encodeURIComponent para manejar caracteres especiales en URLs
    const url = `${API_URL}/weather/${encodeURIComponent(cityId)}`;
    const response = await fetch(url);
    if (!response.ok) {
        throw new Error('Error al obtener el clima de la ciudad');
    }
    return await response.json();
}