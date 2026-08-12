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
	<path d="M4 5h16v11H8.5L4 19.5V5Z" stroke-linecap="round" stroke-linejoin="round" />
</svg>
