<script lang="ts">
	let { size = 20 }: { size?: number } = $props();
</script>

<span class="flex shrink-0 items-center gap-[0.3em]" style="font-size: {size}px">
	<span
		class="font-[family-name:var(--font-logo)] leading-none tracking-[-0.005em]"
		style="font-weight: 480; font-optical-sizing: auto;">elci</span
	>
	<svg
		viewBox="2 14 96 56"
		class="mt-[0.06em] h-[0.58em] w-auto shrink-0 origin-[50%_55%] self-start rotate-[-6deg]"
		aria-hidden="true"
		focusable="false"
	>
		<path
			d="M 12 58 Q 32 25 50 48 Q 68 25 88 58"
			fill="none"
			stroke="currentColor"
			stroke-width="8"
			stroke-linecap="round"
			stroke-linejoin="round"
		/>
		<circle cx="50" cy="46" r="3.6" fill="currentColor" />
	</svg>
</span>
