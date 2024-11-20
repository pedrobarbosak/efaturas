import { writable } from 'svelte/store';
import { redirect } from "@sveltejs/kit";
import { browser } from '$app/environment';

export let ssr = false;

export const formKey = "form"

const userKey = "session"
const userEmpty = { isLoggedIn: false, username: "", password: "", token: "" };

export const user = persistentStore(userKey, userEmpty);

export async function logout() {
    user.set(userEmpty)
    redirect(302, "/")
}

export async function logoutAndClear() {
    user.set(userEmpty)
    localStorage.removeItem(userKey)
    localStorage.removeItem(formKey)
}

function persistentStore(key: string, defaultValue: any) {
    const existing = browser ? localStorage.getItem(key) : null;
    const store = writable(existing ? JSON.parse(existing) : defaultValue);

    if (browser) {
        store.subscribe(value => {
            localStorage.setItem(key, JSON.stringify(value));
        });
    }

    return store;
}