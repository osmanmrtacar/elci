import { env } from '$env/dynamic/public';

export const siteUrl = (env.PUBLIC_SITE_URL as string | undefined) ?? 'https://elci.io';
export const defaultTitle = 'elci - One place to plan and publish your socials';
export const defaultDescription =
	'Plan this week in one calendar, not three apps. Write once in elci and publish to TikTok and Instagram today, more platforms coming to the same flow.';

export function jsonLdOrganization(): string {
	return JSON.stringify({
		'@context': 'https://schema.org',
		'@type': 'Organization',
		name: 'elci',
		url: siteUrl,
		logo: `${siteUrl}/og.png`,
		sameAs: []
	});
}

export function jsonLdSoftwareApplication(): string {
	return JSON.stringify({
		'@context': 'https://schema.org',
		'@type': 'SoftwareApplication',
		name: 'elci',
		applicationCategory: 'SocialNetworkingApplication',
		operatingSystem: 'Web',
		url: siteUrl,
		description: defaultDescription,
		offers: {
			'@type': 'Offer',
			price: '0',
			priceCurrency: 'USD'
		}
	});
}

export function jsonLdFAQ(faqs: { q: string; a: string }[]): string {
	return JSON.stringify({
		'@context': 'https://schema.org',
		'@type': 'FAQPage',
		mainEntity: faqs.map((f) => ({
			'@type': 'Question',
			name: f.q,
			acceptedAnswer: {
				'@type': 'Answer',
				text: f.a
			}
		}))
	});
}
