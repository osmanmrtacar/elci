<script lang="ts">
	import { reveal } from '$lib/actions/reveal';
	import { auth } from '$lib/auth.svelte';
	import { resolve } from '$app/paths';
	import { track } from '$lib/analytics';
</script>

<section class="mx-auto max-w-4xl px-6 pt-8 pb-24 text-center md:pb-32">
	<div use:reveal class="reveal">
		<h2 class="text-3xl font-bold tracking-tight md:text-4xl">
			One calendar for everything you ship.
		</h2>
		<p class="mx-auto mt-4 max-w-lg text-(--color-text-muted)">
			Connect TikTok and Instagram today. Your next platforms will live in the same place.
		</p>
		<a
			href={auth.isAuthenticated ? resolve('/dashboard') : resolve('/login')}
			data-analytics="final-cta"
			onclick={() => track('final-cta')}
			class="mt-8 inline-block rounded-full bg-(--color-accent) px-7 py-3.5 text-base font-semibold text-white shadow-[0_12px_30px_-12px_var(--color-accent)] transition hover:opacity-90"
		>
			{auth.isAuthenticated ? 'Go to dashboard' : 'Connect your accounts'}
		</a>
		<p class="mt-3 text-xs text-(--color-text-muted)">Free during early access · No credit card</p>
		<a
			href="#faq"
			class="mt-6 inline-block text-sm font-semibold text-(--color-text-muted) hover:text-(--color-text)"
			>Read the FAQ <span aria-hidden="true">→</span></a
		>
	</div>
</section>
