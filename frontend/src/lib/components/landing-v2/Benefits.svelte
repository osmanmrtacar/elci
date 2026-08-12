<script lang="ts">
	import { reveal } from '$lib/actions/reveal';

	const benefits = [
		{
			number: '01',
			title: 'Nothing pre-checked, nothing assumed',
			body: 'Privacy, comments, duet, and stitch are never defaulted. elci greys out whatever your TikTok account already has disabled, before you ever hit publish.',
			path: 'M12 3 4 6.5v5c0 4.7 3.2 8.4 8 9.5 4.8-1.1 8-4.8 8-9.5v-5L12 3Z'
		},
		{
			number: '02',
			title: 'One draft, tuned per platform',
			body: 'Write once, then adjust the caption, media, or settings for TikTok, Instagram, or X individually — without duplicating the whole post.',
			path: 'M12 4 3 8.5 12 13l9-4.5L12 4ZM3 13l9 4.5 9-4.5M3 17.5 12 22l9-4.5'
		},
		{
			number: '03',
			title: 'Real uploads, not just links',
			body: 'Drop in a video and elci handles the chunked upload straight to the platform — no separate file host, no broken URLs.',
			path: 'M12 16V4m0 0-4 4m4-4 4 4M5 16v3a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2v-3'
		},
		{
			number: '04',
			title: 'You always know what happened',
			body: 'elci polls each platform until your post is actually live, and picks up right where it left off if it ever restarts mid-publish.',
			path: 'M3 12h4l2-7 4 14 2-7h6'
		}
	];
</script>

<section class="mx-auto max-w-6xl px-6 py-20 md:py-24">
	<div use:reveal class="reveal max-w-2xl">
		<span
			class="eyebrow inline-flex items-center gap-2 rounded-full border border-border px-3 py-1 text-xs text-(--color-accent)"
		>
			Why elci
		</span>
		<h2 class="mt-3 text-3xl font-semibold tracking-tight md:text-4xl">
			Every platform has its own rulebook. elci reads all of them.
		</h2>
	</div>

	<div class="mt-12 grid gap-6 md:grid-cols-2">
		{#each benefits as b, i (b.title)}
			<div
				use:reveal
				class="reveal relative rounded-2xl border border-border bg-(--color-surface) p-6"
				style="transition-delay: {i * 80}ms"
			>
				<div class="flex items-start justify-between">
					<span class="mono text-sm font-medium text-(--color-text-muted)">{b.number}</span>
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
							<path d={b.path} />
						</svg>
					</div>
				</div>
				<h3 class="mt-6 text-lg font-semibold">{b.title}</h3>
				<p class="mt-2 text-sm leading-relaxed text-(--color-text-muted)">{b.body}</p>
			</div>
		{/each}
	</div>
</section>
