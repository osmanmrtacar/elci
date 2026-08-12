<script lang="ts">
	import type { SettingsPanelProps } from '../types';

	let { caption, mediaKind, infoError }: SettingsPanelProps = $props();

	const maxCaptionLen = 2200;
</script>

<div class="space-y-3 text-sm">
	{#if infoError}
		<p class="rounded-lg bg-(--color-status-failed)/10 p-3 text-(--color-status-failed)">
			Couldn't load this Instagram account: {infoError}
		</p>
	{/if}

	{#if caption.length > maxCaptionLen}
		<p class="text-(--color-status-failed)">
			Caption is {caption.length - maxCaptionLen} characters over Instagram's {maxCaptionLen} limit.
		</p>
	{/if}

	<p class="text-(--color-text-muted)">
		{#if mediaKind === 'video'}
			This will publish as a Reel.
		{:else}
			Instagram has no extra per-post settings for a single image — it publishes with the caption
			above.
		{/if}
	</p>
</div>
