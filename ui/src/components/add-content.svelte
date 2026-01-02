<script>
    import {locales, t} from "$lib/i18n.js";
    import Tags from "svelte-tags-input";
    import {tags} from "$lib/tags.js";
    import {createPost} from "$lib/posts.js";

    let form = $state({
       title: "",
       link: "",
       lang: "fr",
       tags: [],
       content: ""
    });

    function submit() {
        createPost(form.title, form.link, form.lang, form.tags, form.content)
            .then(() => {
                form = {
                    title: "",
                    link: "",
                    lang: "fr",
                    tags: [],
                    content: ""
                }
            })
    }
</script>

<h1>{$t('post.page_title')}</h1>
<form onsubmit={submit}>
    <div class="form-flex">
        <div>
            <div class="form-group">
                <label for="title">{$t('post.title')}</label><input type="text" name="title" id="title" placeholder={$t('post.title_placeholder')} bind:value={form.title} required/>
            </div>

            <div class="form-group">
                <label for="link">{$t('post.link')}</label><input type="text" name="link" id="link" placeholder={$t('post.link_placeholder')}  bind:value={form.link} required/>
            </div>

            <div class="form-group">
                <label for="lang">{$t('post.lang')}</label><select  bind:value={form.lang}>
                {#each locales as locale}
                    <option value={locale}>{$t(`lang.${locale}`)}</option>
                {/each}
            </select>
            </div>

            <div class="form-group">
                <label for="tags">{$t('post.tags')}</label><Tags autoComplete={$tags} id="tags" name="tags" onlyUnique={true} bind:tags={form.tags}/>
            </div>
        </div>
        <div>
            <label for="content" style="text-align: left">{$t('post.comment')}</label>
            <br/>
            <textarea name="content" id="content" cols="43" rows="3" maxlength="180" bind:value={form.content}></textarea>
            <br/>
            <br/>
            <button type="submit">{$t('post.create_button')}</button>
        </div>

    </div>
</form>