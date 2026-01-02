<script>

	import { t } from '$lib/i18n.js';
	import { addBookmark, removeBookmark, bookmarks } from '$lib/bookmarks.js';
	import { get } from 'svelte/store';
	import { currentUser } from '$lib/auth.js';
	import { deletePost } from '$lib/posts.js';

	let isFavorite = $state(false);

	let {
		item = {
			id: '',
			title: '',
			content: '',
			link: '',
			tags: [],
			creationDate: '',
			user: {
				id: '',
				display_name: ''
			}
		}
	} = $props();


	bookmarks.subscribe(books => {
		if (item) {
			isFavorite = books.filter(e => e.id === item.id).length === 1;
		}
	});

	$effect(() => {
		isFavorite = get(bookmarks).filter(e => e.id === item.id).length === 1;
	});


</script>

<div class="card">
	<div>
		<div class="card-title">
			<div>
				<h3>
					{#if item.icon}
						<img src={item.icon} width="16" />
					{/if}
					{#if item.tags && item.tags.filter(e => e.startsWith("lang:")).length !== 0}
						{item.tags.filter(e => e.startsWith("lang:").length)}
						{$t(`lang.${item.tags.filter(e => e.startsWith("lang:"))[0].replace("lang:", "")}`)}
					{/if}
					{#if (item.link)}
						<a href={item.link} target="_blank" rel="noopener noreferrer">
							{item.title}
						</a>
					{:else}
						{item.title}
					{/if}

				</h3>
				<sub>{$t('item.by', { user: item.user.display_name, date: new Date(item.creation_date).toLocaleDateString() })}</sub>
			</div>
			<div>
				{#if isFavorite}
					<span class="fas fa-bookmark" on:click={() => {removeBookmark(item.id)}}></span>
				{:else }
					<span class="far fa-bookmark" on:click={() => {addBookmark(item.id)}}></span>
				{/if}
			</div>
		</div>
	</div>
	<p>{item.content}</p>
	{#if item.image}
		<div class="article-image" style={`background-image: url(${item.image})`} />
	{/if}
	<div class="card-footer">
		<div>
			<p>
				{#if item.tags}
					<span class="fas fa-tags"></span>
					{#each item.tags.filter(e => !e.startsWith("lang:")) as tag}<span class="tag">#{tag}</span> {/each}
				{/if}
			</p>
		</div>
		{#if $currentUser.is_admin || item.user.id === $currentUser.id }
			<div class="delete-icon">
				<span class="far fa-trash-can" on:click={() => {deletePost(item.id)}}></span>
			</div>
		{/if}
	</div>
</div>