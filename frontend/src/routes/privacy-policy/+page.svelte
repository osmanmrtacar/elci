<script lang="ts">
	import Header from '$lib/components/landing/Header.svelte';
	import Footer from '$lib/components/landing/Footer.svelte';
</script>

<svelte:head>
	<title>Privacy Policy — elci</title>
</svelte:head>

<Header />
<main class="mx-auto max-w-2xl px-6 py-16">
	<h1 class="text-3xl font-bold tracking-tight">Privacy Policy</h1>
	<p class="mt-2 text-sm text-(--color-text-muted)">Last updated December 27, 2025</p>

	<div class="mt-10 space-y-8 text-sm leading-relaxed text-(--color-text-muted)">
		<section>
			<h2 class="text-base font-semibold text-(--color-text)">Introduction</h2>
			<p class="mt-2">
				Welcome to Elci.io ("Service", "Application", "we", "us", or "our"). This Privacy Policy
				explains how we collect, use, disclose, and safeguard your information when you use Elci.io,
				our social media publisher application that integrates with TikTok, Instagram (Meta), and X
				(formerly Twitter).
			</p>
			<p class="mt-2">
				By using Elci.io, you agree to the collection and use of information in accordance with this
				policy.
			</p>
		</section>

		<section>
			<h2 class="text-base font-semibold text-(--color-text)">Information We Collect</h2>

			<h3 class="mt-4 text-sm font-semibold text-(--color-text)">1. Information You Provide</h3>
			<p class="mt-2"><strong class="text-(--color-text)">Account Information:</strong></p>
			<ul class="mt-2 list-disc space-y-1 pl-5">
				<li>OAuth tokens from TikTok, Instagram, and X platforms</li>
				<li>Username and display name from connected platforms</li>
				<li>Profile information (avatar, bio) from connected platforms</li>
				<li>Email address (if provided by the platform)</li>
			</ul>
			<p class="mt-3"><strong class="text-(--color-text)">Content Information:</strong></p>
			<ul class="mt-2 list-disc space-y-1 pl-5">
				<li>Videos and media files you upload</li>
				<li>Captions and text content you create</li>
				<li>Post scheduling preferences</li>
				<li>Platform selection choices</li>
			</ul>

			<h3 class="mt-4 text-sm font-semibold text-(--color-text)">
				2. Information We Collect Automatically
			</h3>
			<p class="mt-2"><strong class="text-(--color-text)">Usage Data:</strong></p>
			<ul class="mt-2 list-disc space-y-1 pl-5">
				<li>Log data (IP address, browser type, pages visited)</li>
				<li>Device information (device type, operating system)</li>
				<li>Post publishing history and status</li>
				<li>API interaction logs</li>
			</ul>
			<p class="mt-3"><strong class="text-(--color-text)">Cookies and Tracking:</strong></p>
			<ul class="mt-2 list-disc space-y-1 pl-5">
				<li>Session cookies for authentication</li>
				<li>Analytics cookies to improve our service</li>
				<li>OAuth state tokens for security</li>
			</ul>
		</section>

		<section>
			<h2 class="text-base font-semibold text-(--color-text)">How We Use Your Information</h2>

			<h3 class="mt-4 text-sm font-semibold text-(--color-text)">Platform Integration</h3>
			<ul class="mt-2 list-disc space-y-1 pl-5">
				<li><strong class="text-(--color-text)">TikTok:</strong> To publish videos to your TikTok account as authorized by you</li>
				<li><strong class="text-(--color-text)">Instagram:</strong> To publish media to your Instagram Business or Creator account</li>
				<li><strong class="text-(--color-text)">X (Twitter):</strong> To publish posts and media to your X account</li>
			</ul>

			<h3 class="mt-4 text-sm font-semibold text-(--color-text)">Service Operations</h3>
			<ul class="mt-2 list-disc space-y-1 pl-5">
				<li>Authenticate your identity across platforms</li>
				<li>Process and publish your content to selected platforms</li>
				<li>Monitor post status and provide publishing confirmations</li>
				<li>Maintain and improve Elci.io</li>
				<li>Provide customer support</li>
				<li>Detect and prevent fraud or abuse</li>
			</ul>
		</section>

		<section>
			<h2 class="text-base font-semibold text-(--color-text)">Data Sharing and Third Parties</h2>

			<h3 class="mt-4 text-sm font-semibold text-(--color-text)">Third-Party Platforms</h3>
			<p class="mt-2">
				Elci.io shares your information with the following third-party platforms only as necessary
				to provide our Service:
			</p>

			<p class="mt-3"><strong class="text-(--color-text)">TikTok:</strong></p>
			<ul class="mt-2 list-disc space-y-1 pl-5">
				<li>We share video content and captions you authorize us to publish</li>
				<li>User consent is obtained before posting with clear disclosure</li>
				<li>We comply with TikTok's Developer Guidelines</li>
			</ul>

			<p class="mt-3"><strong class="text-(--color-text)">Instagram (Meta):</strong></p>
			<ul class="mt-2 list-disc space-y-1 pl-5">
				<li>We share media and captions you authorize us to publish</li>
				<li>All data processing complies with Meta's Platform Terms</li>
				<li>We do not share data of minors or sensitive personal data</li>
			</ul>

			<p class="mt-3"><strong class="text-(--color-text)">X (Twitter):</strong></p>
			<ul class="mt-2 list-disc space-y-1 pl-5">
				<li>We share posts and media you authorize us to publish</li>
				<li>All usage respects users' reasonable expectations of privacy</li>
				<li>We do not use data for surveillance, tracking, or monitoring purposes</li>
			</ul>
		</section>

		<section>
			<h2 class="text-base font-semibold text-(--color-text)">Your Rights and Choices</h2>

			<h3 class="mt-4 text-sm font-semibold text-(--color-text)">Data Access and Portability</h3>
			<p class="mt-2">You have the right to:</p>
			<ul class="mt-2 list-disc space-y-1 pl-5">
				<li>Access the personal information we hold about you</li>
				<li>Request a copy of your data in a portable format</li>
				<li>Review your posting history and connected platforms</li>
			</ul>

			<h3 class="mt-4 text-sm font-semibold text-(--color-text)">Data Deletion</h3>
			<p class="mt-2">You can request deletion of your data by:</p>
			<ul class="mt-2 list-disc space-y-1 pl-5">
				<li>Disconnecting platforms from your account settings</li>
				<li>
					Contacting us at
					<a href="mailto:support@elci.io" class="text-(--color-accent)">support@elci.io</a>
				</li>
				<li>We will delete your data within 30 days of your request</li>
			</ul>

			<h3 class="mt-4 text-sm font-semibold text-(--color-text)">GDPR Rights (EU Users)</h3>
			<ul class="mt-2 list-disc space-y-1 pl-5">
				<li>Right to rectification of inaccurate data</li>
				<li>Right to erasure ("right to be forgotten")</li>
				<li>Right to restrict processing</li>
				<li>Right to data portability</li>
				<li>Right to object to processing</li>
				<li>Right to withdraw consent</li>
			</ul>

			<h3 class="mt-4 text-sm font-semibold text-(--color-text)">CCPA Rights (California Users)</h3>
			<ul class="mt-2 list-disc space-y-1 pl-5">
				<li>Know what personal information is collected</li>
				<li>Know whether personal information is sold or disclosed</li>
				<li>Say no to the sale of personal information</li>
				<li>Access your personal information</li>
				<li>Request deletion of personal information</li>
				<li>Equal service and price, even if you exercise your privacy rights</li>
			</ul>
		</section>

		<section>
			<h2 class="text-base font-semibold text-(--color-text)">Data Security</h2>
			<p class="mt-2">
				We implement appropriate technical and organizational measures to protect your information:
			</p>
			<ul class="mt-2 list-disc space-y-1 pl-5">
				<li>Encryption of data in transit (HTTPS/TLS)</li>
				<li>Encrypted storage of sensitive data</li>
				<li>Access controls and authentication</li>
				<li>Regular security audits and updates</li>
			</ul>
			<p class="mt-3">
				However, no method of transmission over the Internet is 100% secure. While we strive to
				protect your information, we cannot guarantee its absolute security.
			</p>
		</section>

		<section>
			<h2 class="text-base font-semibold text-(--color-text)">Platform-Specific Policies</h2>
			<p class="mt-2">Your use of our Service is also governed by the privacy policies of:</p>
			<ul class="mt-2 list-disc space-y-1 pl-5">
				<li>
					<a
						href="https://www.tiktok.com/legal/page/eea/privacy-policy/en"
						target="_blank"
						rel="noopener noreferrer"
						class="text-(--color-accent)">TikTok Privacy Policy</a
					>
				</li>
				<li>
					<a
						href="https://www.facebook.com/privacy/policy/"
						target="_blank"
						rel="noopener noreferrer"
						class="text-(--color-accent)">Meta Privacy Policy</a
					> (for Instagram)
				</li>
				<li>
					<a
						href="https://twitter.com/en/privacy"
						target="_blank"
						rel="noopener noreferrer"
						class="text-(--color-accent)">X Privacy Policy</a
					>
				</li>
			</ul>
		</section>

		<section>
			<h2 class="text-base font-semibold text-(--color-text)">Children's Privacy</h2>
			<p class="mt-2">
				Elci.io is not intended for users under 13 years of age. We do not knowingly collect
				information from children under 13. We prohibit sharing or processing data that belongs to
				or relates to minors, in compliance with platform policies.
			</p>
		</section>

		<section>
			<h2 class="text-base font-semibold text-(--color-text)">Contact Us</h2>
			<p class="mt-2">
				If you have questions about this Privacy Policy or wish to exercise your rights, please
				contact us:
			</p>
			<p class="mt-2">
				<strong class="text-(--color-text)">Email:</strong>
				<a href="mailto:support@elci.io" class="text-(--color-accent)">support@elci.io</a>
			</p>
		</section>

		<section>
			<h2 class="text-base font-semibold text-(--color-text)">Changes to This Privacy Policy</h2>
			<p class="mt-2">
				We may update this Privacy Policy from time to time. We will notify you of any changes by
				posting the new Privacy Policy on this page and updating the "Last Updated" date.
			</p>
			<p class="mt-3"><strong class="text-(--color-text)">This Privacy Policy is designed to comply with:</strong></p>
			<ul class="mt-2 list-disc space-y-1 pl-5">
				<li>General Data Protection Regulation (GDPR)</li>
				<li>California Consumer Privacy Act (CCPA)</li>
				<li>TikTok Developer Guidelines</li>
				<li>Meta Platform Terms</li>
				<li>X Developer Policy</li>
			</ul>
			<p class="mt-3"><strong class="text-(--color-text)">Version:</strong> 1.0</p>
		</section>
	</div>
</main>
<Footer />
