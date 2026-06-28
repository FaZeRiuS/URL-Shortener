<script lang="ts">
    import type {
        ShortenRequest,
        ShortenResponse,
        StatsResponse,
        ErrorResponse,
    } from "./types";

    let longUrl = $state("");
    let shortUrl = $state("");

    let statsCode = $state("");
    let statsCount = $state(0);
    let statsOriginalUrl = $state("");

    let shortenErrorMessage = $state("");
    let statsErrorMessage = $state("");

    function shorten() {
        shortUrl = "";
        shortenErrorMessage = "";

        if (longUrl.trim() == "") {
            shortenErrorMessage = "URL cannot be empty";
            return;
        } else if (
            !longUrl.startsWith("http://") &&
            !longUrl.startsWith("https://")
        ) {
            shortenErrorMessage = "Invalid URL";
            return;
        }

        fetch("/shorten", {
            headers: {
                "Content-Type": "application/json",
            },
            method: "POST",
            body: JSON.stringify({ url: longUrl }),
        })
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
                const shortenResponse = data as ShortenResponse;
                shortUrl = shortenResponse.short_url;
            })
            .catch((err) => {
                shortenErrorMessage = err.message;
            });
    }

    function getStats() {
        statsCount = 0;
        statsErrorMessage = "";
        statsOriginalUrl = "";

        if (statsCode.length != 6) {
            statsErrorMessage = "Code must be 6 characters";
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
                statsErrorMessage = err.message;
            });
    }
</script>

<main>
    <h1>Мій URL Shortener на Go+Svelte</h1>
    <div id="shorten">
        <input
            type="text"
            id="longUrl"
            placeholder="Вставте сюди довге посилання..."
            style="width: 300px"
            bind:value={longUrl}
        />

        <button onclick={shorten}>Скоротити</button>
        {#if shortUrl}
            <p id="result">
                Коротке посилання: <a href={shortUrl} target="_blank"
                    >{shortUrl}</a
                >
            </p>
        {/if}
        {#if shortenErrorMessage}
            <p id="error">Помилка: {shortenErrorMessage}</p>
        {/if}
    </div>
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
        {#if statsCount}
            <p id="stats">Кількість переходів: {statsCount}</p>
            <br />
            <p id="originalUrl">
                Оригінальне посилання: <a
                    href={statsOriginalUrl}
                    target="_blank">{statsOriginalUrl}</a
                >
            </p>
        {/if}
        {#if statsErrorMessage}
            <p id="error">Помилка: {statsErrorMessage}</p>
        {/if}
    </div>
</main>
