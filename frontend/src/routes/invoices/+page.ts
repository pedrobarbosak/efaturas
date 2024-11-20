import {logout, user} from "@/store";
import { get } from 'svelte/store';
import {getCategories, getInvoices} from "@/api";
import {initInvoices} from "@/utils";

export const ssr = false;
export const prerender = true;

export async function load({ params }) {
    const data = get(user)

    const categories = await getCategories(data.token)
    if (categories.hasError || categories.status != 200) {
        await logout()
    }

    const invoices = await getInvoices(data.token)
    if (invoices.hasError || invoices.status != 200) {
        await logout()
    }

    return {
        invoices: initInvoices(invoices.data || []),
        categories: categories.data || [],
    };
}
