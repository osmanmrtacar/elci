<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { SvelteSet } from 'svelte/reactivity';
	import { providers, getProvider } from '$lib/providers/registry';
	import { accountInfo as fetchAccountInfo } from '$lib/api/connections';
	import { uploadMedia } from '$lib/api/media';
	import { createPost, deletePost, publishPost, schedulePost } from '$lib/api/posts';
	import type { TargetInput } from '$lib/api/posts';
	import { ApiError } from '$lib/api/client';
	import type { AccountInfo, Connection, MediaKind, Platform, ValidationErrors } from '$lib/types';

	let { connections }: { connections: Connection[] } = $props();

	function connectionFor(id: Platform) {
		return connections.find((c) => c.platform === id && c.isActive);
	}
	let connectedIds = $derived(
		new Set(connections.filter((c) => c.isActive).map((c) => c.platform))
	);

	let defaultCaption = $state('');
	let media = $state<{ id: number; publicUrl: string; kind: MediaKind } | null>(null);
	let uploading = $state(false);
	let uploadError = $state<string | null>(null);

	let selected = new SvelteSet<Platform>();
	let overrides = $state<
		Record<string, { captionOverride: string | null; settings: Record<string, unknown> }>
	>({});
	let activeTab = $state<'global' | Platform>('global');

	let accountInfoCache = $state<Partial<Record<Platform, AccountInfo>>>({});
	let infoErrors = $state<Partial<Record<Platform, string>>>({});
	let infoLoading = $state<Partial<Record<Platform, boolean>>>({});

	let draftId = $state<number | null>(null);
	let submitError = $state<string | null>(null);
	let validationErrors = $state<ValidationErrors['errors'] | null>(null);
	let submitting = $state(false);
	let scheduleAt = $state('');

	function togglePlatform(id: Platform) {
		if (!connectedIds.has(id)) return;
		if (selected.has(id)) {
			selected.delete(id);
			delete overrides[id];
			if (activeTab === id) activeTab = 'global';
		} else {
			selected.add(id);
			overrides[id] = { captionOverride: null, settings: {} };
		}
	}

	function effectiveCaption(id: Platform) {
		return overrides[id]?.captionOverride ?? defaultCaption;
	}

	// Fetch live account info for every selected platform as soon as it's
	// selected — not just whichever tab happens to be active — so every
	// stacked preview shows the real handle immediately rather than a
	// placeholder until its settings tab has been opened at least once.
	$effect(() => {
		for (const id of selected) {
			if (accountInfoCache[id] || infoLoading[id]) continue;
			infoLoading[id] = true;
			delete infoErrors[id];
			fetchAccountInfo(id)
				.then((info) => {
					accountInfoCache[id] = info;
				})
				.catch((e) => {
					infoErrors[id] = e instanceof Error ? e.message : 'failed to load';
				})
				.finally(() => {
					infoLoading[id] = false;
				});
		}
	});

	// TikTok requires the publish action itself to be disabled (not just
	// warned) when content disclosure is on but neither brand option is
	// picked. Written generically over every selected platform's settings
	// rather than hardcoded to "tiktok" — it simply never matches for
	// platforms that don't have these keys.
	let blocksPublish = $derived(
		[...selected].some((id) => {
			const s = overrides[id]?.settings ?? {};
			return Boolean(s.disclose_content) && !s.your_brand && !s.branded_content;
		})
	);

	async function handleFile(e: Event) {
		const file = (e.target as HTMLInputElement).files?.[0];
		if (!file) return;
		uploading = true;
		uploadError = null;
		const kind: MediaKind = file.type.startsWith('video/') ? 'video' : 'image';
		try {
			media = await uploadMedia(file, kind);
		} catch (err) {
			uploadError = err instanceof Error ? err.message : 'upload failed';
		} finally {
			uploading = false;
		}
	}

	async function submit(action: 'draft' | 'schedule' | 'publish') {
		submitError = null;
		validationErrors = null;
		if (selected.size === 0) {
			submitError = 'Select at least one platform.';
			return;
		}
		if (action === 'schedule' && !scheduleAt) {
			submitError = 'Pick a date and time to schedule for.';
			return;
		}

		submitting = true;
		try {
			if (draftId) {
				await deletePost(draftId).catch(() => {});
				draftId = null;
			}

			const targets: TargetInput[] = [...selected].map((id) => ({
				platform: id,
				captionOverride: overrides[id]?.captionOverride ?? null,
				settings: overrides[id]?.settings ?? {}
			}));
			const post = await createPost({
				defaultCaption,
				defaultMediaKind: media?.kind ?? null,
				defaultMediaUrls: media ? [media.publicUrl] : [],
				targets
			});
			draftId = post.id;

			if (action === 'schedule') {
				await schedulePost(post.id, new Date(scheduleAt));
			} else if (action === 'publish') {
				await publishPost(post.id);
			}

			if (action === 'draft') {
				goto(resolve('/dashboard'));
			} else {
				goto(resolve('/(app)/posts/[id]', { id: String(post.id) }));
			}
		} catch (err) {
			if (
				err instanceof ApiError &&
				err.status === 422 &&
				err.body &&
				typeof err.body === 'object' &&
				'errors' in err.body
			) {
				validationErrors = (err.body as ValidationErrors).errors;
				const firstBad = Object.keys(validationErrors)[0] as Platform | undefined;
				if (firstBad) activeTab = firstBad;
			} else {
				submitError = err instanceof Error ? err.message : 'Something went wrong.';
			}
		} finally {
			submitting = false;
		}
	}
