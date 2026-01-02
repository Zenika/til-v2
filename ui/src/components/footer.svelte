<script>
    import {locales, locale, t} from "$lib/i18n.js";
    import {changeTheme} from "$lib/theme.js";
</script>

<div class="footer-option">
    <a onclick={changeTheme}><span class="fas fa-sun nopad"></span></a>

    <select bind:value={$locale}>
        {#each locales as locale}
            <option value={locale}>{$t(`lang.${locale}`)}</option>
        {/each}
    </select>

</div>

