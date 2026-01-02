<script>
    import Content from "../../components/content.svelte";
    import {bookmarks} from "$lib/bookmarks.js";
    import {t} from "$lib/i18n";
</script>

{#if $bookmarks.length === 0}
    <div class="nothing-to-show">
        <span class="fas fa-sad-cry"></span><br/>
        <p>{$t('bookmarks.empty_list')}</p>
    </div>
{:else }
    <div class="cards">
        {#each $bookmarks as article}
            <Content item={article}></Content>
        {/each}
    </div>
{/if}

