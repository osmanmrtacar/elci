import { api } from './client';
import type { MediaKind, Platform, Post } from '$lib/types';

export interface TargetInput {
	platform: Platform;
	captionOverride?: string | null;
	settings?: Record<string, unknown>;
}

export interface CreatePostInput {
	defaultCaption: string;
	defaultMediaKind?: MediaKind | null;
	defaultMediaUrls?: string[];
	targets: TargetInput[];
}

export function createPost(input: CreatePostInput) {
	return api.post<Post>('/api/v1/posts', input);
}

export function listPosts() {
	return api.get<Post[]>('/api/v1/posts');
}

export function getPost(id: number) {
	return api.get<Post>(`/api/v1/posts/${id}`);
}

export function schedulePost(id: number, scheduledAt: Date) {
	return api.post<void>(`/api/v1/posts/${id}/schedule`, { scheduledAt: scheduledAt.toISOString() });
}

export function publishPost(id: number) {
	return api.post<void>(`/api/v1/posts/${id}/publish`);
}

export function deletePost(id: number) {
	return api.delete<void>(`/api/v1/posts/${id}`);
}
