import { env } from '$env/dynamic/public';
import { api } from './client';
import type { MediaKind } from '$lib/types';

interface PresignedUpload {
	key: string;
	uploadUrl: string;
	publicUrl: string;
}

// Reads a <video>'s duration client-side so the backend can enforce
// platform duration limits (e.g. TikTok's max_video_post_duration_sec)
// without re-processing the file itself.
function videoDurationSeconds(file: File): Promise<number> {
	return new Promise((resolve, reject) => {
		const url = URL.createObjectURL(file);
		const videoEl = document.createElement('video');
		videoEl.preload = 'metadata';
		videoEl.style.display = 'none';

		const cleanup = () => {
			URL.revokeObjectURL(url);
			videoEl.remove();
		};
		videoEl.onloadedmetadata = () => {
			const duration = Math.round(videoEl.duration);
			cleanup();
			resolve(duration);
		};
		videoEl.onerror = () => {
			cleanup();
			reject(new Error('could not read video metadata'));
		};

		// Metadata loading isn't reliable across browsers for a <video> that
		// was never attached to the document, so it must be in the DOM (even
		// if hidden) for `loadedmetadata` to fire consistently.
		document.body.appendChild(videoEl);
		videoEl.src = url;
	});
}

export async function uploadMedia(
	file: File,
	kind: MediaKind
): Promise<{ id: number; publicUrl: string; kind: MediaKind }> {
	const presigned = await api.post<PresignedUpload>('/api/v1/media/presign', {
		filename: file.name,
		contentType: file.type
	});

	const putRes = await fetch(presigned.uploadUrl, {
		method: 'PUT',
		headers: { 'Content-Type': file.type },
		body: file
	});
	if (!putRes.ok) throw new Error('upload to storage failed');

	const durationSeconds = kind === 'video' ? await videoDurationSeconds(file) : undefined;

	return api.post('/api/v1/media/confirm', {
		key: presigned.key,
		publicUrl: presigned.publicUrl,
		kind,
		contentType: file.type,
		sizeBytes: file.size,
		durationSeconds
	});
}

export const publicApiUrl = env.PUBLIC_API_URL ?? 'http://localhost:8080';
