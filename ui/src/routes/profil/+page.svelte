<script>
    import {currentUser} from "$lib/auth.js";
    import {t} from "$lib/i18n";
    import Tags from "svelte-tags-input";
    import {tags} from "$lib/tags.js";
    import {PUBLIC_TIL_SERVER_URL} from "$env/static/public";
    import toast from "svelte-french-toast";
    import {get} from "svelte/store";
    import {renewFeedKey, updateTags} from "$lib/user.js";

    let personalTags = $state([]);

    let feedUrlBase = `${PUBLIC_TIL_SERVER_URL.startsWith('http') ? `${PUBLIC_TIL_SERVER_URL}/rss?key=` : `${window.location.protocol}//${window.location.host}/api/rss?key=`}`;

    currentUser.subscribe(e => {
        personalTags = e.automatic_tags_filter
    })

    function copyRssToClipBoard() {
        navigator.clipboard.writeText(`${feedUrlBase}${$currentUser.feed_key}`)
        toast.success(get(t)('clipboard.copied'), {
            position: "top-right"
        });
    }
</script>

<div class="profile">
    <h1>
        {#if $currentUser.is_admin}👑{:else }👋{/if} {$currentUser.display_name}</h1>

    <hr>

    <h1>{$t('profile.automatic_tags.title')}</h1>
    <p>{$t('profile.automatic_tags.description')}</p>

    <Tags autoComplete={$tags} bind:tags={personalTags} onlyUnique={true}/>
    <button onclick={() => {updateTags(personalTags)}}><i class="fa-solid fa-floppy-disk"></i> {$t('profile.automatic_tags.save_button')}</button>

    <hr>
    <h1>{$t('profile.rss_feed.title')}</h1>
    <p>{$t('profile.rss_feed.description')}</p>
    <p><code>{feedUrlBase}{$currentUser.feed_key}</code> <i class="fa-solid fa-copy link" onclick={copyRssToClipBoard}></i></p>
    <p>{$t('profile.rss_feed.tags_manual')} <code>&tags=tag1,tag2,[...]</code>.</p>

    <div class="alert">
        <p><b>{$t('profile.rss_feed.warning.content')}</b><br/><br/>{$t('profile.rss_feed.warning.disconnection_warning')}</p>
        <p>
            <button onclick={renewFeedKey}><i class="fa-solid fa-biohazard"></i> {$t('profile.rss_feed.warning.renew_button')}</button>
        </p>
    </div>

</div>