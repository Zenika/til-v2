import axios from 'axios';
import { PUBLIC_TIL_SERVER_URL } from '$env/static/public';
import toast from 'svelte-french-toast';
import { goto } from '$app/navigation';
import { get, writable } from 'svelte/store';
import { currentUser, token } from '$lib/auth.js';
import { t } from '$lib/i18n.js';
import { getTags } from '$lib/tags.js';
import { EventSource } from 'eventsource';

export const posts = writable([]);
export const selectedTags = writable([]);

let loading = false;
let totalPages = 1;
let retrievedPages = 0;
let eventSrc;

currentUser.subscribe((newUser) => {
	if (newUser.id !== '') {
		posts.set([]);
		if(newUser.automatic_tags_filter) {
			selectedTags.set(newUser.automatic_tags_filter);
		} else {
			totalPages = 1;
			retrievedPages = 0;
			getNextPage()
		}

	}
});

selectedTags.subscribe(() => {
	posts.set([]);
	totalPages = 1;
	retrievedPages = 0;
	getNextPage()
});

token.subscribe((newToken) => {
	if (newToken !== '') {
		if (eventSrc) {
			eventSrc.close();
		}

		eventSrc = new EventSource(`${PUBLIC_TIL_SERVER_URL}/posts/stream`, {
			fetch: (input, init) =>
				fetch(input, {
					...init,
					headers: {
						...init.headers,
						Authorization: newToken
					}
				})
		});
		eventSrc.addEventListener('updated', (article) => {
			let data = JSON.parse(article.data)
			if(get(posts).filter(e => e.id === data.id).length !== 0) {
				let newPosts = get(posts)
				let postIndex = newPosts.findIndex(e => e.id === data.id)
				newPosts[postIndex] = data
				posts.set(newPosts)
			}
		})

		eventSrc.addEventListener('created', (article) => {
			let data = JSON.parse(article.data);
			if (get(selectedTags).length === 0 || data.tags.some((e) => get(selectedTags).includes(e))) {
				let newPosts = get(posts);
				newPosts.unshift(data);
				posts.set(newPosts);
			}
		})
	}
});

export function getNextPage() {
	if(!loading && get(token) !== '') {
		getPostPage(retrievedPages, get(selectedTags));
	}
}


export function deletePost(id) {
	axios
		.delete(`${PUBLIC_TIL_SERVER_URL}/posts/${id}`, { headers: { Authorization: get(token) } })
		.then(() => {
			toast.success(get(t)('post.deleted'), {
				position: 'top-right'
			});
			let postsFilter = get(posts).filter((e) => e.id !== id);
			posts.set(postsFilter);
		})
		.catch((err) => {
			toast.error(get(t)('generic_error', { err: err.response.data }), {
				position: 'top-right'
			});
		});
}

export function createPost(title, link, lang, tags, content) {
	let form = {
		title: title,
		link: link,
		tags: tags,
		content: content
	};

	form.tags.push(`lang:${lang}`);

	return axios
		.post(`${PUBLIC_TIL_SERVER_URL}/posts`, form, { headers: { Authorization: get(token) } })
		.then(() => {
			toast.success(get(t)('post.created'), {
				position: 'top-right'
			});
			getTags();
			goto('/');
		})
		.catch((err) => {
			toast.error(get(t)('generic_error', { err: err.response.data }), {
				position: 'top-right'
			});
		});
}


function getPostPage(page, tags) {
	if (page < totalPages && page >= retrievedPages) {
		loading = true;
		retrievedPages = page + 1;
		axios
			.get(
				`${PUBLIC_TIL_SERVER_URL}/posts?page=${page}${tags && tags.length > 0 ? '&tags=' + tags.join(',') : ''}`,
				{ headers: { Authorization: get(token) } }
			)
			.then((res) => {
				let newPosts = get(posts);
				newPosts.push(...res.data.items);
				posts.set(newPosts);
				totalPages = res.data.total_pages;
			})
			.finally(() => {
				loading = false;
			})
	}
}