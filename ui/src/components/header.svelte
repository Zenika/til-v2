<script>
    import {locales, locale, t} from "$lib/i18n.js";
    import {changeTheme} from "$lib/theme.js";
    import {logout} from "$lib/auth.js";
</script>

<nav class="main">
    <div>
        <ul>
            <li><img id="logo" src="/til.svg" width="50" alt={$t('a11y.logo')}/></li>
            <li><a href="/"><span class="fas fa-home"></span> <span class="menu-text">{$t('menu.home')}</span></a></li>
            <li><a href="/publish"><span class="fas fa-pencil"></span> <span class="menu-text">{$t('menu.publish')}</span></a></li>
            <li><a href="/bookmarks"><span class="fas fa-bookmark"></span> <span class="menu-text">{$t('menu.bookmarks')}</span></a></li>
            <li><a href="/extension"><span class="fas fa-window-restore"></span> <span class="menu-text">{$t('menu.extension')}</span></a></li>
            <li><a href="/profil"> <span class="fas fa-user"></span> <span class="menu-text">{$t('menu.profile')}</span></a></li>
        </ul>
    </div>
    <div>
        <ul>
            <li><a onclick={logout}><span class="fas fa-power-off nopad"></span></a></li>
            <li class="header-option"><a onclick={changeTheme}><span class="fas fa-sun nopad"></span></a></li>
            <li class="header-option">
                <select bind:value={$locale}>
                    {#each locales as locale}
                        <option value={locale}>{$t(`lang.${locale}`)}</option>
                    {/each}
                </select>
            </li>
        </ul>
    </div>
</nav>