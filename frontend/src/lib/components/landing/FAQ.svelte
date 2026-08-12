<script lang="ts">
	import { reveal } from '$lib/actions/reveal';

	const faqs = [
		{
			q: 'Will elci post something I didn’t approve?',
			a: 'Never. Nothing goes out unless you hit Schedule or Publish. No auto-posting in the background.'
		},
		{
			q: 'Do you store my TikTok or Instagram password?',
			a: 'Nope. You log in on TikTok or Instagram directly. We only get a revocable token, never your password.'
		},
		{
			q: 'Why can I only post privately on TikTok right now?',
			a: 'TikTok keeps new apps in a private-only sandbox until they pass a full audit. Daily caps too. It’s TikTok’s rule, same for any app at this stage.'
		},
		{
			q: 'Can I choose who sees each post?',
			a: 'Yep, only from the privacy options your account actually has right now. We never pre-select one for you.'
		},
		{
			q: 'What happens if TikTok or Instagram rejects a post?',
			a: 'We show you the exact error from the platform and keep your draft intact, so you can fix it and try again. No lost captions.'
		},
		{
			q: 'Do I need to upload my video somewhere else first?',
			a: 'No. Just drop it in elci. We do the chunked upload straight to TikTok or Instagram for you.'
		},
		{
			q: 'Can I schedule for later?',
			a: 'Yes. Pick a date and time and we’ll queue it and publish when it hits.'
		},
		{
			q: 'Which accounts can I connect?',
			a: 'One TikTok and one Instagram account, each connected through that platform’s own login.'
		},
		{
			q: 'Can an AI agent post through elci?',
			a: 'Yes, through the MCP server and REST API. Agent posts go through the same privacy and disclosure checks as a human draft.'
		},
		{
			q: 'Do you have a public API?',
			a: 'Yep, REST with OAuth2, plus native support for n8n, Make and Zapier.'
		}
	];

	let openIndex = $state<number | null>(0);
</script>

<section id="faq" class="mx-auto max-w-3xl px-6 py-20 md:py-28">
	<div use:reveal class="reveal">
		<p class="label mb-3 text-xs text-(--color-accent)">Questions</p>
		<h2 class="text-3xl font-bold tracking-tight md:text-4xl">Good to know</h2>
	</div>

	<div class="mt-10 divide-y divide-(--color-border) rounded-2xl border border-border">
		{#each faqs as item, i (item.q)}
			<div>
				<button
					type="button"
					id="faq-q-{i}"
					aria-controls="faq-a-{i}"
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
						class="h-4 w-4 flex-shrink-0 text-(--color-text-muted) transition-transform"
						style={openIndex === i ? 'transform: rotate(180deg)' : ''}
						aria-hidden="true"
						focusable="false"
					>
						<path d="m6 9 6 6 6-6" stroke-linecap="round" stroke-linejoin="round" />
					</svg>
				</button>
				{#if openIndex === i}
					<div id="faq-a-{i}" role="region" aria-labelledby="faq-q-{i}" class="px-6 pb-5">
						<p class="text-sm leading-relaxed text-(--color-text-muted)">{item.a}</p>
					</div>
				{/if}
			</div>
		{/each}
	</div>
</section>
