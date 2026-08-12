export type AnalyticsEvent =
	'header-cta' | 'hero-primary' | 'howitworks-cta' | 'agent-teaser-cta' | 'final-cta';

type GtagFn = (command: string, event: string, params?: Record<string, string>) => void;
type PlausibleFn = (event: string, opts?: { props?: Record<string, string> }) => void;

export function track(event: AnalyticsEvent, extra?: Record<string, string>) {
	if (typeof window === 'undefined') return;
	try {
		const gtag = (window as unknown as { gtag?: GtagFn }).gtag;
		gtag?.('event', event, extra);
	} catch {}
	try {
		const plausible = (window as unknown as { plausible?: PlausibleFn }).plausible;
		plausible?.(event, { props: extra });
	} catch {}
	if (import.meta.env.DEV) {
		console.debug('[analytics]', event, extra);
	}
}
