import type { AccountInfo, MediaKind } from '$lib/types';

// The shared contract every platform's SettingsPanel/PreviewPanel implements
// — this is what lets the Composer mount whichever one is active without
// knowing anything platform-specific itself.
export interface SettingsPanelProps {
	caption: string;
	mediaKind: MediaKind | null;
	settings: Record<string, unknown>;
	accountInfo: AccountInfo | null;
	infoError: string | null;
}

export interface PreviewPanelProps {
	caption: string;
	mediaKind: MediaKind | null;
	mediaUrl: string | null;
	settings: Record<string, unknown>;
	accountInfo: AccountInfo | null;
}
