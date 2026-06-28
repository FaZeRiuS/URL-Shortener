<script lang="ts">
    import type { ShortenResponse, ErrorResponse } from "../types";

    let longUrl = $state("");
    let shortUrl = $state("");

    let ErrorMessage = $state("");

    async function shorten() {
        shortUrl = "";
        ErrorMessage = "";

        if (longUrl.trim() == "") {
            ErrorMessage = "URL cannot be empty";
            return;
        } else if (
            !longUrl.startsWith("http://") &&
            !longUrl.startsWith("https://")
        ) {
            ErrorMessage = "Invalid URL";
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
                ErrorMessage = err.message;
            });
    }
</script>

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
            Коротке посилання: <a href={shortUrl} target="_blank">{shortUrl}</a>
        </p>
    {/if}
    {#if ErrorMessage}
        <p id="error">Помилка: {ErrorMessage}</p>
    {/if}
</div>
