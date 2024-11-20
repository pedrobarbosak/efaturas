import {get} from "svelte/store";
import { user } from '$lib/store';
import { redirect } from '@sveltejs/kit';

export const ssr = false;
export const prerender = true;

export async function load({url}) {
    const data = get(user)

    if (!data.isLoggedIn && url.pathname !== "/") {
        throw redirect(302, '/');
    }

    return {};
}
