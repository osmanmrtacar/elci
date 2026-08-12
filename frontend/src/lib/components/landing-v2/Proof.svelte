<script lang="ts">
	import { reveal } from '$lib/actions/reveal';
	import TikTokIcon from '$lib/components/icons/TikTokIcon.svelte';
	import InstagramIcon from '$lib/components/icons/InstagramIcon.svelte';
	import XIcon from '$lib/components/icons/XIcon.svelte';
</script>

<section class="border-y border-border bg-(--color-surface)/60">
	<div use:reveal class="reveal mx-auto max-w-4xl px-6 py-16 text-center">
		<p class="text-sm font-medium text-(--color-text-muted)">
			Built directly against each platform's official publishing API
		</p>
		<div class="mt-6 flex items-center justify-center gap-10 text-(--color-text)">
			<TikTokIcon class="h-6 w-6" />
			<InstagramIcon class="h-6 w-6" />
			<XIcon class="h-6 w-6" />
		</div>
		<p class="mx-auto mt-6 max-w-xl text-sm leading-relaxed text-(--color-text-muted)">
			Every posting flow follows the platform's own published developer guidelines — including
			TikTok's Content Posting API UX requirements, line by line.
		</p>
	</div>
</section>
