<script lang="ts">
	import ComposerMock from './ComposerMock.svelte';
	import { auth } from '$lib/auth.svelte';
	import { resolve } from '$app/paths';
	import { track } from '$lib/analytics';
</script>

<section id="product" class="mx-auto max-w-6xl px-6 pt-16 pb-20 md:pt-24 md:pb-28">
	<div class="grid items-center gap-14 md:grid-cols-[1.1fr_1fr]">
		<div>
			<p class="label mb-5 text-xs text-(--color-accent)">One calendar for all your posts</p>
			<h1 class="max-w-xl text-4xl font-bold tracking-tight md:text-5xl">
				One place to ship. Not to juggle.
			</h1>
			<p class="mt-5 max-w-lg text-lg text-(--color-text-muted)">
				Plan this week in one calendar, not three apps. Write once in elci and publish to
				TikTok and Instagram today, with more platforms joining the same flow. Check, schedule
				and track from one screen.
			</p>
			<div class="mt-8 flex flex-wrap items-center gap-4">
				<a
					href={auth.isAuthenticated ? resolve('/dashboard') : resolve('/login')}
					data-analytics="hero-primary"
					onclick={() => track('hero-primary')}
					class="rounded-full bg-(--color-accent-2) px-6 py-3 text-base font-semibold text-(--color-navy) shadow-[0_12px_30px_-12px_var(--color-accent-2)] transition hover:opacity-90"
				>
					{auth.isAuthenticated ? 'Go to dashboard' : 'Connect your accounts'}
				</a>
				<a
					href="#how-it-works"
					class="text-sm font-semibold text-(--color-text-muted) hover:text-(--color-text)"
				>
					See how it works <span aria-hidden="true">→</span>
				</a>
			</div>
			<p class="mt-8 text-xs text-(--color-text-muted)">
				Live with TikTok and Instagram. More platforms on the way.
			</p>
			<div class="mt-3 flex flex-wrap gap-3 text-xs text-(--color-text-muted)">
				<span class="inline-flex items-center gap-1.5">
					<svg
						viewBox="0 0 16 16"
						class="h-3.5 w-3.5"
						fill="none"
						stroke="currentColor"
						stroke-width="1.6"
						aria-hidden="true"
						focusable="false"
					>
						<path d="M3 8.5 6 11.5 13 4.5" stroke-linecap="round" stroke-linejoin="round" />
					</svg>
					One draft to manage
				</span>
				<span class="inline-flex items-center gap-1.5">
					<svg
						viewBox="0 0 16 16"
						class="h-3.5 w-3.5"
						fill="none"
						stroke="currentColor"
						stroke-width="1.6"
						aria-hidden="true"
						focusable="false"
					>
						<path d="M3 8.5 6 11.5 13 4.5" stroke-linecap="round" stroke-linejoin="round" />
					</svg>
					One calendar to plan
				</span>
				<span class="inline-flex items-center gap-1.5">
					<svg
						viewBox="0 0 16 16"
						class="h-3.5 w-3.5"
						fill="none"
						stroke="currentColor"
						stroke-width="1.6"
						aria-hidden="true"
						focusable="false"
					>
						<path d="M3 8.5 6 11.5 13 4.5" stroke-linecap="round" stroke-linejoin="round" />
					</svg>
					No app hopping
				</span>
			</div>
			<p class="mt-6 text-xs font-medium text-(--color-text-muted)">
				Create, schedule and see what shipped, all in one place
			</p>
		</div>

		<div class="flex justify-center md:justify-end">
			<ComposerMock />
		</div>
	</div>
</section>
