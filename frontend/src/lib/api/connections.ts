import { api } from './client';
import type { AccountInfo, Connection, Platform } from '$lib/types';

export function listConnections() {
	return api.get<{ userId: number; connections: Connection[] }>('/api/v1/me');
}

export function disconnect(platform: Platform) {
	return api.delete<void>(`/api/v1/connections/${platform}`);
}

export function accountInfo(platform: Platform) {
	return api.get<AccountInfo>(`/api/v1/connections/${platform}/info`);
}
