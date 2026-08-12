<script lang="ts">
	import { resolve } from '$app/paths';
	import Logo from '$lib/components/Logo.svelte';

	const year = new Date().getFullYear();
</script>

<footer class="border-t border-border">
	<div
		class="mx-auto flex max-w-6xl flex-col gap-6 px-6 py-10 md:flex-row md:items-center md:justify-between"
	>
		<div class="text-(--color-logo)">
			<Logo size={16} />
		</div>
		<p class="text-sm text-(--color-text-muted)">© {year} elci.io. All rights reserved.</p>
		<nav class="flex items-center gap-6 text-sm text-(--color-text-muted)">
			<a href={resolve('/privacy-policy')} class="hover:text-(--color-text)">Privacy</a>
			<a href={resolve('/terms-of-service')} class="hover:text-(--color-text)">Terms</a>
			<a href="mailto:hello@elci.io" class="hover:text-(--color-text)">Contact</a>
		</nav>
	</div>
	<nav
		aria-label="Footer"
		class="mx-auto flex max-w-6xl flex-wrap gap-6 border-t border-border/50 px-6 py-6 text-sm text-(--color-text-muted)"
	>
		<a href="#product" class="hover:text-(--color-text)">Product</a>
		<a href="#how-it-works" class="hover:text-(--color-text)">How it works</a>
		<a href="#agents" class="hover:text-(--color-text)">For automation</a>
		<a href="#faq" class="hover:text-(--color-text)">FAQ</a>
	</nav>
</footer>
