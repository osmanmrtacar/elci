export type Platform = 'tiktok' | 'instagram' | 'x';

export type MediaKind = 'image' | 'video';

export interface Connection {
	platform: Platform;
	username: string;
	displayName: string;
	avatarUrl: string;
	isActive: boolean;
	connectedAt: string;
}

export interface AccountInfo {
	platformUserId: string;
	username: string;
	displayName: string;
	avatarUrl: string;
	extra: Record<string, unknown>;
}

export type TargetStatus = 'pending' | 'processing' | 'published' | 'failed';

export interface PostTarget {
	id: number;
	platform: Platform;
	captionOverride: string | null;
	mediaUrlsOverride: string[] | null;
	settings: Record<string, unknown>;
	status: TargetStatus;
	platformPostId: string;
	error: string;
	publishedAt: string | null;
}

export type PostStatus = 'draft' | 'scheduled' | 'publishing' | 'published' | 'partial' | 'failed';

export interface Post {
	id: number;
	defaultCaption: string;
	defaultMediaKind: MediaKind | null;
	defaultMediaUrls: string[] | null;
	status: PostStatus;
	scheduledAt: string | null;
	createdAt: string;
	targets: PostTarget[];
}

export interface ValidationErrors {
	errors: Record<Platform, { field: string; message: string }[]>;
}
