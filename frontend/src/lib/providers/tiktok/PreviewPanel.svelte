<script lang="ts">
	import HeartIcon from '$lib/components/icons/HeartIcon.svelte';
	import CommentIcon from '$lib/components/icons/CommentIcon.svelte';
	import BookmarkIcon from '$lib/components/icons/BookmarkIcon.svelte';
	import ShareIcon from '$lib/components/icons/ShareIcon.svelte';
	import type { PreviewPanelProps } from '../types';

	let { caption, mediaKind, mediaUrl, accountInfo }: PreviewPanelProps = $props();
</script>

<!-- A literal recreation of the TikTok For You feed card — full-bleed video,
caption overlaid at the bottom, action rail on the right — since that's what
"preview" means to a creator, not an abstract settings summary. -->
<div
	class="relative aspect-[9/16] w-full max-w-[280px] overflow-hidden rounded-2xl bg-black shadow-[0_30px_60px_-25px_rgba(16,24,40,0.35)]"
>
	{#if mediaUrl}
		{#if mediaKind === 'video'}
			<video src={mediaUrl} class="h-full w-full object-cover" autoplay muted loop playsinline
			></video>
		{:else}
			<img src={mediaUrl} alt="" class="h-full w-full object-cover" />
		{/if}
	{:else}
		<div
			class="flex h-full w-full items-center justify-center bg-gradient-to-br from-(--color-navy) to-(--color-accent)"
		>
			<p class="font-mono text-[11px] text-white/70">No media yet</p>
		</div>
	{/if}

	<div
		class="pointer-events-none absolute inset-x-0 bottom-0 flex flex-col gap-1.5 bg-gradient-to-t from-black/85 via-black/25 to-transparent p-4 pt-14 pr-14"
	>
		<p class="text-sm font-semibold text-white">@{accountInfo?.username || 'username'}</p>
		<p class="line-clamp-4 text-sm text-white/90">{caption || 'Your caption will appear here.'}</p>
	</div>

	<div
		class="pointer-events-none absolute inset-y-0 right-2.5 flex flex-col items-center justify-end gap-4 pb-6"
	>
		<span class="h-8 w-8 overflow-hidden rounded-full border-2 border-white bg-white/20">
			{#if accountInfo?.avatarUrl}
				<img src={accountInfo.avatarUrl} alt="" class="h-full w-full object-cover" />
			{/if}
		</span>
		<HeartIcon class="h-6 w-6 text-white drop-shadow" />
		<CommentIcon class="h-6 w-6 text-white drop-shadow" />
		<BookmarkIcon class="h-6 w-6 text-white drop-shadow" />
		<ShareIcon class="h-6 w-6 text-white drop-shadow" />
	</div>
</div>
