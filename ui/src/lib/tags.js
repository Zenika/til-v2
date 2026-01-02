import axios from "axios";
import {get, writable} from "svelte/store";
import {PUBLIC_TIL_SERVER_URL} from "$env/static/public";
import {token} from "$lib/auth.js";

export let tags = writable([]);


token.subscribe(newToken => {
    if (newToken !== "") {
        getTags()
    }
})

export function getTags() {
    axios.get(`${PUBLIC_TIL_SERVER_URL}/tags`, {headers: {Authorization: get(token)}})
        .then(res => {
            tags.set(res.data.filter(e => !e.startsWith("lang:")));
        })
}