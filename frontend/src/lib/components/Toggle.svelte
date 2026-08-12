<script lang="ts">
	let {
		checked,
		onchange,
		disabled = false,
		label
	}: {
		checked: boolean;
		onchange: (value: boolean) => void;
		disabled?: boolean;
		label: string;
	} = $props();
</script>

<label class="flex items-center justify-between gap-3 py-1 text-sm" class:opacity-40={disabled}>
	<span>{label}</span>
	<button
		type="button"
		role="switch"
		aria-checked={checked}
		aria-label={label}
		{disabled}
		onclick={() => onchange(!checked)}
		class="relative h-6 w-10 shrink-0 rounded-full transition-colors disabled:cursor-not-allowed"
		style="background-color: {checked ? 'var(--color-accent)' : 'var(--color-border)'}"
	>
		<span
			class="absolute top-0.5 left-0.5 h-5 w-5 rounded-full bg-white shadow transition-transform"
			style="transform: translateX({checked ? '16px' : '0px'})"
		></span>
	</button>
</label>
