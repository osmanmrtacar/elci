<script lang="ts">
	import { onMount } from 'svelte';
	import { resolve } from '$app/paths';
	import { listConnections, disconnect } from '$lib/api/connections';
	import { listPosts } from '$lib/api/posts';
	import { api } from '$lib/api/client';
	import { providers, getProvider } from '$lib/providers/registry';
	import StatusPill from '$lib/components/StatusPill.svelte';
	import type { Connection, Post } from '$lib/types';

	let connections = $state<Connection[]>([]);
	let posts = $state<Post[]>([]);
	let loading = $state(true);
	let error = $state<string | null>(null);

	async function load() {
		loading = true;
		error = null;
		try {
			const [me, postList] = await Promise.all([listConnections(), listPosts()]);
			connections = me.connections;
			posts = postList;
		} catch (e) {
			error = e instanceof Error ? e.message : 'Could not load your account';
		} finally {
			loading = false;
		}
	}

	onMount(load);

	function connectionFor(id: string) {
		return connections.find((c) => c.platform === id && c.isActive);
	}

	async function handleDisconnect(id: 'tiktok' | 'instagram' | 'x') {
		if (!confirm(`Disconnect ${id}? You'll need to reconnect to post there again.`)) return;
		await disconnect(id);
		await load();
	}

	function formatDate(iso: string | null) {
		if (!iso) return '';
		return new Date(iso).toLocaleString(undefined, {
			month: 'short',
			day: 'numeric',
			hour: 'numeric',
			minute: '2-digit'
		});
	}
</script>

<svelte:head>
	<title>Dashboard — elci</title>
</svelte:head>

<section class="mb-10">
	<h1 class="mb-4 font-[family-name:var(--font-display)] text-xl font-bold">Connected accounts</h1>
	<div class="grid gap-3 sm:grid-cols-3">
		{#each providers as provider (provider.id)}
			{@const conn = connectionFor(provider.id)}
			<div class="flex items-center justify-between rounded-xl border border-border p-4">
				<div class="flex items-center gap-2">
					<provider.icon class="h-4 w-4 shrink-0" />
					<span class="leading-tight">
						<span class="block text-sm font-semibold">{provider.label}</span>
						{#if conn}
							<span class="block font-mono text-[11px] text-(--color-text-muted)"
								>@{conn.username}</span
							>
						{/if}
					</span>
				</div>
				{#if conn}
					<button
						type="button"
						onclick={() => handleDisconnect(provider.id)}
						class="text-xs text-(--color-text-muted) hover:text-(--color-status-failed)"
					>
						Disconnect
					</button>
				{:else}
					<!-- eslint-disable-next-line svelte/no-navigation-without-resolve -- external backend OAuth URL, not a SvelteKit route -->
					<a href={api.authURL(provider.id)} class="text-xs font-semibold text-(--color-accent)">
						Connect
					</a>
				{/if}
			</div>
		{/each}
	</div>
</section>

<section>
	<div class="mb-4 flex items-center justify-between">
		<h1 class="font-[family-name:var(--font-display)] text-xl font-bold">Posts</h1>
		<a
			href={resolve('/compose')}
			class="rounded-full bg-(--color-navy) px-4 py-2 text-sm font-semibold text-white transition hover:opacity-90"
		>
			New post
		</a>
	</div>

	{#if loading}
		<p class="text-sm text-(--color-text-muted)">Loading…</p>
	{:else if error}
		<p class="text-sm text-(--color-status-failed)">{error}</p>
	{:else if posts.length === 0}
		<p
			class="rounded-xl border border-dashed border-border p-8 text-center text-sm text-(--color-text-muted)"
		>
			No posts yet — create your first one.
		</p>
	{:else}
		<div class="divide-y divide-border rounded-xl border border-border">
			{#each posts as post (post.id)}
				<a
					href={resolve('/(app)/posts/[id]', { id: String(post.id) })}
					class="flex items-center justify-between gap-4 p-4 hover:bg-(--color-border)/20"
				>
					<div class="min-w-0 flex-1">
						<p class="truncate text-sm">{post.defaultCaption || '(no caption)'}</p>
						<div class="mt-1 flex items-center gap-2 text-xs text-(--color-text-muted)">
							{#each post.targets as target (target.id)}
								{@const p = getProvider(target.platform)}
								<p.icon class="h-3.5 w-3.5" />
							{/each}
							<span>
								{post.scheduledAt ? formatDate(post.scheduledAt) : formatDate(post.createdAt)}
							</span>
						</div>
					</div>
					<StatusPill status={post.status} />
				</a>
			{/each}
		</div>
	{/if}
</section>
