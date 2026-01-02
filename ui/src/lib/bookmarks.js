import {get, writable} from "svelte/store";
import {token} from "$lib/auth.js";
import axios from "axios";
import {PUBLIC_TIL_SERVER_URL} from "$env/static/public";
import toast from "svelte-french-toast";
import {t} from "$lib/i18n.js";

export const  bookmarks = writable([]);

token.subscribe(newToken => {
    if(newToken !== "") {
        getBookmarks()
    }
})


function getBookmarks() {
    axios.get(`${PUBLIC_TIL_SERVER_URL}/bookmarks`, {headers: {'Authorization': get(token)}})
        .then(res => {
            if(res.status===204) {
                bookmarks.set([])
            } else {
                bookmarks.set(res.data)
            }
        })
}

export function addBookmark(bookmarkId){
    axios.put(`${PUBLIC_TIL_SERVER_URL}/bookmarks/${bookmarkId}`, {}, {headers: {'Authorization': get(token)}})
        .then(() => {
            toast.success(get(t)('bookmarks.added'), {
                position: "top-right"
            });
            getBookmarks()
        })
        .catch(err => {
            toast.error(get(t)('bookmarks.generic_error', {err: err.response.data}), {
                position: "top-right"
            });
        })
}

export function removeBookmark(bookmarkId){
    axios.delete(`${PUBLIC_TIL_SERVER_URL}/bookmarks/${bookmarkId}`, {headers: {'Authorization': get(token)}})
        .then(() => {
            toast.success(get(t)('bookmarks.deleted'), {
                position: "top-right"
            });
            getBookmarks()
        })
        .catch(err => {
            toast.error(get(t)('bookmarks.generic_error', {err: err.response.data}), {
                position: "top-right"
            });
        })
}