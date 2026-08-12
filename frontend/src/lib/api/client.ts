import { env } from '$env/dynamic/public';
import { auth } from '$lib/auth.svelte';

const API_URL = env.PUBLIC_API_URL ?? 'http://localhost:8080';

export class ApiError extends Error {
	constructor(
		public status: number,
		public body: unknown
	) {
		super(
			typeof body === 'object' && body && 'error' in body ? String(body.error) : `HTTP ${status}`
		);
	}
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
	const headers = new Headers(init?.headers);
	if (auth.token) headers.set('Authorization', `Bearer ${auth.token}`);
	if (init?.body) headers.set('Content-Type', 'application/json');

	const res = await fetch(`${API_URL}${path}`, { ...init, headers });
	const raw = await res.text();
	const body = raw ? JSON.parse(raw) : undefined;

	if (!res.ok) throw new ApiError(res.status, body);
	return body as T;
}

export const api = {
	get: <T>(path: string) => request<T>(path),
	post: <T>(path: string, body?: unknown) =>
		request<T>(path, { method: 'POST', body: body ? JSON.stringify(body) : undefined }),
	delete: <T>(path: string) => request<T>(path, { method: 'DELETE' }),
	authURL: (platform: string) =>
		`${API_URL}/api/v1/auth/${platform}/login${auth.token ? `?token=${encodeURIComponent(auth.token)}` : ''}`
};
