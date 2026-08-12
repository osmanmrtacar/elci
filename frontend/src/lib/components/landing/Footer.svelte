<script lang="ts">
	import { resolve } from '$app/paths';

	const year = new Date().getFullYear();
</script>

<footer class="border-t border-border">
	<div
		class="mx-auto flex max-w-6xl flex-col gap-6 px-6 py-10 md:flex-row md:items-center md:justify-between"
	>
		<div class="flex items-center gap-[0.3em] text-(--color-logo)">
			<span
				class="font-[family-name:var(--font-logo)] text-[16px] leading-none tracking-[-0.005em]"
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
