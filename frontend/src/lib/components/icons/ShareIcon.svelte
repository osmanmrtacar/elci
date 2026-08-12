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
	<path d="M4 12 20.5 4 16 20l-4-6.5L4 12Z" stroke-linecap="round" stroke-linejoin="round" />
</svg>
