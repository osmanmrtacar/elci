<script lang="ts">
	import { reveal } from '$lib/actions/reveal';
	import { auth } from '$lib/auth.svelte';
	import { resolve } from '$app/paths';
</script>

<section class="mx-auto max-w-6xl px-6 py-20 md:py-24">
	<div use:reveal class="reveal mx-auto max-w-2xl text-center">
		<span
			class="eyebrow inline-flex items-center gap-2 rounded-full border border-border px-3 py-1 text-xs text-(--color-accent)"
		>
			Pricing
		</span>
		<h2 class="mt-3 text-3xl font-semibold tracking-tight md:text-4xl">
			Free during early access.
		</h2>
	</div>

	<div
		use:reveal
		class="reveal mx-auto mt-10 max-w-lg rounded-2xl border border-border bg-(--color-surface) p-8 text-center"
	>
		<p class="text-(--color-text-muted)">
			elci is in early access. Connect your accounts and post free while we finish the rest of the
			roadmap — no credit card, no trial countdown.
		</p>
		<p class="mt-4 text-sm text-(--color-text-muted)">
			Early access users get advance notice and a lower rate before any paid plan goes live.
		</p>
		<a
			href={auth.isAuthenticated ? resolve('/dashboard') : resolve('/login')}
			class="mt-6 inline-block rounded-full bg-(--color-accent) px-6 py-3 text-base font-semibold text-white shadow-[0_12px_30px_-12px_var(--color-accent)] transition hover:opacity-90"
		>
			{auth.isAuthenticated ? 'Go to dashboard' : 'Connect your accounts'}
		</a>
		<p class="mt-3 text-xs text-(--color-text-muted)">No credit card required.</p>
	</div>
</section>
