<script lang="ts">
	import XIcon from '$lib/components/icons/XIcon.svelte';
	import CommentIcon from '$lib/components/icons/CommentIcon.svelte';
	import RetweetIcon from '$lib/components/icons/RetweetIcon.svelte';
	import HeartIcon from '$lib/components/icons/HeartIcon.svelte';
	import ShareIcon from '$lib/components/icons/ShareIcon.svelte';
	import type { PreviewPanelProps } from '../types';

	let { caption, mediaKind, mediaUrls, accountInfo }: PreviewPanelProps = $props();
</script>

<!-- A literal recreation of a post in the X timeline — avatar, name/handle,
text, media, then the reply/repost/like/share row. -->
<div
	class="w-full max-w-[280px] rounded-2xl border border-border bg-(--color-surface) p-4 shadow-[0_30px_60px_-25px_rgba(16,24,40,0.35)]"
>
	<div class="flex gap-2.5">
		<span
			class="grid h-9 w-9 shrink-0 place-items-center overflow-hidden rounded-full bg-(--color-navy) text-white"
		>
			{#if accountInfo?.avatarUrl}
				<img src={accountInfo.avatarUrl} alt="" class="h-full w-full object-cover" />
			{:else}
				<XIcon class="h-3.5 w-3.5" />
			{/if}
		</span>
		<div class="min-w-0 flex-1">
			<div class="flex flex-wrap items-baseline gap-1 text-sm">
				<span class="font-semibold">{accountInfo?.displayName || 'Display name'}</span>
				<span class="text-(--color-text-muted)">@{accountInfo?.username || 'handle'}</span>
			</div>

			<p class="mt-0.5 text-sm leading-snug break-words">
				{caption || 'Your post will appear here.'}
			</p>

			{#if mediaUrls.length > 0}
				<div class="mt-2 overflow-hidden rounded-xl border border-border">
					{#if mediaKind === 'video'}
						<video
							src={mediaUrls[0]}
							class="max-h-56 w-full object-cover"
							autoplay
							muted
							loop
							playsinline
						></video>
					{:else if mediaUrls.length === 1}
						<img src={mediaUrls[0]} alt="" class="max-h-56 w-full object-cover" />
					{:else if mediaUrls.length === 2}
						<div class="grid h-44 grid-cols-2 gap-0.5">
							{#each mediaUrls as url (url)}
								<img src={url} alt="" class="h-full w-full object-cover" />
							{/each}
						</div>
					{:else if mediaUrls.length === 3}
						<div class="grid h-44 grid-cols-2 grid-rows-2 gap-0.5">
							<img src={mediaUrls[0]} alt="" class="row-span-2 h-full w-full object-cover" />
							<img src={mediaUrls[1]} alt="" class="h-full w-full object-cover" />
							<img src={mediaUrls[2]} alt="" class="h-full w-full object-cover" />
						</div>
					{:else}
						<div class="grid h-44 grid-cols-2 grid-rows-2 gap-0.5">
							{#each mediaUrls.slice(0, 4) as url (url)}
								<img src={url} alt="" class="h-full w-full object-cover" />
							{/each}
						</div>
					{/if}
				</div>
			{/if}

			<div class="mt-3 flex items-center justify-between text-(--color-text-muted)">
				<CommentIcon class="h-[18px] w-[18px]" />
				<RetweetIcon class="h-[18px] w-[18px]" />
				<HeartIcon class="h-[18px] w-[18px]" />
				<ShareIcon class="h-[18px] w-[18px]" />
			</div>
		</div>
	</div>
</div>
