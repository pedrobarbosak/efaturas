import {redirect} from "@sveltejs/kit";
import {browser} from "$app/environment";
import {goto} from "$app/navigation";

export default class HttpClient {
	private readonly baseUrl: string;

	constructor(baseUrl: string) {
		this.baseUrl = baseUrl;
	}

	private async request<T>(
		url: string,
		options: RequestInit
	): Promise<{ hasError: boolean; data: T | null; error: string | null; status: number }> {
		try {
			const response = await fetch(`${this.baseUrl}${url}`, options);
			if (!response.ok) {
				if (response.status == 401) {
					console.log("OUT!")
					if (window) {
						window.location.href = "/"
					} else {
						redirect(302, "/")
					}
				}

				const errorText = await response.text();
				return { hasError: true, data: null, error: errorText || response.statusText || "Request failed: " + response.status, status: response.status }; // You can replace this with a more appropriate error handling logic.
			}
			const data = await response.json();
			return { hasError: false, data: data as T, error: null, status: response.status };
		} catch (error) {
			console.error('Request failed', error);
			return { hasError: true, data: null, error: 'Request failed', status: 500 };
		}
	}

	public async get<T>(
		url: string,
		headers: Record<string, string> = {}
	): Promise<{ hasError: boolean; data: T | null; error: string | null; status: number }> {
		return this.request<T>(url, {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json',
				...headers
			}
		});
	}

	public async post<K, T>(
		url: string,
		body: K,
		headers: Record<string, string> = {}
	): Promise<{ hasError: boolean; data: T | null; error: string | null; status: number }> {
		return this.request<T>(url, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
				...headers
			},
			body: JSON.stringify(body)
		});
	}

	public async put<K, T>(
		url: string,
		body: K,
		headers: Record<string, string> = {}
	): Promise<{ hasError: boolean; data: T | null; error: string | null }> {
		let b = JSON.stringify(body)
		if (body instanceof Map) {
			b = JSON.stringify(Object.fromEntries(body))
		}

		return this.request<T>(url, {
			method: 'PUT',
			headers: {
				'Content-Type': 'application/json',
				...headers
			},
			body: b
		});
	}

	public async delete<T>(
		url: string,
		headers: Record<string, string> = {}
	): Promise<{ hasError: boolean; data: T | null; error: string | null }> {
		return this.request<T>(url, {
			method: 'DELETE',
			headers: {
				'Content-Type': 'application/json',
				...headers
			}
		});
	}
}
