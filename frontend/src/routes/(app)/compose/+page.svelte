<script lang="ts">
	import { onMount } from 'svelte';
	import Composer from '$lib/composer/Composer.svelte';
	import { listConnections } from '$lib/api/connections';
	import type { Connection } from '$lib/types';

	let connections = $state<Connection[]>([]);
	let loading = $state(true);
	let error = $state<string | null>(null);

	onMount(async () => {
		try {
			const me = await listConnections();
			connections = me.connections;
		} catch (e) {
			error = e instanceof Error ? e.message : 'Could not load your connections';
		} finally {
			loading = false;
		}
	});
</script>

<svelte:head>
	<title>New post — elci</title>
</svelte:head>

<h1 class="mb-6 font-[family-name:var(--font-display)] text-xl font-bold">New post</h1>

{#if loading}
	<p class="text-sm text-(--color-text-muted)">Loading…</p>
{:else if error}
	<p class="text-sm text-(--color-status-failed)">{error}</p>
{:else if !connections.some((c) => c.isActive)}
	<p
		class="rounded-xl border border-dashed border-border p-8 text-center text-sm text-(--color-text-muted)"
	>
		Connect at least one platform from the dashboard before composing a post.
	</p>
{:else}
	<Composer {connections} />
{/if}
