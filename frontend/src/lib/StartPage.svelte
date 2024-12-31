<script>
    let state = $state("init")

    async function getGameCode() {
        const url = "/api/new/game";
        try {
            const response = await fetch(url);
            if (!response.ok) {
                throw new Error(`Response status: ${response.status}`);
            }

            const json = await response.json();
            return json;
        } catch (error) {
            console.error(error.message);
            return null
        }
    }


</script>

<div>
    {#if state == "init"}
        <button on:click={()=>{state = "join"}}>join</button>
        <button on:click={()=>{state = "new"}}>new game</button>
    {/if}

    {#if state == "new"}
        {#await getGameCode()}
            <p>Getting game code...</p>
        {:then ret}
            {ret}
        {/await}
    {/if}

    {#if state == "join"}
        <label for="gcode">Game Code</label>
        <input type="text" id="gcode" name="gcode">
    {/if}
</div>

<style>
    div {
        display: flex;
        align-items: center;
        justify-content: center;
    }
    button {
        margin: 0.2rem;
    }
</style>