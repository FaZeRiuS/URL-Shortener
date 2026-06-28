<script lang="ts">
    import type { ErrorResponse, StatsResponse } from "../types";

    let statsCount = $state(0);
    let statsOriginalUrl = $state("");
    let statsCode = $state("");

    let ErrorMessage = $state("");

    async function getStats() {
        statsCount = 0;
        ErrorMessage = "";
        statsOriginalUrl = "";

        if (statsCode.length != 6) {
            ErrorMessage = "Code must be 6 characters";
            return;
        }

        fetch(`/stats/${statsCode}`)
            .then((res) => {
                if (!res.ok) {
                    return res.json().then((data) => {
                        const errorData = data as ErrorResponse;
                        throw new Error(errorData.error);
                    });
                }
                return res.json();
            })
            .then((data) => {
                const statsResponse = data as StatsResponse;
                statsCount = statsResponse.count;
                statsOriginalUrl = statsResponse.original_url;
            })
            .catch((err) => {
                ErrorMessage = err.message;
            });
    }
</script>

<h2>Порахувати кількість переходів</h2>
<div id="stats">
    <br />
    <input
        type="text"
        id="statsCode"
        placeholder="Вставте сюди код..."
        bind:value={statsCode}
    />
    <button onclick={getStats}>Отримати статистику</button>
    {#if statsOriginalUrl}
        <p id="stats">Кількість переходів: {statsCount}</p>
        <br />
        <p id="originalUrl">
            Оригінальне посилання: <a href={statsOriginalUrl} target="_blank"
                >{statsOriginalUrl}</a
            >
        </p>
    {/if}
    {#if ErrorMessage}
        <p id="error">Помилка: {ErrorMessage}</p>
    {/if}
</div>
