<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { deletePost, getPost, publishPost } from '$lib/api/posts';
	import { getProvider } from '$lib/providers/registry';
	import { ApiError } from '$lib/api/client';
	import StatusPill from '$lib/components/StatusPill.svelte';
	import MediaCarousel from '$lib/providers/MediaCarousel.svelte';
	import type { Post, ValidationErrors } from '$lib/types';

	function describeError(e: unknown, fallback: string): string {
		if (
			e instanceof ApiError &&
			e.status === 422 &&
			e.body &&
			typeof e.body === 'object' &&
			'errors' in e.body
		) {
			const errors = (e.body as ValidationErrors).errors;
			return Object.values(errors)
				.flat()
				.map((err) => err.message)
				.join(' ');
		}
		return e instanceof Error ? e.message : fallback;
	}

	const id = Number(page.params.id);

	let post = $state<Post | null>(null);
	let loading = $state(true);
	let error = $state<string | null>(null);
	let acting = $state(false);
	let pollHandle: ReturnType<typeof setInterval> | undefined;

	async function load() {
		try {
			post = await getPost(id);
		} catch (e) {
			error = e instanceof Error ? e.message : 'Could not load this post';
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		load();
		pollHandle = setInterval(() => {
			if (
				post &&
				(post.status === 'publishing' || post.targets.some((t) => t.status === 'processing'))
			) {
				load();
			}
		}, 5000);
	});

	onDestroy(() => clearInterval(pollHandle));

	async function handlePublish() {
		acting = true;
		try {
			await publishPost(id);
			await load();
		} catch (e) {
			error = describeError(e, 'Could not publish this post');
		} finally {
			acting = false;
		}
	}

	async function handleDelete() {
		if (!confirm('Delete this post? This cannot be undone.')) return;
		acting = true;
		try {
			await deletePost(id);
			goto(resolve('/dashboard'));
		} catch (e) {
			error = e instanceof Error ? e.message : 'Could not delete this post';
			acting = false;
		}
	}

	function formatDate(iso: string | null) {
		if (!iso) return '—';
		return new Date(iso).toLocaleString();
	}
</script>

<svelte:head>
	<title>Post — elci</title>
</svelte:head>

{#if loading}
	<p class="text-sm text-(--color-text-muted)">Loading…</p>
{:else if error && !post}
	<p class="text-sm text-(--color-status-failed)">{error}</p>
{:else if post}
	<div class="mb-6 flex items-center justify-between">
		<h1 class="font-[family-name:var(--font-display)] text-xl font-bold">Post</h1>
		<StatusPill status={post.status} />
	</div>

	<div class="mb-8 rounded-xl border border-border p-5">
		<p class="text-sm whitespace-pre-wrap">{post.defaultCaption || '(no caption)'}</p>
		{#if post.defaultMediaUrls?.[0]}
			<div class="mt-4 max-w-xs overflow-hidden rounded-lg border border-border">
				{#if post.defaultMediaKind === 'video'}
					<!-- svelte-ignore a11y_media_has_caption -->
					<video src={post.defaultMediaUrls[0]} class="w-full" controls></video>
				{:else}
					<div class="relative aspect-square w-full">
						<MediaCarousel mediaUrls={post.defaultMediaUrls} mediaKind={post.defaultMediaKind} />
					</div>
				{/if}
			</div>
		{/if}
		<p class="mt-4 text-xs text-(--color-text-muted)">
			{post.scheduledAt
				? `Scheduled for ${formatDate(post.scheduledAt)}`
				: `Created ${formatDate(post.createdAt)}`}
		</p>
	</div>

	<div class="mb-8 space-y-3">
		{#each post.targets as target (target.id)}
			{@const provider = getProvider(target.platform)}
			<div class="rounded-xl border border-border p-4">
				<div class="flex items-center justify-between">
					<div class="flex items-center gap-2 text-sm font-semibold">
						<provider.icon class="h-4 w-4" />
						{provider.label}
					</div>
					<StatusPill status={target.status} />
				</div>
				{#if target.captionOverride}
					<p class="mt-2 text-sm text-(--color-text-muted)">{target.captionOverride}</p>
				{/if}
				{#if target.error}
					<p class="mt-2 text-sm text-(--color-status-failed)">{target.error}</p>
				{/if}
				{#if target.platform === 'tiktok' && target.status === 'processing'}
					<p class="mt-2 text-xs text-(--color-text-muted)">
						TikTok is still processing this — it can take a few minutes to appear on the profile.
					</p>
				{/if}
				{#if target.publishedAt}
					<p class="mt-2 text-xs text-(--color-text-muted)">
						Published {formatDate(target.publishedAt)}
					</p>
				{/if}
			</div>
		{/each}
	</div>

	{#if error}
		<p class="mb-4 text-sm text-(--color-status-failed)">{error}</p>
	{/if}

	<div class="flex gap-3 border-t border-border pt-6">
		{#if post.status === 'draft' || post.status === 'scheduled'}
			<button
				type="button"
				disabled={acting}
				onclick={handlePublish}
				class="rounded-full bg-(--color-navy) px-4 py-2 text-sm font-semibold text-white disabled:opacity-50"
			>
				Publish now
			</button>
		{/if}
		<button
			type="button"
			disabled={acting}
			onclick={handleDelete}
			class="rounded-full border border-border px-4 py-2 text-sm font-semibold text-(--color-status-failed) disabled:opacity-50"
		>
			Delete
		</button>
	</div>
{/if}
