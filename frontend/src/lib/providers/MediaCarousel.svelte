<script lang="ts">
	import type { MediaKind } from '$lib/types';

	let { mediaUrls, mediaKind }: { mediaUrls: string[]; mediaKind: MediaKind | null } = $props();

	let index = $state(0);
	$effect(() => {
		if (index >= mediaUrls.length) index = 0;
	});

	function prev(e: MouseEvent) {
		e.stopPropagation();
		index = (index - 1 + mediaUrls.length) % mediaUrls.length;
	}
	function next(e: MouseEvent) {
		e.stopPropagation();
		index = (index + 1) % mediaUrls.length;
	}
</script>

{#if mediaUrls.length > 0}
	{#if mediaKind === 'video'}
		<video src={mediaUrls[0]} class="h-full w-full object-cover" autoplay muted loop playsinline
		></video>
	{:else}
		<img src={mediaUrls[index]} alt="" class="h-full w-full object-cover" />
	{/if}

	{#if mediaUrls.length > 1}
		<button
			type="button"
			onclick={prev}
			aria-label="Previous image"
			class="absolute top-1/2 left-1.5 grid h-6 w-6 -translate-y-1/2 place-items-center rounded-full bg-black/50 text-white"
		>
			<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="h-3.5 w-3.5">
				<path d="M15 6l-6 6 6 6" stroke-linecap="round" stroke-linejoin="round" />
			</svg>
		</button>
		<button
			type="button"
			onclick={next}
			aria-label="Next image"
			class="absolute top-1/2 right-1.5 grid h-6 w-6 -translate-y-1/2 place-items-center rounded-full bg-black/50 text-white"
		>
			<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="h-3.5 w-3.5">
				<path d="M9 6l6 6-6 6" stroke-linecap="round" stroke-linejoin="round" />
			</svg>
		</button>
		<span
			class="absolute top-2 right-2 rounded-full bg-black/50 px-1.5 py-0.5 font-mono text-[10px] text-white"
		>
			{index + 1}/{mediaUrls.length}
		</span>
	{/if}
{/if}