</script>

<div class="mb-6">
	<span class="mb-2 block text-sm font-medium">Post to</span>
	<div class="flex flex-wrap gap-2">
		{#each providers as provider (provider.id)}
			{@const conn = connectionFor(provider.id)}
			{@const connected = Boolean(conn)}
			<button
				type="button"
				disabled={!connected}
				onclick={() => togglePlatform(provider.id)}
				class="flex items-center gap-2 rounded-xl border px-3 py-2 text-left transition disabled:cursor-not-allowed disabled:opacity-40 {selected.has(
					provider.id
				)
					? 'border-(--color-accent) bg-(--color-accent)/10'
					: 'border-border hover:border-(--color-text-muted)'}"
			>
				<provider.icon
					class="h-4 w-4 shrink-0 {selected.has(provider.id) ? 'text-(--color-accent)' : ''}"
				/>
				<span class="leading-tight">
					<span
						class="block text-sm font-medium {selected.has(provider.id)
							? 'text-(--color-accent)'
							: 'text-(--color-text)'}">{provider.label}</span
					>
					<span class="block font-mono text-[11px] text-(--color-text-muted)">
						{conn ? `@${conn.username}` : 'not connected'}
					</span>
				</span>
			</button>
		{/each}
	</div>
</div>

<div class="grid gap-8 lg:grid-cols-[1fr_360px]">
	<div class="rounded-2xl border border-border bg-(--color-surface) p-6">
		<div class="mb-5 flex flex-wrap gap-1 border-b border-border">
			<button
				type="button"
				onclick={() => (activeTab = 'global')}
				class="border-b-2 px-3 py-2 text-sm font-semibold {activeTab === 'global'
					? 'border-(--color-accent) text-(--color-text)'
					: 'border-transparent text-(--color-text-muted)'}"
			>
				Global
			</button>
			{#each [...selected] as id (id)}
				{@const p = getProvider(id)}
				<button
					type="button"
					onclick={() => (activeTab = id)}
					class="flex items-center gap-1.5 border-b-2 px-3 py-2 text-sm font-semibold {activeTab ===
					id
						? 'border-(--color-accent) text-(--color-text)'
						: 'border-transparent text-(--color-text-muted)'}"
				>
					<p.icon class="h-3.5 w-3.5" />
					{p.label}
				</button>
			{/each}
		</div>

		{#if activeTab === 'global'}
			<div class="space-y-5">
				<div>
					<label for="caption" class="mb-1 block text-sm font-medium">Caption</label>
					<textarea
						id="caption"
						bind:value={defaultCaption}
						rows="5"
						placeholder="What do you want to say?"
						class="w-full rounded-lg border border-border bg-(--color-bg) px-3 py-2 text-sm focus:border-(--color-accent) focus:outline-none"
					></textarea>
				</div>

				<div>
					<span class="mb-1 block text-sm font-medium">Media</span>
					{#if media}
						<div class="flex items-center gap-3 rounded-lg border border-border p-3">
							{#if media.kind === 'video'}
								<video src={media.publicUrl} class="h-16 w-16 rounded-md object-cover" muted
								></video>
							{:else}
								<img src={media.publicUrl} alt="" class="h-16 w-16 rounded-md object-cover" />
							{/if}
							<button
								type="button"
								onclick={() => (media = null)}
								class="text-xs text-(--color-text-muted) hover:text-(--color-status-failed)"
							>
								Remove
							</button>
						</div>
					{:else}
						<label
							class="flex cursor-pointer items-center gap-2 rounded-lg border border-dashed border-border px-3 py-2.5 text-sm text-(--color-text-muted) hover:border-(--color-accent) hover:text-(--color-text) {uploading
								? 'pointer-events-none opacity-60'
								: ''}"
						>
							<svg
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
								stroke-width="1.6"
								class="h-4 w-4 shrink-0"
							>
								<path
									d="M12 16V4m0 0 4 4m-4-4-4 4"
									stroke-linecap="round"
									stroke-linejoin="round"
								/>
								<path
									d="M4 16v3a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-3"
									stroke-linecap="round"
									stroke-linejoin="round"
								/>
							</svg>
							{uploading ? 'Uploading…' : 'Upload an image or video'}
							<input
								type="file"
								accept="image/*,video/*"
								onchange={handleFile}
								disabled={uploading}
								class="sr-only"
							/>
						</label>
						{#if uploadError}<p class="mt-1 text-xs text-(--color-status-failed)">
								{uploadError}
							</p>{/if}
					{/if}
				</div>
			</div>
		{:else}
			{@const id = activeTab}
			{@const provider = getProvider(id)}
			{@const usingDefault = overrides[id]?.captionOverride === null}
			<div class="space-y-5">
				{#if validationErrors?.[id]}
					<div
						class="rounded-lg bg-(--color-status-failed)/10 p-3 text-sm text-(--color-status-failed)"
					>
						{#each validationErrors[id] as err (err.field)}
							<p>{err.message}</p>
						{/each}
					</div>
				{/if}

				<label class="flex items-center gap-2 text-sm">
					<input
						type="checkbox"
						checked={usingDefault}
						onchange={(e) => {
							overrides[id].captionOverride = e.currentTarget.checked ? null : defaultCaption;
						}}
					/>
					Use the global caption
				</label>
				{#if !usingDefault}
					<textarea
						bind:value={overrides[id].captionOverride}
						rows="4"
						class="w-full rounded-lg border border-border bg-(--color-bg) px-3 py-2 text-sm focus:border-(--color-accent) focus:outline-none"
					></textarea>
				{/if}

				<provider.SettingsPanel
					caption={effectiveCaption(id)}
					mediaKind={media?.kind ?? null}
					bind:settings={overrides[id].settings}
					accountInfo={accountInfoCache[id] ?? null}
					infoError={infoErrors[id] ?? null}
				/>
			</div>
		{/if}

		{#if submitError}
			<p class="mt-4 text-sm text-(--color-status-failed)">{submitError}</p>
		{/if}
		{#if blocksPublish}
			<p class="mt-4 text-sm text-(--color-status-failed)">
				Resolve the content-disclosure settings on the platform tab above before you can save,
				schedule, or publish.
			</p>
		{/if}

		<div class="mt-8 flex flex-wrap items-center justify-between gap-3 border-t border-border pt-6">
			<button
				type="button"
				disabled={submitting || blocksPublish}
				onclick={() => submit('draft')}
				class="rounded-full border border-border px-4 py-2 text-sm font-semibold disabled:opacity-50"
			>
				Save draft
			</button>
			<div class="flex flex-wrap items-center gap-3">
				<input
					type="datetime-local"
					bind:value={scheduleAt}
					class="rounded-full border border-border bg-(--color-bg) px-3 py-2 font-mono text-sm tabular-nums"
				/>
				<button
					type="button"
					disabled={submitting || blocksPublish}
					onclick={() => submit('schedule')}
					class="rounded-full border border-border px-4 py-2 text-sm font-semibold disabled:opacity-50"
				>
					Schedule
				</button>
				<button
					type="button"
					disabled={submitting || blocksPublish}
					onclick={() => submit('publish')}
					class="rounded-full bg-(--color-navy) px-4 py-2 text-sm font-semibold text-white disabled:opacity-50"
				>
					Publish now
				</button>
			</div>
		</div>
	</div>

	<div class="space-y-6 lg:sticky lg:top-24 lg:self-start">
		{#if selected.size === 0}
			<div
				class="rounded-2xl border border-dashed border-border p-8 text-center text-sm text-(--color-text-muted)"
			>
				Select a platform above and start writing to see a preview.
			</div>
		{:else}
			{#each [...selected] as id (id)}
				{@const provider = getProvider(id)}
				<div class="relative">
					<span
						class="absolute -top-2 -right-2 z-10 grid h-7 w-7 place-items-center rounded-full border border-border bg-(--color-surface) shadow-sm"
					>
						<provider.icon class="h-3.5 w-3.5" />
					</span>
					{#if validationErrors?.[id]}
						<div
							class="mb-2 rounded-lg bg-(--color-status-failed)/10 p-3 text-sm text-(--color-status-failed)"
						>
							{#each validationErrors[id] as err (err.field)}
								<p>{err.message}</p>
							{/each}
						</div>
					{/if}
					<provider.PreviewPanel
						caption={effectiveCaption(id)}
						mediaKind={media?.kind ?? null}
						mediaUrl={media?.publicUrl ?? null}
						settings={overrides[id]?.settings ?? {}}
						accountInfo={accountInfoCache[id] ?? null}
					/>
				</div>
			{/each}
		{/if}
	</div>
</div>
