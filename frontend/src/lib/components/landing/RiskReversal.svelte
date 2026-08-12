<script lang="ts">
	import { reveal } from '$lib/actions/reveal';
</script>

<section class="mx-auto max-w-3xl px-6 py-4 text-center">
	<p
		use:reveal
		class="reveal flex items-center justify-center gap-2 text-sm text-(--color-text-muted)"
	>
		<svg
			viewBox="0 0 24 24"
			fill="none"
			stroke="currentColor"
			stroke-width="1.7"
			class="h-4 w-4 flex-shrink-0"
			aria-hidden="true"
			focusable="false"
		>
			<path
				d="M12 3 4 6.5v5c0 4.7 3.2 8.4 8 9.5 4.8-1.1 8-4.8 8-9.5v-5L12 3Z"
				stroke-linecap="round"
				stroke-linejoin="round"
			/>
			<path d="M9 12 11 14 15 10" stroke-linecap="round" stroke-linejoin="round" />
		</svg>
		Disconnect any account at any time. Revoking access takes one click, and elci never posts without
		a draft you approved first.
	</p>
</section>
