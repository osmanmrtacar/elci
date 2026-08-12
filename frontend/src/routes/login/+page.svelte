<script lang="ts">
	import { resolve } from '$app/paths';
	import { providers } from '$lib/providers/registry';
	import { api } from '$lib/api/client';
	import { page } from '$app/state';

	const errorMessages: Record<string, string> = {
		missing_code_or_state: 'That sign-in link looks incomplete. Please try again.',
		connection_failed: 'We could not complete that connection. Please try again.'
	};
	let error = $derived(page.url.searchParams.get('error'));
</script>

<svelte:head>
	<title>Sign in — elci</title>
</svelte:head>

<main class="flex min-h-screen flex-col items-center justify-center px-6 text-center">
	<a
		href={resolve('/')}
		class="flex items-center gap-2 font-[family-name:var(--font-display)] text-lg font-bold"
	>
		elci
		<span class="inline-block h-1.5 w-1.5 rounded-full bg-(--color-accent)"></span>
	</a>
	<h1 class="mt-6 text-2xl font-bold tracking-tight">Sign in to elci</h1>
	<p class="mt-2 max-w-sm text-sm text-(--color-text-muted)">
		Connect a platform account to sign in — we'll create your elci account automatically the first
		time.
	</p>

	{#if error}
		<p
			class="mt-4 max-w-sm rounded-lg bg-(--color-status-failed)/10 p-3 text-sm text-(--color-status-failed)"
		>
			{errorMessages[error] ?? 'Something went wrong. Please try again.'}
		</p>
	{/if}

	<div class="mt-8 flex w-full max-w-xs flex-col gap-3">
		{#each providers as provider (provider.id)}
			<!-- eslint-disable svelte/no-navigation-without-resolve -- external backend OAuth URL, not a SvelteKit route -->
			<a
				href={api.authURL(provider.id)}
				class="flex items-center justify-center gap-2 rounded-full border border-border bg-(--color-surface) px-5 py-2.5 text-sm font-semibold transition hover:opacity-80"
			>
				<provider.icon class="h-4 w-4" />
				Continue with {provider.label}
			</a>
			<!-- eslint-enable svelte/no-navigation-without-resolve -->
		{/each}
	</div>

	<a href={resolve('/')} class="mt-6 text-sm text-(--color-text-muted) hover:text-(--color-text)">
		Cancel
	</a>
</main>
