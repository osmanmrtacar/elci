<script lang="ts">
	import TikTokIcon from '$lib/components/icons/TikTokIcon.svelte';
	import HeartIcon from '$lib/components/icons/HeartIcon.svelte';
	import CommentIcon from '$lib/components/icons/CommentIcon.svelte';
	import BookmarkIcon from '$lib/components/icons/BookmarkIcon.svelte';
	import ShareIcon from '$lib/components/icons/ShareIcon.svelte';
</script>

<!-- Mirrors the real TikTok preview card in the composer exactly (full-bleed
media, caption overlaid at the bottom, action rail on the right) — the hero
should never show a mockup the product doesn't actually produce. -->
<div
	class="relative aspect-[9/16] w-full max-w-[280px] overflow-hidden rounded-2xl bg-black shadow-[0_30px_60px_-25px_rgba(16,24,40,0.35)]"
>
	<div
		class="flex h-full w-full items-center justify-center bg-gradient-to-br from-(--color-navy) to-(--color-accent)"
	>
		<svg viewBox="0 0 24 24" fill="white" class="h-12 w-12 opacity-90">
			<path d="M8 5.5v13l11-6.5-11-6.5Z" />
		</svg>
	</div>

	<span
		class="mono absolute top-3 left-3 rounded-full px-2.5 py-1 text-[11px] font-medium tracking-wide text-(--color-status-scheduled)"
		style="background-color: color-mix(in srgb, var(--color-status-scheduled) 20%, black)"
	>
		Scheduled · Fri 9:00 AM
	</span>

	<div
		class="pointer-events-none absolute inset-x-0 bottom-0 flex flex-col gap-1.5 bg-gradient-to-t from-black/85 via-black/25 to-transparent p-4 pt-12 pr-12"
	>
		<p class="text-sm font-semibold text-white">@sunnyridgecoffee</p>
		<p class="line-clamp-4 text-sm text-white/90">
			Our oat cortado is back for fall 🍂 first 50 in the shop get one free.
		</p>
	</div>

	<div
		class="pointer-events-none absolute inset-y-0 right-2.5 flex flex-col items-center justify-end gap-4 pb-6"
	>
		<span
			class="grid h-8 w-8 place-items-center overflow-hidden rounded-full border-2 border-white bg-(--color-navy) text-white"
		>
			<TikTokIcon class="h-4 w-4" />
		</span>
		<HeartIcon class="h-6 w-6 text-white drop-shadow" />
		<CommentIcon class="h-6 w-6 text-white drop-shadow" />
		<BookmarkIcon class="h-6 w-6 text-white drop-shadow" />
		<ShareIcon class="h-6 w-6 text-white drop-shadow" />
	</div>
</div>
