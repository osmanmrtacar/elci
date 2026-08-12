<script lang="ts">
	import { theme } from '$lib/theme.svelte';
	import { auth } from '$lib/auth.svelte';
	import { resolve } from '$app/paths';

	let menuOpen = $state(false);

	const links = [
		{ href: '#product', label: 'Product' },
		{ href: '#how-it-works', label: 'How it works' },
		{ href: '#agents', label: 'For agents' },
		{ href: '#faq', label: 'FAQ' }
	];

	function closeMenu() {
		menuOpen = false;
	}
</script>

<div class="pointer-events-none fixed inset-x-0 top-0 z-40 mt-6 flex justify-center px-4">
	<header
		class="pointer-events-auto w-max max-w-[calc(100vw-2rem)] rounded-full border border-border bg-(--color-surface)/80 shadow-[0_8px_30px_-15px_rgba(16,24,40,0.25)] backdrop-blur-xl transition-all duration-700 ease-[cubic-bezier(0.32,0.72,0,1)]"
	>
		<div class="flex items-center gap-6 px-4 py-2.5 md:gap-8 md:px-5">
			<a href={resolve('/')} class="flex items-center gap-2 text-base font-semibold">
				elci
				<span class="inline-block h-1.5 w-1.5 rounded-full bg-(--color-accent)"></span>
			</a>

			<nav class="hidden items-center gap-6 text-sm font-medium text-(--color-text-muted) md:flex">
				{#each links as link (link.href)}
					<a href={link.href} class="pill-nav-link hover:text-(--color-text)">{link.label}</a>
				{/each}
			</nav>

			<div class="hidden items-center gap-3 md:flex">
				<button
					type="button"
					onclick={() => theme.toggle()}
					aria-label="Toggle color theme"
					class="grid h-8 w-8 place-items-center rounded-full border border-border text-(--color-text-muted) transition hover:text-(--color-text)"
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
				<a
					href={auth.isAuthenticated ? resolve('/dashboard') : resolve('/login')}
					class="rounded-full bg-(--color-navy) px-4 py-2 text-sm font-semibold text-white transition hover:opacity-90"
				>
					{auth.isAuthenticated ? 'Dashboard' : 'Get started'}
				</a>
			</div>

			<button
				type="button"
				class="relative grid h-8 w-8 place-items-center md:hidden"
				aria-label={menuOpen ? 'Close menu' : 'Open menu'}
				aria-expanded={menuOpen}
				onclick={() => (menuOpen = !menuOpen)}
			>
				<span
					class="hamburger-line absolute h-[1.5px] w-4 bg-(--color-text)"
					style={menuOpen
						? 'transform: translateY(0) rotate(45deg)'
						: 'transform: translateY(-4px)'}
				></span>
				<span
					class="hamburger-line absolute h-[1.5px] w-4 bg-(--color-text)"
					style={menuOpen
						? 'transform: translateY(0) rotate(-45deg)'
						: 'transform: translateY(4px)'}
				></span>
			</button>
		</div>
	</header>
</div>

<div
	class="fixed inset-0 z-30 bg-(--color-bg)/95 backdrop-blur-3xl transition-opacity duration-700 ease-[cubic-bezier(0.32,0.72,0,1)] md:hidden {menuOpen
		? 'pointer-events-auto opacity-100'
		: 'pointer-events-none opacity-0'}"
	aria-hidden={!menuOpen}
>
	<nav
		class="flex h-full flex-col items-center justify-center gap-8 text-2xl font-medium {menuOpen
			? 'nav-open'
			: ''}"
	>
		{#each links as link, i (link.href)}
			<a
				href={link.href}
				class="nav-word"
				style="transition-delay: {100 + i * 50}ms"
				onclick={closeMenu}
				tabindex={menuOpen ? 0 : -1}
			>
				{link.label}
			</a>
		{/each}
		<a
			href={auth.isAuthenticated ? resolve('/dashboard') : resolve('/login')}
			class="nav-word mt-4 rounded-full bg-(--color-accent) px-6 py-3 text-base font-semibold text-white"
			style="transition-delay: {100 + links.length * 50}ms"
			onclick={closeMenu}
			tabindex={menuOpen ? 0 : -1}
		>
			{auth.isAuthenticated ? 'Dashboard' : 'Connect your accounts'}
		</a>
	</nav>
</div>
