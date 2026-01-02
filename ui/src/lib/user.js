import axios from "axios";
import {PUBLIC_TIL_SERVER_URL} from "$env/static/public";
import {get} from "svelte/store";
import {currentUser, refreshUser, token} from "$lib/auth.js";
import toast from "svelte-french-toast";
import {t} from "$lib/i18n.js";

export function renewFeedKey() {
    axios.put(`${PUBLIC_TIL_SERVER_URL}/users/self/renew`, {}, {headers: {Authorization: get(token)}})
        .then(() => {
            toast.success(get(t)('profile.renewed_feed_key'), {
                position: "top-right"
            });
            refreshUser()
        })
        .catch(err => {
            toast.error(get(t)('generic_error', {err: err.response.data}), {
                position: "top-right"
            });
        });
}

export function updateTags(tags) {
    let u = get(currentUser)
    u.automatic_tags_filter = tags

    axios.put(`${PUBLIC_TIL_SERVER_URL}/users/self`, u, {headers: {Authorization: get(token)}})
        .then(() => {
            toast.success(get(t)('profile.updated'), {
                position: "top-right"
            });
            refreshUser()
        })
        .catch(err => {
            toast.error(get(t)('generic_error', {err: err.response.data}), {
                position: "top-right"
            });
        });

}