<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { auth } from '$lib/auth.svelte';
	import { theme } from '$lib/theme.svelte';
	import Logo from '$lib/components/Logo.svelte';

	let { children } = $props();
	let ready = $state(false);

	onMount(() => {
		if (!auth.isAuthenticated) {
			goto(resolve('/login'));
			return;
		}
		ready = true;
	});

	function logout() {
		auth.logout();
		goto(resolve('/login'));
	}
</script>

{#if ready}
	<div class="min-h-screen">
		<header class="sticky top-0 z-30 border-b border-border/70 bg-(--color-bg)/85 backdrop-blur">
			<div class="mx-auto flex max-w-5xl items-center justify-between gap-6 px-6 py-4">
				<a href={resolve('/dashboard')} aria-label="elci home" class="shrink-0 text-(--color-logo)">
					<Logo />
				</a>

				<nav class="flex items-center gap-4 text-sm font-medium">
					<a
						href={resolve('/dashboard')}
						class="hover:text-(--color-text) {page.url.pathname === resolve('/dashboard')
							? 'text-(--color-text)'
							: 'text-(--color-text-muted)'}"
					>
						Dashboard
					</a>
					<a
						href={resolve('/compose')}
						class="rounded-full bg-(--color-navy) px-4 py-2 text-sm font-semibold text-white transition hover:opacity-90"
					>
						New post
					</a>
					<button
						type="button"
						onclick={() => theme.toggle()}
						aria-label="Toggle color theme"
						class="grid h-9 w-9 place-items-center rounded-full border border-border text-(--color-text-muted) transition hover:text-(--color-text)"
					>
						{#if theme.current === 'dark'}
							<svg
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
								stroke-width="1.7"
								class="h-4 w-4"
							>
								<circle cx="12" cy="12" r="4.5" />
								<path
									d="M12 2.5v2M12 19.5v2M4.2 4.2l1.4 1.4M18.4 18.4l1.4 1.4M2.5 12h2M19.5 12h2M4.2 19.8l1.4-1.4M18.4 5.6l1.4-1.4"
									stroke-linecap="round"
								/>
							</svg>
						{:else}
							<svg viewBox="0 0 24 24" fill="currentColor" class="h-4 w-4">
								<path d="M20 14.5A8.5 8.5 0 0 1 9.5 4a8.5 8.5 0 1 0 10.5 10.5Z" />
							</svg>
						{/if}
					</button>
					<button
						type="button"
						onclick={logout}
						class="text-(--color-text-muted) hover:text-(--color-text)"
					>
						Log out
					</button>
				</nav>
			</div>
		</header>

		<main class="mx-auto max-w-5xl px-6 py-10">
			{@render children()}
		</main>
	</div>
{/if}
