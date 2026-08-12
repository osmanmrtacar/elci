<script lang="ts">
	import Header from '$lib/components/landing/Header.svelte';
	import Hero from '$lib/components/landing/Hero.svelte';
	import TaglineReveal from '$lib/components/landing/TaglineReveal.svelte';
	import Benefits from '$lib/components/landing/Benefits.svelte';
	import HowItWorks from '$lib/components/landing/HowItWorks.svelte';
	import Proof from '$lib/components/landing/Proof.svelte';
	import AgentTeaser from '$lib/components/landing/AgentTeaser.svelte';
	import FAQ from '$lib/components/landing/FAQ.svelte';
	import RiskReversal from '$lib/components/landing/RiskReversal.svelte';
	import FinalCTA from '$lib/components/landing/FinalCTA.svelte';
	import Footer from '$lib/components/landing/Footer.svelte';
	import {
		siteUrl,
		defaultTitle,
		defaultDescription,
		jsonLdOrganization,
		jsonLdSoftwareApplication,
		jsonLdFAQ
	} from '$lib/seo';

	const faqData = [
		{
			q: 'Will elci post something I didn’t approve?',
			a: 'Never. Nothing goes out unless you hit Schedule or Publish. No auto-posting in the background.'
		},
		{
			q: 'Do you store my TikTok or Instagram password?',
			a: 'Nope. You log in on TikTok or Instagram directly. We only get a revocable token, never your password.'
		},
		{
			q: 'Why can I only post privately on TikTok right now?',
			a: 'TikTok keeps new apps in a private-only sandbox until they pass a full audit. Daily caps too. It’s TikTok’s rule, same for any app at this stage.'
		},
		{
			q: 'Can I choose who sees each post?',
			a: 'Yep, only from the privacy options your account actually has right now. We never pre-select one for you.'
		},
		{
			q: 'What happens if TikTok or Instagram rejects a post?',
			a: 'We show you the exact error from the platform and keep your draft intact, so you can fix it and try again. No lost captions.'
		},
		{
			q: 'Do I need to upload my video somewhere else first?',
			a: 'No. Just drop it in elci. We do the chunked upload straight to TikTok or Instagram for you.'
		},
		{
			q: 'Can I schedule for later?',
			a: 'Yes. Pick a date and time and we’ll queue it and publish when it hits.'
		},
		{
			q: 'Which accounts can I connect?',
			a: 'One TikTok and one Instagram account, each connected through that platform’s own login.'
		},
		{
			q: 'Can an AI agent post through elci?',
			a: 'Yes, through the MCP server and REST API. Agent posts go through the same privacy and disclosure checks as a human draft.'
		},
		{
			q: 'Do you have a public API?',
			a: 'Yep, REST with OAuth2, plus native support for n8n, Make and Zapier.'
		}
	];

	const orgJsonLd = jsonLdOrganization();
	const softwareJsonLd = jsonLdSoftwareApplication();
	const combinedOrgSoftware = JSON.stringify({
		'@context': 'https://schema.org',
		'@graph': [JSON.parse(orgJsonLd), JSON.parse(softwareJsonLd)]
	});
	const faqJsonLd = jsonLdFAQ(faqData);
</script>

<svelte:head>
	<title>{defaultTitle}</title>
	<meta name="description" content={defaultDescription} />
	<link rel="canonical" href="{siteUrl}/" />
	<meta name="robots" content="index, follow" />
	<meta name="theme-color" content="#fafaf9" />
	<meta property="og:type" content="website" />
	<meta property="og:url" content="{siteUrl}/" />
	<meta property="og:title" content={defaultTitle} />
	<meta property="og:description" content={defaultDescription} />
	<meta property="og:image" content="{siteUrl}/og.png" />
	<meta property="og:image:width" content="1200" />
	<meta property="og:image:height" content="630" />
	<meta name="twitter:card" content="summary_large_image" />
	<meta name="twitter:title" content={defaultTitle} />
	<meta name="twitter:description" content={defaultDescription} />
	<meta name="twitter:image" content="{siteUrl}/og.png" />
	{@html `<script type="application/ld+json">${combinedOrgSoftware}</script>`}
	{@html `<script type="application/ld+json">${faqJsonLd}</script>`}
</svelte:head>

<a
	href="#main-content"
	class="sr-only focus:not-sr-only focus:fixed focus:top-4 focus:left-4 focus:z-50 focus:rounded-full focus:bg-(--color-navy) focus:px-4 focus:py-2 focus:text-white"
	>Skip to content</a
>
<Header />
<main id="main-content">
	<Hero />
	<TaglineReveal />
	<Benefits />
	<HowItWorks />
	<Proof />
	<AgentTeaser />
	<FAQ />
	<RiskReversal />
	<FinalCTA />
</main>
<Footer />
