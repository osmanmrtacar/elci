<script lang="ts">
	import InstagramIcon from '$lib/components/icons/InstagramIcon.svelte';
	import HeartIcon from '$lib/components/icons/HeartIcon.svelte';
	import CommentIcon from '$lib/components/icons/CommentIcon.svelte';
	import ShareIcon from '$lib/components/icons/ShareIcon.svelte';
	import BookmarkIcon from '$lib/components/icons/BookmarkIcon.svelte';
	import MediaCarousel from '../MediaCarousel.svelte';
	import type { PreviewPanelProps } from '../types';

	let { caption, mediaKind, mediaUrls, accountInfo }: PreviewPanelProps = $props();
</script>

<!-- A literal recreation of an Instagram feed post — header, square media,
action row, then caption — rather than an abstract settings summary. -->
<div
	class="w-full max-w-[280px] overflow-hidden rounded-2xl border border-border bg-(--color-surface) shadow-[0_30px_60px_-25px_rgba(16,24,40,0.35)]"
>
	<div class="flex items-center gap-2 p-3">
		<span
			class="grid h-7 w-7 shrink-0 place-items-center overflow-hidden rounded-full bg-(--color-navy) text-white"
		>
			{#if accountInfo?.avatarUrl}
				<img src={accountInfo.avatarUrl} alt="" class="h-full w-full object-cover" />
			{:else}
				<InstagramIcon class="h-3.5 w-3.5" />
			{/if}
		</span>
		<p class="text-sm font-semibold">{accountInfo?.username || 'username'}</p>
	</div>

	<div class="relative aspect-square w-full bg-(--color-border)">
		{#if mediaUrls.length > 0}
			<MediaCarousel {mediaUrls} {mediaKind} />
		{:else}
			<div class="flex h-full items-center justify-center">
				<p class="font-mono text-[11px] text-(--color-text-muted)">No media yet</p>
			</div>
		{/if}
	</div>

	<div class="flex items-center gap-3 px-3 pt-2.5 pb-1">
		<HeartIcon class="h-6 w-6" />
		<CommentIcon class="h-6 w-6" />
		<ShareIcon class="h-6 w-6" />
		<BookmarkIcon class="ml-auto h-6 w-6" />
	</div>

	<p class="px-3 pb-3 text-sm leading-snug break-words">
		<span class="font-semibold">{accountInfo?.username || 'username'}</span>
		{caption || 'Your caption will appear here.'}
	</p>
</div>
