import {derived, writable} from 'svelte/store';
import translations from '../resources/translations.js';

if(!localStorage.getItem('lang')) {
    localStorage.setItem('lang', 'fr');
}


export const locale = writable(localStorage.getItem('lang'));
export const locales = Object.keys(translations);

locale.subscribe(lang => {
    localStorage.setItem("lang", lang)
})

function translate(locale, key, vars) {
    // Let's throw some errors if we're trying to use keys/locales that don't exist.
    if (!key) throw new Error('no key provided to $t()');
    if (!locale) throw new Error(`no translation for key "${key}"`);

    let path = key.split(".")

    let finalText = translations[locale]
    for (let i in path) {
        if(finalText === undefined) {
            break
        }
        finalText = finalText[path[i]]
    }

    if (!finalText) {
        finalText = `${locale}.${key}`;
        console.warn(`Missing translation for ${locale}.${key}`)
        if(finalText.includes(".error.401") || finalText.includes(".error.403")) {
            finalText = translations[locale]["generics"]["unauthorized"];
        } else if(finalText.includes(".error.")) {
            finalText = translations[locale]["generics"]["generic_error"];
        }
    }

    // Replace any passed in variables in the translation string.
    Object.keys(vars).map((k) => {
        const regex = new RegExp(`{{${k}}}`, 'g');
        finalText = finalText.replace(regex, vars[k]);
    });

    return finalText;
}

export const t = derived(
    locale,
    ($locale) =>
        (key, vars = {}) =>
            translate($locale, key, vars)
);
