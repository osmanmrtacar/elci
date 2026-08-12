<script lang="ts">
	import { reveal } from '$lib/actions/reveal';
	import { auth } from '$lib/auth.svelte';
	import { resolve } from '$app/paths';
	import { track } from '$lib/analytics';
</script>

<section id="agents" class="mx-auto max-w-6xl px-6 py-16 md:py-20">
	<div use:reveal class="reveal mx-auto max-w-2xl text-center">
		<p class="label mb-3 text-xs text-(--color-accent)">For automation too</p>
		<h2 class="text-3xl font-bold tracking-tight md:text-4xl">Agents use the same safe path.</h2>
		<p class="mx-auto mt-4 max-w-lg text-sm leading-relaxed text-(--color-text-muted)">
			Point Claude, ChatGPT, or your own scripts at elci's REST API and MCP server. Every agent post
			runs through the same privacy, disclosure, and upload checks.
		</p>
	</div>
	<div class="mt-8 flex flex-wrap justify-center gap-3 text-sm">
		<span
			class="rounded-full border border-border bg-(--color-surface) px-3 py-1.5 font-mono text-xs"
			>MCP server</span
		>
		<span
			class="rounded-full border border-border bg-(--color-surface) px-3 py-1.5 font-mono text-xs"
			>REST API · OAuth2</span
		>
		<span
			class="rounded-full border border-border bg-(--color-surface) px-3 py-1.5 font-mono text-xs"
			>n8n · Make · Zapier</span
		>
	</div>
	<div use:reveal class="reveal mt-8 text-center">
		<a
			href={auth.isAuthenticated ? resolve('/dashboard') : resolve('/login')}
			data-analytics="agent-teaser-cta"
			onclick={() => track('agent-teaser-cta')}
			class="rounded-full bg-(--color-navy) px-6 py-3 text-sm font-semibold text-white transition hover:opacity-90"
			>Get your API key</a
		>
	</div>
</section>
