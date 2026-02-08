<script lang="ts">
    import { onMount } from 'svelte';
    import { getCities } from '$lib/api';
    import type { Cities } from '$lib/types';

    let error = '';
    let cities: Cities[] = [];
    let loading = true;
    let selectedIds: string[] = [];

    onMount(async () => {
        try {
            cities = await getCities();
        } catch (e) {
            error = e instanceof Error ? e.message : 'Error desconocido';
            console.error('Error cargando ciudades:', e);
        } finally {
            loading = false;
        }
    });

    function handleCompare() {
        if (selectedIds.length < 2) {
            error = 'Selecciona al menos dos ciudades para comparar.';
            return;
        }
        const url = `/compare?cityIds=${selectedIds.join(',')}`;
        window.location.href = url;
    }
</script>
<svelte:head>
    <link rel="stylesheet" href="/cities.css">
</svelte:head>
<main>
    <h1>Lista de Ciudades</h1>
    
    {#if loading}
        <div class="state">Cargando ciudades...</div>
    {:else if error}
        <div class="state error">❌ {error}</div>
        <div class="state hint">Verifica que el servidor Go esté corriendo en puerto 8080.</div>
    {:else if cities.length === 0}
        <div class="state">No hay ciudades disponibles.</div>
    {:else}
        <div class="cities-grid">
            {#each cities as city}
                <div class="city-item">
                    <a href="/weather/{city.id}" class="city-card" role="button">
                        <h3>{city.name}</h3>
                        <p>📍 {city.latitude.toFixed(2)}°, {city.longitude.toFixed(2)}°</p>
                    </a>
                    <label class="select-row">
                        <input type="checkbox" bind:group={selectedIds} value={city.id} />
                        <span>Comparar</span>
                    </label>
                </div>
            {/each}
        </div>
        
        {#if cities.length > 1}
            <button class="compare-btn" on:click={handleCompare} disabled={selectedIds.length < 2}>
                Comparar ciudades
            </button>
        {/if}
    {/if}
</main>
