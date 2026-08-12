<script lang="ts">
	import { theme } from '$lib/theme.svelte';
	import { auth } from '$lib/auth.svelte';
	import { resolve } from '$app/paths';
	import { track } from '$lib/analytics';
	import { onMount } from 'svelte';

	let mobileOpen = $state(false);
	let scrolled = $state(false);

	onMount(() => {
		const onScroll = () => {
			scrolled = window.scrollY > 8;
		};
		onScroll();
		window.addEventListener('scroll', onScroll, { passive: true });
		return () => window.removeEventListener('scroll', onScroll);
	});

	$effect(() => {
		if (typeof document !== 'undefined') {
			document.body.style.overflow = mobileOpen ? 'hidden' : '';
		}
	});

	function closeMenu() {
		mobileOpen = false;
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape') closeMenu();
	}
</script>

<header
	aria-label="Primary"
	class="sticky top-0 z-30 border-b border-border/70 bg-(--color-bg)/85 backdrop-blur {scrolled
		? 'shadow-sm'
		: ''}"
>
	<div class="mx-auto flex max-w-6xl items-center justify-between gap-6 px-6 py-4">
		<a
			href={resolve('/')}
			aria-label="elci home"
			class="flex shrink-0 items-center gap-[0.3em] text-(--color-logo)"
		>
			<span
				class="font-[family-name:var(--font-logo)] text-[20px] leading-none tracking-[-0.005em]"
				style="font-weight: 480; font-optical-sizing: auto;">elci</span
			>
			<svg
				viewBox="2 14 96 56"
				class="mt-[0.06em] h-[0.58em] w-auto shrink-0 origin-[50%_55%] self-start rotate-[-6deg]"
				aria-hidden="true"
				focusable="false">
				<path
					d="M 12 58 Q 32 25 50 48 Q 68 25 88 58"
					fill="none"
					stroke="currentColor"
					stroke-width="8"
					stroke-linecap="round"
					stroke-linejoin="round"
				/>
				<circle cx="50" cy="46" r="3.6" fill="currentColor" />
			</svg>
		</a>
		<nav aria-label="Primary" class="hidden items-center gap-8 text-sm font-medium text-(--color-text-muted) md:flex">
			<ul class="flex items-center gap-8">
				<li><a href="#product" class="hover:text-(--color-text)">Product</a></li>
				<li><a href="#how-it-works" class="hover:text-(--color-text)">How it works</a></li>
				<li><a href="#faq" class="hover:text-(--color-text)">FAQ</a></li>
			</ul>
		</nav>

		<div class="flex items-center gap-3">
			<button
				type="button"
				onclick={() => theme.toggle()}
				aria-label="Toggle color theme"
				aria-pressed={theme.current === 'dark'}
				class="grid h-9 w-9 place-items-center rounded-full border border-border text-(--color-text-muted) transition hover:text-(--color-text)"
			>
				{#if theme.current === 'dark'}
					<svg
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="1.7"
						class="h-4 w-4"
						aria-hidden="true"
						focusable="false"
					>
						<circle cx="12" cy="12" r="4.5" />
						<path
							d="M12 2.5v2M12 19.5v2M4.2 4.2l1.4 1.4M18.4 18.4l1.4 1.4M2.5 12h2M19.5 12h2M4.2 19.8l1.4-1.4M18.4 5.6l1.4-1.4"
							stroke-linecap="round"
						/>
					</svg>
				{:else}
					<svg
						viewBox="0 0 24 24"
						fill="currentColor"
						class="h-4 w-4"
						aria-hidden="true"
						focusable="false"
					>
						<path d="M20 14.5A8.5 8.5 0 0 1 9.5 4a8.5 8.5 0 1 0 10.5 10.5Z" />
					</svg>
				{/if}
			</button>
			<a
				href={auth.isAuthenticated ? resolve('/dashboard') : resolve('/login')}
				data-analytics="header-cta"
				onclick={() => track('header-cta')}
				class="rounded-full bg-(--color-navy) px-4 py-2 text-sm font-semibold text-white transition hover:opacity-90"
			>
				{auth.isAuthenticated ? 'Dashboard' : 'Get started'}
			</a>
			<button
				type="button"
				class="grid h-9 w-9 place-items-center rounded-full border border-border text-(--color-text-muted) md:hidden"
				aria-label="Open navigation menu"
				aria-expanded={mobileOpen}
				aria-controls="mobile-nav"
				onclick={() => (mobileOpen = !mobileOpen)}
			>
				<svg
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="1.8"
					class="h-5 w-5"
					aria-hidden="true"
					focusable="false"
				>
					{#if mobileOpen}
						<path d="M6 6 18 18M18 6 6 18" stroke-linecap="round" />
					{:else}
						<path d="M4 7h16M4 12h16M4 17h16" stroke-linecap="round" />
					{/if}
				</svg>
			</button>
		</div>
	</div>

	{#if mobileOpen}
		<!-- svelte-ignore a11y_interactive_supports_focus -->
		<div
			id="mobile-nav"
			role="dialog"
			aria-modal="true"
			aria-label="Navigation menu"
			tabindex="-1"
			class="fixed inset-0 z-40 flex flex-col bg-(--color-bg) px-6 pt-20 pb-8 md:hidden"
			onkeydown={handleKeydown}
		>
			<button
				type="button"
				aria-label="Close navigation menu"
				class="absolute top-4 right-6 grid h-9 w-9 place-items-center rounded-full border border-border"
				onclick={closeMenu}
			>
				<svg
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="1.8"
					class="h-5 w-5"
					aria-hidden="true"
					focusable="false"
				>
					<path d="M6 6 18 18M18 6 6 18" stroke-linecap="round" />
				</svg>
			</button>
			<nav aria-label="Mobile" class="flex flex-col gap-6 text-lg font-medium">
				<!-- svelte-ignore a11y_autofocus -->
				<a href="#product" onclick={closeMenu} autofocus class="hover:text-(--color-text)"
					>Product</a
				>
				<a href="#how-it-works" onclick={closeMenu} class="hover:text-(--color-text)"
					>How it works</a
				>
				<a href="#faq" onclick={closeMenu} class="hover:text-(--color-text)">FAQ</a>
				<a
					href={auth.isAuthenticated ? resolve('/dashboard') : resolve('/login')}
					onclick={() => {
						track('header-cta');
						closeMenu();
					}}
					data-analytics="header-cta"
					class="mt-4 inline-block rounded-full bg-(--color-navy) px-6 py-3 text-center text-sm font-semibold text-white"
				>
					{auth.isAuthenticated ? 'Dashboard' : 'Get started'}</a
				>
			</nav>
		</div>
	{/if}
</header>
