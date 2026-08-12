<script lang="ts">
	import { reveal } from '$lib/actions/reveal';

	const faqs = [
		{
			q: 'Does elci ever post without my review?',
			a: 'No. elci only publishes what you explicitly schedule or hit publish on — nothing goes out on its own.'
		},
		{
			q: 'Do you store my TikTok, Instagram, or X password?',
			a: "No. Every account connects through that platform's own login screen. elci only ever receives a revocable access token, never your password."
		},
		{
			q: 'Why can I only post privately on TikTok?',
			a: "Until an app completes TikTok's audit, their API restricts every post to private viewing and caps how many accounts can post per day — that's TikTok's rule, applied the same way to any app in this position."
		},
		{
			q: 'Can I choose who sees each post?',
			a: 'Yes — from exactly the privacy options your account allows at that moment. elci never pre-selects one for you.'
		},
		{
			q: 'What happens if a platform rejects a post?',
			a: "elci shows you the platform's own error message and keeps your draft intact, so you can fix it and try again instead of losing your caption and media."
		},
		{
			q: 'Do I need to upload my video somewhere else first?',
			a: 'No — upload it directly in elci. It handles the chunked upload to the platform for you.'
		},
		{
			q: 'Can I schedule a post for later?',
			a: 'Yes. Pick a date and time and elci queues it, then publishes automatically.'
		},
		{
			q: 'Which accounts can I connect?',
			a: 'One TikTok, one Instagram, and one X account, each authorized directly with that platform.'
		},
		{
			q: 'Can an AI agent post through elci?',
			a: 'Yes. elci exposes an MCP server and a REST API, so agents like Claude or ChatGPT — or your own scripts — can draft, schedule, and publish through the same rule checks a person using the dashboard would go through.'
		},
		{
			q: 'Do you have a public API?',
			a: 'Yes — a REST API with OAuth2 for building your own integrations, plus native support for tools like n8n, Make, and Zapier.'
		}
	];

	let openIndex = $state<number | null>(0);
</script>

<section id="faq" class="mx-auto max-w-3xl px-6 py-20 md:py-24">
	<div use:reveal class="reveal">
		<span
			class="eyebrow inline-flex items-center gap-2 rounded-full border border-border px-3 py-1 text-xs text-(--color-accent)"
		>
			Questions
		</span>
		<h2 class="mt-3 text-3xl font-semibold tracking-tight md:text-4xl">Good to know</h2>
	</div>

	<div class="mt-10 divide-y divide-(--color-border) rounded-2xl border border-border">
		{#each faqs as item, i (item.q)}
			<div>
				<button
					type="button"
					class="flex w-full items-center justify-between gap-4 px-6 py-5 text-left text-base font-semibold"
					onclick={() => (openIndex = openIndex === i ? null : i)}
					aria-expanded={openIndex === i}
				>
					{item.q}
					<svg
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						class="h-4 w-4 flex-shrink-0 text-(--color-text-muted) transition-transform duration-500 ease-[cubic-bezier(0.32,0.72,0,1)]"
						style={openIndex === i ? 'transform: rotate(180deg)' : ''}
					>
						<path d="m6 9 6 6 6-6" stroke-linecap="round" stroke-linejoin="round" />
					</svg>
				</button>
				{#if openIndex === i}
					<p class="px-6 pb-5 text-sm leading-relaxed text-(--color-text-muted)">{item.a}</p>
				{/if}
			</div>
		{/each}
	</div>
</section>
