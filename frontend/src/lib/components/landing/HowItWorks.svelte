<script lang="ts">
	import { reveal } from '$lib/actions/reveal';
	import { auth } from '$lib/auth.svelte';
	import { resolve } from '$app/paths';
	import { track } from '$lib/analytics';

	const steps = [
		{
			title: 'Connect your accounts',
			body: 'Sign in to TikTok and Instagram through their own login screens. elci never sees your password, just a revocable token.'
		},
		{
			title: 'Write once, adjust per platform',
			body: 'Draft once, then tweak the caption, cover or settings only where TikTok or Instagram actually needs it. Preview exactly what each will look like.'
		},
		{
			title: 'Publish or schedule, then track it',
			body: "elci checks the platform's live rules before sending anything, then polls until the post is confirmed up. No guessing."
		}
	];
</script>

<section id="how-it-works" class="mx-auto max-w-6xl px-6 py-20 md:py-28">
	<div use:reveal class="reveal max-w-2xl">
		<p class="label mb-3 text-xs text-(--color-accent)">How it works</p>
		<h2 class="text-3xl font-bold tracking-tight md:text-4xl">
			From draft to published in three steps
		</h2>
	</div>

	<ol class="relative mt-14 grid gap-10 md:grid-cols-3 md:gap-8">
		<div
			class="pointer-events-none absolute top-6 right-[calc(16px+1rem)] left-[calc(16px+1rem)] hidden h-px bg-(--color-border) md:block"
			aria-hidden="true"
		></div>
		{#each steps as step, i (step.title)}
			<li use:reveal class="reveal relative" style="transition-delay: {i * 90}ms">
				<span class="font-mono text-sm font-medium text-(--color-accent)">0{i + 1}</span>
				<h3 class="mt-3 text-lg font-semibold">{step.title}</h3>
				<p class="mt-2 text-sm leading-relaxed text-(--color-text-muted)">{step.body}</p>
			</li>
		{/each}
	</ol>
	<a
		href={auth.isAuthenticated ? resolve('/dashboard') : resolve('/login')}
		data-analytics="howitworks-cta"
		onclick={() => track('howitworks-cta')}
		class="mt-10 inline-block rounded-full bg-(--color-navy) px-6 py-3 text-sm font-semibold text-white transition hover:opacity-90"
	>
		Get started. No credit card
	</a>
</section>
