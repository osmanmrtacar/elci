<script lang="ts">
	import { reveal } from '$lib/actions/reveal';
	import TikTokIcon from '$lib/components/icons/TikTokIcon.svelte';
	import InstagramIcon from '$lib/components/icons/InstagramIcon.svelte';
</script>

<section class="border-y border-border bg-(--color-surface)/60">
	<div use:reveal class="reveal mx-auto max-w-4xl px-6 py-16 text-center">
		<p class="text-sm font-medium text-(--color-text-muted)">
			Built on the official TikTok and Instagram APIs, not a workaround
		</p>
		<div class="mt-6 flex items-center justify-center gap-10 text-(--color-text)">
			<TikTokIcon class="h-6 w-6" />
			<InstagramIcon class="h-6 w-6" />
		</div>
		<p class="mx-auto mt-6 max-w-xl text-sm leading-relaxed text-(--color-text-muted)">
			We follow TikTok’s Content Posting API checklist line by line, privacy, disclosures,
			comment/duet/stitch so your draft can’t slip into the wrong setting.
		</p>
		<div class="mt-8 grid gap-4 text-left md:grid-cols-3">
			<div
				class="flex items-start gap-2 rounded-xl border border-border bg-(--color-surface) px-4 py-3"
			>
				<svg
					viewBox="0 0 16 16"
					class="mt-0.5 h-4 w-4 flex-shrink-0"
					fill="none"
					stroke="var(--color-status-published)"
					stroke-width="1.7"
					aria-hidden="true"
					focusable="false"
				>
					<path d="M3 8.5 6 11.5 13 4.5" stroke-linecap="round" stroke-linejoin="round" />
				</svg>
				<span class="text-sm leading-relaxed text-(--color-text-muted)"
					>Privacy, duet, stitch never pre-checked</span
				>
			</div>
			<div
				class="flex items-start gap-2 rounded-xl border border-border bg-(--color-surface) px-4 py-3"
			>
				<svg
					viewBox="0 0 16 16"
					class="mt-0.5 h-4 w-4 flex-shrink-0"
					fill="none"
					stroke="var(--color-status-published)"
					stroke-width="1.7"
					aria-hidden="true"
					focusable="false"
				>
					<path d="M3 8.5 6 11.5 13 4.5" stroke-linecap="round" stroke-linejoin="round" />
				</svg>
				<span class="text-sm leading-relaxed text-(--color-text-muted)"
					>Chunked upload straight to TikTok / Instagram</span
				>
			</div>
			<div
				class="flex items-start gap-2 rounded-xl border border-border bg-(--color-surface) px-4 py-3"
			>
				<svg
					viewBox="0 0 16 16"
					class="mt-0.5 h-4 w-4 flex-shrink-0"
					fill="none"
					stroke="var(--color-status-published)"
					stroke-width="1.7"
					aria-hidden="true"
					focusable="false"
				>
					<path d="M3 8.5 6 11.5 13 4.5" stroke-linecap="round" stroke-linejoin="round" />
				</svg>
				<span class="text-sm leading-relaxed text-(--color-text-muted)"
					>Polls until the post is actually live</span
				>
			</div>
		</div>
	</div>
</section>
