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

    function rankingText() {
        const list = [...weatherCards].sort((a, b) => b.temperature - a.temperature);
        return list.map((c, index) => `${index + 1}. ${c.city_name} (${c.temperature.toFixed(1)}°C)`).join(' - ');
    }

    function averageText() {
        if (!weatherCards.length) return '';
        const avgTemp = weatherCards.reduce((a, b) => a + b.temperature, 0) / weatherCards.length;
        const avgHum = weatherCards.reduce((a, b) => a + b.humidity, 0) / weatherCards.length;
        const avgWind = weatherCards.reduce((a, b) => a + b.wind_speed, 0) / weatherCards.length;
        return `Temp promedio: ${avgTemp.toFixed(1)}°C, Humedad: ${avgHum.toFixed(0)}%, Viento: ${avgWind.toFixed(0)} km/h`;
    }

    function extremesText() {
        if (!weatherCards.length) return '';
        const hotter = weatherCards.reduce((a, b) => (a.temperature > b.temperature ? a : b));
        const colder = weatherCards.reduce((a, b) => (a.temperature < b.temperature ? a : b));
        const windy = weatherCards.reduce((a, b) => (a.wind_speed > b.wind_speed ? a : b));
        return `Más caliente: ${hotter.city_name}, Más fría: ${colder.city_name}, Más ventosa: ${windy.city_name}`;
    }

    function byConditionText() {
        if (!weatherCards.length) return '';
        const counts: Record<string, number> = {};
        for (const c of weatherCards) {
            counts[c.condition] = (counts[c.condition] || 0) + 1;
        }
        return Object.entries(counts)
            .map(([k, v]) => `${k}: ${v} ciudad${v === 1 ? '' : 'es'}`)
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
