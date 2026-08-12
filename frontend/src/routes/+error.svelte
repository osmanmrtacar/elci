<script lang="ts">
	import { page } from '$app/state';
	import { resolve } from '$app/paths';
</script>

<svelte:head>
	<title>{page.status} — elci</title>
</svelte:head>

<main class="flex min-h-screen flex-col items-center justify-center px-6 text-center">
	<p class="font-mono text-sm text-(--color-accent)">{page.status}</p>
	<h1 class="mt-3 text-2xl font-bold tracking-tight">
		{page.status === 404 ? "This page doesn't exist" : 'Something went wrong'}
	</h1>
	<p class="mt-2 max-w-sm text-sm text-(--color-text-muted)">
		{page.error?.message ?? 'Try heading back to the homepage.'}
	</p>
	<a
		href={resolve('/')}
		class="mt-6 rounded-full bg-(--color-navy) px-5 py-2.5 text-sm font-semibold text-white hover:opacity-90"
	>
		Back to elci
	</a>
</main>
