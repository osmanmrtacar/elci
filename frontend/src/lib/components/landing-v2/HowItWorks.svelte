<script lang="ts">
	import { reveal } from '$lib/actions/reveal';
	import { auth } from '$lib/auth.svelte';
	import { resolve } from '$app/paths';
	import ComposerMock from './ComposerMock.svelte';

	const steps = [
		{
			title: 'Connect your accounts',
			body: 'Sign in to TikTok, Instagram, and X directly through their own login screens. elci never sees or stores your password.'
		},
		{
			title: 'Write once, adjust per platform',
			body: 'Compose a draft, then fork the caption, media, or settings only where a platform actually needs it. Preview exactly what each one will look like.'
		},
		{
			title: 'Publish or schedule, then track it',
			body: "elci checks the platform's live rules before sending anything, then polls until the post is confirmed up — no guessing."
		}
	];
</script>

<section id="how-it-works" class="mx-auto max-w-6xl px-6 py-20 md:py-24">
	<div class="grid gap-12 md:grid-cols-[1fr_1.2fr] md:gap-16">
		<div class="hidden md:block">
			<div class="sticky top-28 flex justify-center">
				<ComposerMock />
			</div>
		</div>

		<div>
			<div use:reveal class="reveal max-w-xl">
				<span
					class="eyebrow inline-flex items-center gap-2 rounded-full border border-border px-3 py-1 text-xs text-(--color-accent)"
				>
					How it works
				</span>
				<h2 class="mt-3 text-3xl font-semibold tracking-tight md:text-4xl">
					From draft to published in three steps
				</h2>
			</div>

			<ol class="mt-10 divide-y divide-(--color-border) border-y border-border">
				{#each steps as step, i (step.title)}
					<li
						use:reveal
						class="reveal flex items-start gap-6 py-6"
						style="transition-delay: {i * 90}ms"
					>
						<span class="mono pt-1 text-sm font-medium text-(--color-accent)">0{i + 1}</span>
						<div>
							<h3 class="text-lg font-semibold">{step.title}</h3>
							<p class="mt-2 text-sm leading-relaxed text-(--color-text-muted)">{step.body}</p>
						</div>
					</li>
				{/each}
			</ol>

			<a
				href={auth.isAuthenticated ? resolve('/dashboard') : resolve('/login')}
				class="mt-8 inline-block rounded-full bg-(--color-navy) px-6 py-3 text-base font-semibold text-white transition hover:opacity-90"
			>
				{auth.isAuthenticated ? 'Go to dashboard' : 'Get started'}
			</a>
		</div>
	</div>
</section>
