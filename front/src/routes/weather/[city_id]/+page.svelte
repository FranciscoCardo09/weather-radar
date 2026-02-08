<script lang="ts">
    import { page } from '$app/stores';
    import { getWeatherForCity } from '$lib/api';
    import type { WeatherData } from '$lib/types';

    let loading = true;
    let error: string | null = null;
    let weather: WeatherData | null = null;
    let lastCityId: string | null = null;

    $: cityId = $page.params.city_id;

    async function loadWeather(id: string) {
        loading = true;
        error = null;
        try {
            weather = await getWeatherForCity(id);
            lastCityId = id;
        } catch (e) {
            error = e instanceof Error ? e.message : 'Error desconocido';
            weather = null;
        } finally {
            loading = false;
        }
    }

    $: if (cityId && cityId !== lastCityId) {
        loadWeather(cityId);
    }
</script>

<svelte:head>
    <link rel="stylesheet" href="/weather.css">
</svelte:head>

<main class="page">
    <a class="back" href="/">Volver</a>

    {#if loading}
        <p class="state">Cargando clima...</p>
    {:else if error}
        <p class="state error">❌ {error}</p>
    {:else if weather}
        <section class="card">
            <h1>{weather.city_name}</h1>
            <ul>
                <li><strong>Temperatura:</strong> {weather.temperature.toFixed(1)}°C</li>
                <li><strong>Humedad:</strong> {weather.humidity.toFixed(0)}%</li>
                <li><strong>Viento:</strong> {weather.wind_speed.toFixed(1)} km/h</li>
                <li><strong>Condición:</strong> {weather.condition}</li>
            </ul>
        </section>
    {:else}
        <p class="state">No hay datos disponibles.</p>
    {/if}
</main>
