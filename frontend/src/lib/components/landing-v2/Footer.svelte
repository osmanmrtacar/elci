<script lang="ts">
	import { resolve } from '$app/paths';

	const year = new Date().getFullYear();
</script>

<footer class="border-t border-border">
	<div
		class="mx-auto flex max-w-6xl flex-col gap-6 px-6 py-10 md:flex-row md:items-center md:justify-between"
	>
		<div class="flex items-center gap-2 font-semibold">
			elci
			<span class="inline-block h-1.5 w-1.5 rounded-full bg-(--color-accent)"></span>
		</div>
		<p class="text-sm text-(--color-text-muted)">© {year} elci.io. All rights reserved.</p>
		<nav class="flex items-center gap-6 text-sm text-(--color-text-muted)">
			<a href={resolve('/privacy-policy')} class="hover:text-(--color-text)">Privacy</a>
			<a href={resolve('/terms-of-service')} class="hover:text-(--color-text)">Terms</a>
			<a href="mailto:hello@elci.io" class="hover:text-(--color-text)">Contact</a>
		</nav>
	</div>
</footer>
