<script lang="ts">
	import ComposerMock from './ComposerMock.svelte';
	import { auth } from '$lib/auth.svelte';
	import { resolve } from '$app/paths';
	import TikTokIcon from '$lib/components/icons/TikTokIcon.svelte';
	import InstagramIcon from '$lib/components/icons/InstagramIcon.svelte';
	import XIcon from '$lib/components/icons/XIcon.svelte';

	const facts = [
		'Nothing pre-checked or assumed',
		'Real chunked uploads, no link hosts',
		'Polls until the post is confirmed live',
		'Scriptable through an API and MCP server'
	];
</script>

<section id="product" class="mx-auto max-w-6xl px-6 pt-24 pb-20 md:pt-32 md:pb-24">
	<div class="grid items-center gap-12 md:grid-cols-[1.1fr_1fr]">
		<div>
			<span
				class="eyebrow inline-flex items-center gap-2 rounded-full border border-border px-3 py-1 text-xs text-(--color-accent)"
			>
				For creators &amp; small teams
			</span>
			<h1 class="gradient-heading mt-5 max-w-xl text-4xl font-semibold md:text-5xl">
				Post to TikTok, Instagram, and X without breaking their rules.
			</h1>
			<p class="mt-5 max-w-lg text-lg text-(--color-text-muted)">
				elci is a scheduler built directly against each platform's own publishing API — down to
				TikTok's exact privacy, disclosure, and interaction requirements — so a UX mistake never
				gets your account flagged again.
			</p>
			<div class="mt-8 flex flex-wrap items-center gap-4">
				<a
					href={auth.isAuthenticated ? resolve('/dashboard') : resolve('/login')}
					class="rounded-full bg-(--color-accent) px-6 py-3 text-base font-semibold text-white shadow-[0_12px_30px_-12px_var(--color-accent)] transition hover:opacity-90"
				>
					{auth.isAuthenticated ? 'Go to dashboard' : 'Connect your accounts'}
				</a>
				<a
					href="#how-it-works"
					class="text-sm font-semibold text-(--color-text-muted) hover:text-(--color-text)"
				>
					See how it works →
				</a>
			</div>

			<div class="mt-10 border-t border-border pt-6">
				<div class="mb-3 flex items-center gap-3 text-(--color-text-muted)">
					<TikTokIcon class="h-4 w-4" />
					<InstagramIcon class="h-4 w-4" />
					<XIcon class="h-4 w-4" />
					<span class="mono text-xs">Official publishing APIs, not screen automation</span>
				</div>
				<ul
					class="grid grid-cols-1 gap-x-6 gap-y-1.5 text-sm text-(--color-text-muted) sm:grid-cols-2"
				>
					{#each facts as fact (fact)}
						<li class="flex items-baseline gap-2">
							<span class="text-(--color-accent)">·</span>
							{fact}
						</li>
					{/each}
				</ul>
			</div>
		</div>

		<div class="flex justify-center md:justify-end">
			<ComposerMock />
		</div>
	</div>
</section>
