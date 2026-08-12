<script lang="ts">
	import { reveal } from '$lib/actions/reveal';

	const benefits = [
		{
			number: '01',
			title: 'One calendar, every post',
			body: 'See TikTok and Instagram in one weekly view. Drag to reschedule, filter by status, and spot what is next without opening another app.',
			path: 'M3 4h18v13H3zM3 8h18M8 12h8'
		},
		{
			number: '02',
			title: 'One draft, tuned per platform',
			body: 'Write once, then tweak the caption or cover for TikTok vs Instagram without duplicating the whole post.',
			path: 'M12 4 3 8.5 12 13l9-4.5L12 4ZM3 13l9 4.5 9-4.5M3 17.5 12 22l9-4.5'
		},
		{
			number: '03',
			title: 'Real uploads, not just links',
			body: 'Drop in a video and elci handles the chunked upload straight to the platform. No separate file host, no broken URLs.',
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

<section class="mx-auto max-w-6xl px-6 py-20 md:py-28">
	<div use:reveal class="reveal max-w-2xl">
		<p class="label mb-3 text-xs text-(--color-accent)">Why elci</p>
		<h2 class="text-3xl font-bold tracking-tight md:text-4xl">
			One place to run all your posts, not a tab for each app.
		</h2>
	</div>

	<div class="mt-12 grid gap-6 md:grid-cols-2">
		{#each benefits as b, i (b.title)}
			<div
				use:reveal
				class="reveal rounded-2xl border border-border bg-(--color-surface) p-7"
				style="transition-delay: {i * 80}ms"
			>
				<div class="mb-5 flex items-center justify-between">
					<div class="icon-tile grid h-11 w-11 place-items-center rounded-xl">
						<svg
							viewBox="0 0 24 24"
							fill="none"
							stroke="var(--color-accent)"
							stroke-width="1.6"
							class="h-6 w-6"
							stroke-linecap="round"
							stroke-linejoin="round"
							aria-hidden="true"
							focusable="false"
						>
							<path d={b.path} />
						</svg>
					</div>
					<span class="font-mono text-xs font-medium text-(--color-text-muted)">{b.number}</span>
				</div>
				<h3 class="text-lg font-semibold">{b.title}</h3>
				<p class="mt-2 text-sm leading-relaxed text-(--color-text-muted)">{b.body}</p>
			</div>
		{/each}
	</div>
</section>
