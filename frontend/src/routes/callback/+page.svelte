<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { auth } from '$lib/auth.svelte';

	onMount(() => {
		const token = page.url.searchParams.get('token');
		if (token) {
			auth.login(token);
			goto(resolve('/dashboard'));
		} else {
			goto(resolve('/login?error=connection_failed'));
		}
	});
</script>

<main class="flex min-h-screen items-center justify-center">
	<p class="text-sm text-(--color-text-muted)">Finishing sign-in…</p>
</main>
