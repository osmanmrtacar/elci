import type { Component } from 'svelte';
import TikTokIcon from '$lib/components/icons/TikTokIcon.svelte';
import InstagramIcon from '$lib/components/icons/InstagramIcon.svelte';
import XIcon from '$lib/components/icons/XIcon.svelte';
import TikTokSettingsPanel from './tiktok/SettingsPanel.svelte';
import TikTokPreviewPanel from './tiktok/PreviewPanel.svelte';
import InstagramSettingsPanel from './instagram/SettingsPanel.svelte';
import InstagramPreviewPanel from './instagram/PreviewPanel.svelte';
import XSettingsPanel from './x/SettingsPanel.svelte';
import XPreviewPanel from './x/PreviewPanel.svelte';
import type { Platform } from '$lib/types';
import type { SettingsPanelProps, PreviewPanelProps } from './types';

export interface ProviderDefinition {
	id: Platform;
	label: string;
	icon: Component<{ class?: string }>;
	maxCaptionLength: number;
	SettingsPanel: Component<SettingsPanelProps>;
	PreviewPanel: Component<PreviewPanelProps>;
}

export const providers: ProviderDefinition[] = [
	{
		id: 'tiktok',
		label: 'TikTok',
		icon: TikTokIcon,
		maxCaptionLength: 2200,
		SettingsPanel: TikTokSettingsPanel,
		PreviewPanel: TikTokPreviewPanel
	},
	{
		id: 'instagram',
		label: 'Instagram',
		icon: InstagramIcon,
		maxCaptionLength: 2200,
		SettingsPanel: InstagramSettingsPanel,
		PreviewPanel: InstagramPreviewPanel
	},
	{
		id: 'x',
		label: 'X',
		icon: XIcon,
		maxCaptionLength: 280,
		SettingsPanel: XSettingsPanel,
		PreviewPanel: XPreviewPanel
	}
];

export function getProvider(id: Platform): ProviderDefinition {
	const found = providers.find((p) => p.id === id);
	if (!found) throw new Error(`unknown platform: ${id}`);
	return found;
}
