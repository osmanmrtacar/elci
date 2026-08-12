<script lang="ts">
	import { reveal } from '$lib/actions/reveal';
	const tagline = "Your week of posts, finally in one place.";
	const words = tagline.split(' ');
</script>

<section class="border-y border-border bg-(--color-surface)/60">
	<div class="mx-auto max-w-4xl px-6 py-24 text-center md:py-32">
		<p
			use:reveal
			class="reveal font-[family-name:var(--font-display)] text-3xl font-semibold tracking-tight md:text-5xl"
		>
			{#each words as word, i (i)}
				<span class="word inline-block" style="transition-delay: {i * 30}ms">{word}&nbsp;</span>
			{/each}
		</p>
	</div>
</section>

<style>
	.word {
		opacity: 0.28;
		color: var(--color-text-muted);
		will-change: opacity, color;
		transition:
			opacity 0.5s cubic-bezier(0.16, 1, 0.3, 1),
			color 0.5s cubic-bezier(0.16, 1, 0.3, 1);
	}

	:global(.reveal.is-visible) .word {
		opacity: 1;
		color: var(--color-text);
	}
</style>
