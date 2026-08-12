<script lang="ts">
	import './redesign.css';
	import Header from '$lib/components/landing-v2/Header.svelte';
	import Hero from '$lib/components/landing-v2/Hero.svelte';
	import TaglineReveal from '$lib/components/landing-v2/TaglineReveal.svelte';
	import Benefits from '$lib/components/landing-v2/Benefits.svelte';
	import HowItWorks from '$lib/components/landing-v2/HowItWorks.svelte';
	import AgentAPI from '$lib/components/landing-v2/AgentAPI.svelte';
	import Proof from '$lib/components/landing-v2/Proof.svelte';
	import Pricing from '$lib/components/landing-v2/Pricing.svelte';
	import FAQ from '$lib/components/landing-v2/FAQ.svelte';
	import RiskReversal from '$lib/components/landing-v2/RiskReversal.svelte';
	import FinalCTA from '$lib/components/landing-v2/FinalCTA.svelte';
	import Footer from '$lib/components/landing-v2/Footer.svelte';
</script>

<svelte:head>
	<title>elci — Social scheduler built for humans and AI agents</title>
	<meta
		name="description"
		content="elci is a scheduler built directly against TikTok, Instagram, and X's own publishing APIs. Post from the dashboard or through elci's REST API and MCP server — every AI agent or automation still passes the same platform rule checks, including TikTok's exact privacy and disclosure requirements."
	/>
</svelte:head>

<div class="redesign">
	<Header />
	<main>
		<Hero />
		<TaglineReveal />
		<Benefits />
		<HowItWorks />
		<AgentAPI />
		<Proof />
		<Pricing />
		<FAQ />
		<RiskReversal />
		<FinalCTA />
	</main>
	<Footer />
</div>
