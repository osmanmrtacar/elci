<script lang="ts">
	import type { SettingsPanelProps } from '../types';

	let { caption, mediaKind, infoError }: SettingsPanelProps = $props();

	const maxCaptionLen = 280;
</script>

<div class="space-y-3 text-sm">
	{#if infoError}
		<p class="rounded-lg bg-(--color-status-failed)/10 p-3 text-(--color-status-failed)">
			Couldn't load this X account: {infoError}
		</p>
	{/if}

	<p
		class="font-mono text-xs {caption.length > maxCaptionLen ? 'text-(--color-status-failed)' : ''}"
	>
		{caption.length} / {maxCaptionLen}
	</p>

	<p class="text-(--color-text-muted)">
		{#if mediaKind === 'video'}
			One video is allowed per post.
		{:else}
			Up to 4 images are allowed per post — this composer currently sends one.
		{/if}
	</p>
</div>
