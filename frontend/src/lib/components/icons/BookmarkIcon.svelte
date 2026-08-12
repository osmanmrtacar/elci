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
	<path d="M6 3.5h12v17l-6-4-6 4v-17Z" stroke-linecap="round" stroke-linejoin="round" />
</svg>
