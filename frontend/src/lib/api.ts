import HttpClient from './http';
import type { Category } from './models/models';
import {Invoice} from "./models/models";
import {fetchEventSource} from "@microsoft/fetch-event-source";

const baseURL = 'http://127.0.0.1:8080/v1';

const api = new HttpClient(baseURL);

function auth(token: string) {
	return { "Authorization": token }
}

export function login(username: string, password: string) {
	return api.post('/login', { username, password });
}

export function getCategories(token: string) {
	return api.get<Category[]>('/categories', auth(token));
}

export function getInvoices(token: string) {
	return api.get<Invoice[]>('/invoices', auth(token));
}

export function fetchNewInvoices(token: string) {
	return api.post<{},Invoice[]>('/invoices',{}, auth(token));
}

export function processInvoices(token: string) {
	return api.put<{},Invoice[]>('/invoices',{}, auth(token));
}

export function updateInvoices(token: string, invoices: Invoice[]) {
	let body = new Map<number, string>()
	for (let i=0; i<invoices.length; i++) {
		body.set(invoices[i].id, invoices[i].selected)
	}

	return api.put<Map<number, string>, Invoice[]>('/invoices/categories', body, auth(token));
}

export async function getInvoiceSSE(token: string, invoiceID: number, onmessage) {
	 await fetchEventSource(`${baseURL}/invoices/${invoiceID}`, {
		 method: 'GET',
		 headers: auth(token),
		 onmessage: onmessage,
	 });
}