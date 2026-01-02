<script>
	import { onDestroy, onMount } from 'svelte';
	import Content from '../components/content.svelte';
	import { posts, getNextPage, selectedTags } from '$lib/posts.js';
	import { currentUser } from '$lib/auth.js';
	import Tags from 'svelte-tags-input';
	import { t } from '$lib/i18n';
	import { tags } from '$lib/tags.js';
	import { get } from 'svelte/store';

	let cardList = $state();
	let tagsInSelector = $state([]); // Component is quite old so we have to add a compatibility layer. Not perfect, but it works.

	selectedTags.subscribe(newTags => {
		tagsInSelector = newTags;
	})

	// Stub functions to interconnect with our store efficiently.
	function handleTagChanges() {
		selectedTags.set(tagsInSelector)
	}

	function clearTags() {
		selectedTags.set([]);
	}

	function loadUserTags() {
		selectedTags.set(get(currentUser).automatic_tags_filter);
	}

	// Handle wheel for infiniscroll
	function wheel() {
		if(cardList) {
			if (window.innerHeight + window.scrollY + 200 >= cardList.clientHeight) {
				getNextPage();
			}
		}
	}

	onMount(() => {
		window.addEventListener('wheel', wheel, { passive: true });
		window.addEventListener('scroll', wheel, { passive: true });
	});

	onDestroy(() => {
		window.removeEventListener('wheel', wheel);
		window.removeEventListener('scroll', wheel);
	});
</script>

<div class="filter-cards-flex">
	<div class="filter-cards">
		<label for="filter">{$t('index.filter.label')}</label>
		<Tags autoComplete={$tags} onlyUnique={true} bind:tags={tagsInSelector} name="filter" id="filter" removeKeys={[]} onTagAdded={handleTagChanges}
					onTagRemoved={handleTagChanges} onlyAutocomplete={true} />
	</div>
	<div>
		<button onclick={clearTags}><span class="far fa-trash-can"></span> <span class="menu-text">{$t('index.filter.remove_all')}</span></button>
		<button onclick={loadUserTags}><span class="far fa-user"></span> <span class="menu-text">{$t('index.filter.restore')}</span></button>
	</div>
</div>
{#if $posts.length === 0}
	<div class="nothing-to-show">
		<span class="fas fa-sad-cry"></span><br />
		<p>{$t('post.empty_list')}</p>
	</div>
{:else }
	<div class="cards" bind:this={cardList}>
		{#each $posts as article}
			<Content item={article}></Content>
		{/each}
	</div>
{/if}