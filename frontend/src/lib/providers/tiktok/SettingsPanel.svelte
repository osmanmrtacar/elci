<script lang="ts">
	import Toggle from '$lib/components/Toggle.svelte';
	import type { SettingsPanelProps } from '../types';

	let {
		caption,
		mediaKind,
		settings = $bindable(),
		accountInfo,
		infoError
	}: SettingsPanelProps = $props();

	const privacyLabels: Record<string, string> = {
		PUBLIC_TO_EVERYONE: 'Public',
		MUTUAL_FOLLOW_FRIENDS: 'Friends',
		FOLLOWER_OF_CREATOR: 'Followers',
		SELF_ONLY: 'Only me'
	};

	let privacyOptions = $derived(
		(accountInfo?.extra?.privacy_level_options as string[] | undefined) ?? []
	);
	let commentDisabled = $derived(Boolean(accountInfo?.extra?.comment_disabled));
	let duetDisabled = $derived(Boolean(accountInfo?.extra?.duet_disabled));
	let stitchDisabled = $derived(Boolean(accountInfo?.extra?.stitch_disabled));
	let maxDurationSec = $derived(
		accountInfo?.extra?.max_video_post_duration_sec as number | undefined
	);

	let privacyLevel = $derived((settings.privacy_level as string | undefined) ?? '');
	let comment = $derived(Boolean(settings.comment));
	let duet = $derived(Boolean(settings.duet));
	let stitch = $derived(Boolean(settings.stitch));
	let autoAddMusic = $derived(Boolean(settings.auto_add_music));
	let discloseContent = $derived(Boolean(settings.disclose_content));
	let yourBrand = $derived(Boolean(settings.your_brand));
	let brandedContent = $derived(Boolean(settings.branded_content));

	let maxCaptionLen = $derived(mediaKind === 'image' ? 4000 : 2200);

	// "Your photo/video will be labeled as..." — distinct from the "by
	// posting you agree to..." consent line further down; TikTok requires
	// both to be shown.
	let labelPreview = $derived(
		brandedContent
			? "Your photo/video will be labeled as 'Paid partnership'"
			: yourBrand
				? "Your photo/video will be labeled as 'Promotional content'"
				: null
	);

	$effect(() => {
		if (commentDisabled) settings.comment = false;
		if (mediaKind === 'video') {
			if (duetDisabled) settings.duet = false;
			if (stitchDisabled) settings.stitch = false;
		} else {
			settings.duet = false;
			settings.stitch = false;
		}
		// Branded content and a private ("only me") audience are mutually
		// exclusive — if the account's live options mean one of these is
		// already true when the other becomes true, force it back off
		// rather than leaving an invalid combination in place.
		if (brandedContent && privacyLevel === 'SELF_ONLY') settings.branded_content = false;
	});
</script>

