<script lang="ts">
	let { class: className = '' }: { class?: string } = $props();
</script>

<svg
	viewBox="0 0 24 24"
	fill="none"
	stroke="currentColor"
	stroke-width="1.8"
	class={className}
	aria-hidden="true"
>
	<path
		d="M6 7h9.5L13 4.5M18 17H8.5L11 19.5M6 7v6.5a3 3 0 0 0 3 3h1M18 17v-6.5a3 3 0 0 0-3-3h-1"
		stroke-linecap="round"
		stroke-linejoin="round"
	/>
</svg>
