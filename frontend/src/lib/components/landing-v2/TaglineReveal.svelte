<script lang="ts">
	import { reveal } from '$lib/actions/reveal';

	const tagline = "Checked against every platform's own rules — before it ever ships.";
	const words = tagline.split(' ');
</script>

<section class="border-y border-border bg-(--color-surface)/60">
	<div class="mx-auto max-w-4xl px-6 py-24 text-center">
		<p use:reveal class="reveal text-3xl font-semibold tracking-tight md:text-5xl">
			{#each words as word, i (i)}
				<span class="word inline-block" style="transition-delay: {i * 45}ms">{word}&nbsp;</span>
			{/each}
		</p>
	</div>
</section>

<style>
	p {
		font-family: var(--font-serif);
	}

	.word {
		opacity: 0.28;
		color: var(--color-text-muted);
		transition:
			opacity 0.5s cubic-bezier(0.32, 0.72, 0, 1),
			color 0.5s cubic-bezier(0.32, 0.72, 0, 1);
	}

	:global(.reveal.is-visible) .word {
		opacity: 1;
		color: var(--color-text);
	}
</style>
