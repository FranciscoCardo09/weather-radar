<script lang="ts">
    import { onMount } from 'svelte';
    import { page } from '$app/stores';
    import { compareWeather, getCities } from '$lib/api';
    import type { Cities } from '$lib/types';

    let cities: Cities[] = [];
    let cityA = '';
    let cityB = '';
    let loading = true;
    let error: string | null = null;
    let result: boolean | null = null;

    onMount(async () => {
        try {
            cities = await getCities();
            const cityAParam = $page.url.searchParams.get('cityA');
            const cityBParam = $page.url.searchParams.get('cityB');
            if (cityAParam) cityA = cityAParam;
            if (cityBParam) cityB = cityBParam;
            if (cityAParam && cityBParam) {
                handleCompare();
            }
        } catch (e) {
            error = e instanceof Error ? e.message : 'Error desconocido';
        } finally {
            loading = false;
        }
    });

    async function handleCompare() {
        if (!cityA || !cityB || cityA === cityB) {
            error = 'Selecciona dos ciudades distintas.';
            return;
        }
        error = null;
        result = null;
        try {
            const response = await compareWeather(cityA, cityB);
            result = Boolean(response?.comparison);
        } catch (e) {
            error = e instanceof Error ? e.message : 'Error desconocido';
        }
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
            <label>
                Ciudad A
                <select bind:value={cityA}>
                    <option value="" disabled>Selecciona una ciudad</option>
                    {#each cities as city}
                        <option value={city.id}>{city.name}</option>
                    {/each}
                </select>
            </label>

            <label>
                Ciudad B
                <select bind:value={cityB}>
                    <option value="" disabled>Selecciona una ciudad</option>
                    {#each cities as city}
                        <option value={city.id}>{city.name}</option>
                    {/each}
                </select>
            </label>

            <button class="btn" on:click={handleCompare}>Comparar</button>
        </div>
    {/if}

    {#if error}
        <p class="state error">❌ {error}</p>
    {/if}

    {#if result !== null}
        <div class="result {result ? 'ok' : 'diff'}">
            {#if result}
                El clima es igual en ambas ciudades.
            {:else}
                El clima es diferente entre las ciudades.
            {/if}
        </div>
    {/if}
</main>