<div class="space-y-5 text-sm">
	{#if infoError}
		<p class="rounded-lg bg-(--color-status-failed)/10 p-3 text-(--color-status-failed)">
			Couldn't load live TikTok account settings: {infoError}
		</p>
	{:else if !accountInfo}
		<p class="text-(--color-text-muted)">Loading this account's current TikTok settings…</p>
	{:else}
		<p class="text-(--color-text-muted)">
			Posting as <span class="font-medium text-(--color-text)"
				>{accountInfo.displayName || accountInfo.username}</span
			>
		</p>

		{#if caption.length > maxCaptionLen}
			<p class="text-(--color-status-failed)">
				Caption is {caption.length - maxCaptionLen} characters over TikTok's {maxCaptionLen} limit.
			</p>
		{/if}

		<div>
			<label for="tiktok-privacy" class="mb-1 block font-medium">Who can view this</label>
			<div class="relative">
				<select
					id="tiktok-privacy"
					value={privacyLevel}
					onchange={(e) => (settings.privacy_level = e.currentTarget.value)}
					class="w-full appearance-none rounded-lg border border-border bg-(--color-bg) px-3 py-2 pr-9 focus:border-(--color-accent) focus:outline-none"
				>
					<option value="" disabled>Select a privacy setting…</option>
					{#each privacyOptions as level (level)}
						<option
							value={level}
							disabled={level === 'SELF_ONLY' && brandedContent}
							title={level === 'SELF_ONLY' && brandedContent
								? 'Branded content visibility cannot be set to private.'
								: undefined}
						>
							{privacyLabels[level] ?? level}
						</option>
					{/each}
				</select>
				<svg
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="1.8"
					class="pointer-events-none absolute top-1/2 right-3 h-3.5 w-3.5 -translate-y-1/2 text-(--color-text-muted)"
				>
					<path d="m6 9 6 6 6-6" stroke-linecap="round" stroke-linejoin="round" />
				</svg>
			</div>
		</div>

		<div class="space-y-1 rounded-lg border border-border p-3">
			<Toggle
				checked={comment}
				onchange={(v: boolean) => (settings.comment = v)}
				disabled={commentDisabled}
				label="Allow comments"
			/>
			{#if mediaKind === 'video'}
				<Toggle
					checked={duet}
					onchange={(v: boolean) => (settings.duet = v)}
					disabled={duetDisabled}
					label="Allow Duet"
				/>
				<Toggle
					checked={stitch}
					onchange={(v: boolean) => (settings.stitch = v)}
					disabled={stitchDisabled}
					label="Allow Stitch"
				/>
			{:else}
				<Toggle
					checked={autoAddMusic}
					onchange={(v: boolean) => (settings.auto_add_music = v)}
					label="Automatically add music"
				/>
			{/if}
		</div>

		{#if mediaKind === 'video' && maxDurationSec}
			<p class="text-xs text-(--color-text-muted)">
				This account's videos can be up to {maxDurationSec}s long.
			</p>
		{/if}

		<div class="rounded-lg border border-border p-3">
			<Toggle
				checked={discloseContent}
				onchange={(v: boolean) => (settings.disclose_content = v)}
				label="Disclose {mediaKind === 'video' ? 'video' : 'photo'} content"
			/>
			<p class="mt-1 text-xs text-(--color-text-muted)">
				Turn on to disclose that this content promotes goods or services in exchange for something
				of value. Your content could promote yourself, a third party, or both.
			</p>

			{#if discloseContent}
				<div
					class="mt-3 space-y-2 border-t border-border pt-3"
					title={!yourBrand && !brandedContent
						? 'You need to indicate if your content promotes yourself, a third party, or both.'
						: undefined}
				>
					<label class="flex items-start gap-2">
						<input
							type="checkbox"
							checked={yourBrand}
							onchange={(e) => (settings.your_brand = e.currentTarget.checked)}
							class="mt-0.5"
						/>
						<span>
							<strong>Your Brand</strong> — you are promoting yourself or your own business.
						</span>
					</label>
					<label
						class="flex items-start gap-2"
						class:opacity-40={privacyLevel === 'SELF_ONLY'}
						title={privacyLevel === 'SELF_ONLY'
							? 'Branded content visibility cannot be set to private.'
							: undefined}
					>
						<input
							type="checkbox"
							checked={brandedContent}
							disabled={privacyLevel === 'SELF_ONLY'}
							onchange={(e) => (settings.branded_content = e.currentTarget.checked)}
							class="mt-0.5"
						/>
						<span>
							<strong>Branded Content</strong> — you are promoting another brand or a third party.
						</span>
					</label>
					{#if !yourBrand && !brandedContent}
						<p class="text-xs text-(--color-status-failed)">
							Select at least one of "Your Brand" or "Branded Content".
						</p>
					{/if}
					{#if labelPreview}
						<p class="text-xs text-(--color-text-muted)">{labelPreview}</p>
					{/if}
				</div>
			{/if}
		</div>

		<p class="text-xs text-(--color-text-muted)">
			By posting, you agree to TikTok's{brandedContent ? ' Branded Content Policy and' : ''} Music Usage
			Confirmation.
		</p>

		<p class="text-xs text-(--color-text-muted)">
			After you publish, it can take a few minutes for TikTok to finish processing before this
			appears on the profile above.
		</p>
	{/if}
</div>
