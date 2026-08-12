<script lang="ts">
	import TikTokIcon from '$lib/components/icons/TikTokIcon.svelte';
	import HeartIcon from '$lib/components/icons/HeartIcon.svelte';
	import CommentIcon from '$lib/components/icons/CommentIcon.svelte';
	import BookmarkIcon from '$lib/components/icons/BookmarkIcon.svelte';
	import ShareIcon from '$lib/components/icons/ShareIcon.svelte';
</script>

<!-- Mirrors the real TikTok preview card in the composer exactly (full-bleed
media, caption overlaid at the bottom, action rail on the right). The hero
should never show a mockup the product does not actually produce. -->
<div
	role="img"
	aria-label="Preview of a scheduled TikTok post in the elci composer"
	class="relative aspect-[9/16] w-full max-w-[280px] overflow-hidden rounded-2xl bg-black shadow-[0_30px_60px_-25px_rgba(16,24,40,0.35)]"
>
	<img
		src="/mock-video.jpg"
		alt=""
		width="540"
		height="960"
		loading="eager"
		decoding="async"
		class="absolute inset-0 h-full w-full object-cover"
	/>

	<div class="pointer-events-none absolute inset-0 bg-gradient-to-t from-black/30 via-transparent to-transparent"></div>

	<div class="pointer-events-none absolute inset-0 flex items-center justify-center">
		<span class="grid h-14 w-14 place-items-center rounded-full bg-white/90 backdrop-blur">
			<svg
				viewBox="0 0 24 24"
				fill="currentColor"
				class="ml-0.5 h-6 w-6 text-(--color-navy)"
				aria-hidden="true"
				focusable="false"
			>
				<path d="M8 5.5v13l11-6.5-11-6.5Z" />
			</svg>
		</span>
	</div>

	<span
		class="absolute top-3 left-3 rounded-full px-2.5 py-1 font-mono text-[11px] font-medium tracking-wide text-(--color-status-scheduled)"
		style="background-color: color-mix(in srgb, var(--color-status-scheduled) 20%, black)"
	>
		Scheduled · Fri 9:00 AM
	</span>

	<div
		class="pointer-events-none absolute inset-x-0 bottom-0 flex flex-col gap-1.5 bg-gradient-to-t from-black/85 via-black/25 to-transparent p-4 pt-14 pr-14"
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
			aria-hidden="true"
		>
			<span aria-hidden="true"><TikTokIcon class="h-4 w-4" /></span>
		</span>
		<span aria-hidden="true"><HeartIcon class="h-6 w-6 text-white drop-shadow" /></span>
		<span aria-hidden="true"><CommentIcon class="h-6 w-6 text-white drop-shadow" /></span>
		<span aria-hidden="true"><BookmarkIcon class="h-6 w-6 text-white drop-shadow" /></span>
		<span aria-hidden="true"><ShareIcon class="h-6 w-6 text-white drop-shadow" /></span>
	</div>
</div>
