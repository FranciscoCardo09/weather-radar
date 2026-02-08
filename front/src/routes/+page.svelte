<script lang="ts">
    import { onMount } from 'svelte';
    import { getCities, compareWeather, getWeatherForCity } from '$lib/api';
    import type { Cities, WeatherData } from '$lib/types';
    
    let cities: Cities[] = [];
    let cityId1: string = '';
    let cityId2: string = '';
    let comparisonResult: any = null;
    let weatherCity1: WeatherData | null = null;
    let weatherCity2: WeatherData | null = null;
    let errorMessage: string = '';
    
    onMount(async () => {
        cities = await getCities();
    });
    
    async function compare() {
        errorMessage = '';
        comparisonResult = null;
        weatherCity1 = null;
        weatherCity2 = null;
    
        try {
        comparisonResult = await compareWeather(cityId1, cityId2);
        weatherCity1 = await getWeatherForCity(cityId1);
        weatherCity2 = await getWeatherForCity(cityId2);
        } catch (error) {
        errorMessage = 'Error al obtener los datos del clima.';
        }
    }
</script>
<main>
    
</main>