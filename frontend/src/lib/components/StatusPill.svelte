<script lang="ts">
	let { status }: { status: string } = $props();

	const labels: Record<string, string> = {
		draft: 'Draft',
		scheduled: 'Scheduled',
		publishing: 'Publishing',
		processing: 'Processing',
		published: 'Published',
		partial: 'Partially published',
		failed: 'Failed',
		pending: 'Pending'
	};
</script>

<span
	class="rounded-full px-2.5 py-1 font-mono text-[11px] font-medium tracking-wide"
	style="color: var(--color-status-{status}); background-color: color-mix(in srgb, var(--color-status-{status}) 14%, transparent)"
>
	{labels[status] ?? status}
</span>
