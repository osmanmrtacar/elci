<script lang="ts">
	import { reveal } from '$lib/actions/reveal';
	import { auth } from '$lib/auth.svelte';
	import { resolve } from '$app/paths';
</script>

<section class="mx-auto max-w-4xl px-6 pt-8 pb-24 text-center">
	<div use:reveal class="reveal">
		<h2 class="text-3xl font-semibold tracking-tight md:text-4xl">
			Stop guessing what each platform will allow.
		</h2>
		<p class="mx-auto mt-4 max-w-lg text-(--color-text-muted)">
			Connect TikTok, Instagram, and X, and let elci handle the rules.
		</p>
		<a
			href={auth.isAuthenticated ? resolve('/dashboard') : resolve('/login')}
			class="mt-8 inline-block rounded-full bg-(--color-accent) px-6 py-3 text-base font-semibold text-white shadow-[0_12px_30px_-12px_var(--color-accent)] transition hover:opacity-90"
		>
			{auth.isAuthenticated ? 'Go to dashboard' : 'Connect your accounts'}
		</a>
	</div>
</section>
