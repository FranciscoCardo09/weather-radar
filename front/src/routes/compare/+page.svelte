<script lang="ts">
    import { onMount } from 'svelte';
    import { page } from '$app/stores';
    import { compareWeather, getCities } from '$lib/api';
    import type { Cities, WeatherData, CompareResult } from '$lib/types';

    let cities: Cities[] = [];
    let selectedIds: string[] = [];
    let loading = true;
    let error: string | null = null;
    let weatherCards: WeatherData[] = [];
    let summary: CompareResult['summary'] | null = null;

    onMount(async () => {
        try {
            cities = await getCities();
            const idsParam = $page.url.searchParams.get('cityIds');
            if (idsParam) {
                selectedIds = idsParam.split(',').filter(Boolean);
                if (selectedIds.length >= 2) {
                    handleCompare();
                }
            }
        } catch (e) {
            error = e instanceof Error ? e.message : 'Error desconocido';
        } finally {
            loading = false;
        }
    });

    async function handleCompare() {
        if (selectedIds.length < 2) {
            error = 'Selecciona al menos dos ciudades.';
            return;
        }
        
        error = null;
        summary = null;
        weatherCards = [];
        
        try {
            const response = await compareWeather(selectedIds);
            weatherCards = response.cities || [];
            summary = response.summary || null;
            
        } catch (e) {
            error = e instanceof Error ? e.message : 'Error desconocido';
        }
    }

    // FIX: Usar directamente el ranking del backend en lugar de recalcularlo
    function rankingText() {
        if (!summary || !summary.ranking || summary.ranking.length === 0) return '';
        return summary.ranking.map((city, index) => {
            const weather = weatherCards.find(w => w.city_name === city);
            const temp = weather ? weather.temperature.toFixed(1) : 'N/A';
            return `${index + 1}. ${city} (${temp}°C)`;
        }).join(' - ');
    }

    // FIX: Usar promedios del backend en lugar de recalcularlos
    function averageText() {
        if (!summary) return '';
        return `Temp promedio: ${summary.average_temperature.toFixed(1)}°C, Humedad: ${summary.average_humidity.toFixed(0)}%, Viento: ${summary.average_wind_speed.toFixed(0)} km/h`;
    }

    // FIX: Usar extremos del backend en lugar de recalcularlos
    function extremesText() {
        if (!summary) return '';
        return `Más caliente: ${summary.hotter_city}, Más fría: ${summary.colder_city}, Más ventosa: ${summary.windy_city}`;
    }

    // FIX: Usar agrupación por condición del backend
    function byConditionText() {
        if (!summary || !summary.by_condition) return '';
        return Object.entries(summary.by_condition)
            .map(([k, v]) => `${k}: ${v.length} ciudad${v.length === 1 ? '' : 'es'}`)
            .join(', ');
    }
</script>

<svelte:head>
    <link rel="stylesheet" href="/compare.css">
</svelte:head>

<main class="page">
    <a class="back" href="/">Volver</a>
    <h1>Comparar clima</h1>

    {#if loading}
        <p class="state">Cargando ciudades...</p>
    {:else}
        <div class="form">
            {#each cities as city}
                <label class="row">
                    <input type="checkbox" bind:group={selectedIds} value={city.id} />
                    <span>{city.name}</span>
                </label>
            {/each}

            <button class="btn" on:click={handleCompare}>Comparar</button>
        </div>
    {/if}

    {#if error}
        <p class="state error">❌ {error}</p>
    {/if}

    {#if weatherCards.length > 0}
        <div class="cards">
            {#each weatherCards as w}
                <div class="city-card">
                    <h2>{w.city_name}</h2>
                    <p>🌡️ {w.temperature.toFixed(1)}°C</p>
                    <p>💧 {w.humidity.toFixed(0)}% humedad</p>
                    <p>💨 {w.wind_speed.toFixed(0)} km/h</p>
                    <p>☀️ {w.condition}</p>
                </div>
            {/each}
        </div>

        <div class="summary">
            <p><strong>Ranking:</strong> {rankingText()}</p>
            <p><strong>Promedios:</strong> {averageText()}</p>
            <p><strong>Extremos:</strong> {extremesText()}</p>
            <p><strong>Por condición:</strong> {byConditionText()}</p>
        </div>
    {/if}
</main>
