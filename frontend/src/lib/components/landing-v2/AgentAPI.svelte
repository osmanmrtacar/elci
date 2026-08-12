<script lang="ts">
	import { reveal } from '$lib/actions/reveal';
	import { auth } from '$lib/auth.svelte';
	import { resolve } from '$app/paths';

	const cards = [
		{
			title: 'MCP server',
			body: 'Give Claude, ChatGPT, or any MCP-compatible agent direct tools to draft, schedule, and publish — without leaving the conversation.',
			path: 'M12 3v4M12 17v4M4.9 7.8l3.3 2M15.8 14.2l3.3 2M4.9 16.2l3.3-2M15.8 9.8l3.3-2M9 12a3 3 0 1 0 6 0 3 3 0 0 0-6 0Z'
		},
		{
			title: 'REST API with OAuth2',
			body: 'Build your own posting workflows, or wire elci into n8n, Make, or Zapier using the same publishing endpoints the dashboard runs on.',
			path: 'M9 7V4m6 3V4M9 20v-3m6 3v-3M6 7h12a2 2 0 0 1 2 2v6a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V9a2 2 0 0 1 2-2Z'
		},
		{
			title: 'Same rules, enforced automatically',
			body: "Agent-issued posts pass through elci's platform-rule checks too — privacy defaults, disclosure requirements, upload validation — so automation can't get an account flagged either.",
			path: 'M12 3 4 6.5v5c0 4.7 3.2 8.4 8 9.5 4.8-1.1 8-4.8 8-9.5v-5L12 3ZM9 12.3l2 2 4-4.4'
		}
	];
</script>

<section id="agents" class="mx-auto max-w-6xl px-6 py-20 md:py-24">
	<div use:reveal class="reveal mx-auto max-w-2xl text-center">
		<span
			class="eyebrow inline-flex items-center gap-2 rounded-full border border-border px-3 py-1 text-xs text-(--color-accent)"
		>
			For developers &amp; agents
		</span>
		<h2 class="mt-3 text-3xl font-semibold tracking-tight md:text-4xl">
			Built to be driven by agents too.
		</h2>
		<p class="mx-auto mt-4 max-w-lg text-(--color-text-muted)">
			Point Claude, ChatGPT, or your own automation at elci's API and MCP server. Every post they
			create still runs through the same platform checks a human would have to pass.
		</p>
	</div>

	<div class="mt-12 grid gap-6 md:grid-cols-3">
		{#each cards as card, i (card.title)}
			<div
				use:reveal
				class="reveal rounded-2xl border border-border bg-(--color-surface) p-6"
				style="transition-delay: {i * 80}ms"
			>
				<div class="icon-tile grid h-11 w-11 place-items-center rounded-xl">
					<svg
						viewBox="0 0 24 24"
						fill="none"
						stroke="var(--color-accent)"
						stroke-width="1.6"
						class="h-6 w-6"
						stroke-linecap="round"
						stroke-linejoin="round"
					>
						<path d={card.path} />
					</svg>
				</div>
				<h3 class="mt-6 text-lg font-semibold">{card.title}</h3>
				<p class="mt-2 text-sm leading-relaxed text-(--color-text-muted)">{card.body}</p>
			</div>
		{/each}
	</div>

	<div use:reveal class="reveal mt-10 text-center">
		<a
			href={auth.isAuthenticated ? resolve('/dashboard') : resolve('/login')}
			class="rounded-full bg-(--color-navy) px-6 py-3 text-base font-semibold text-white transition hover:opacity-90"
		>
			{auth.isAuthenticated ? 'Go to dashboard' : 'Get your API key'}
		</a>
	</div>
</section>
