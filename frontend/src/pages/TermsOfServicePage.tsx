import './LegalPage.css'

const TermsOfServicePage = () => {
  return (
    <div className="legal-page">
      <div className="legal-container">
        <h1>Terms of Service</h1>
        <p className="last-updated"><strong>Last Updated: December 27, 2025</strong></p>

        <section>
          <h2>1. Agreement to Terms</h2>
          <p>
            By accessing or using our Social Media Publisher application ("Service", "we", "us", or "our"),
            you agree to be bound by these Terms of Service ("Terms"). If you disagree with any part of these
            Terms, you may not access the Service.
          </p>
        </section>

        <section>
          <h2>2. Description of Service</h2>
          <p>Our Service enables you to:</p>
          <ul>
            <li>Connect your social media accounts (TikTok, Instagram, X/Twitter)</li>
            <li>Upload and publish content to multiple platforms simultaneously</li>
            <li>Manage and track your posts across platforms</li>
            <li>Schedule and organize your social media content</li>
          </ul>
        </section>

        <section>
          <h2>3. Eligibility</h2>

          <h3>Age Requirements</h3>
          <p>You must be at least:</p>
          <ul>
            <li><strong>13 years old</strong> to use this Service</li>
            <li><strong>18 years old</strong> (or the age of majority in your jurisdiction) to use the Service without parental consent</li>
            <li>Of legal age to enter into a binding contract in your jurisdiction</li>
          </ul>

          <h3>Account Requirements</h3>
          <ul>
            <li>Have valid accounts on the platforms you wish to connect (TikTok, Instagram, X)</li>
            <li>For Instagram: Have a Business or Creator account</li>
            <li>Comply with each platform's Terms of Service</li>
            <li>Provide accurate and complete registration information</li>
            <li>Maintain the security of your account credentials</li>
          </ul>
        </section>

        <section>
          <h2>4. Platform Integration and Authorization</h2>

          <h3>OAuth Authorization</h3>
          <p>When you connect a platform:</p>
          <ul>
            <li>You authorize us to access your account according to platform permissions</li>
            <li>You grant us permission to publish content on your behalf</li>
            <li>You can revoke authorization at any time through platform settings</li>
            <li>We only access data necessary to provide the Service</li>
          </ul>

          <h3>Platform-Specific Requirements</h3>

          <p><strong>TikTok:</strong></p>
          <ul>
            <li>You acknowledge and agree to <a href="https://www.tiktok.com/legal/page/global/terms-of-service/en" target="_blank" rel="noopener noreferrer">TikTok's Terms of Service</a></li>
            <li>By posting through our Service, you agree to TikTok's Music Usage Confirmation</li>
            <li>You must have full awareness and control of what is being posted to your account</li>
            <li>You provide explicit consent before each post is published</li>
          </ul>

          <p><strong>Instagram (Meta):</strong></p>
          <ul>
            <li>You acknowledge and agree to <a href="https://www.facebook.com/legal/terms" target="_blank" rel="noopener noreferrer">Meta's Platform Terms</a></li>
            <li>Your Instagram account must be a Business or Creator account</li>
            <li>You comply with Instagram's content and community guidelines</li>
            <li>You understand data processing is governed by Meta's Privacy Policy</li>
          </ul>

          <p><strong>X (Twitter):</strong></p>
          <ul>
            <li>You acknowledge and agree to <a href="https://twitter.com/en/tos" target="_blank" rel="noopener noreferrer">X's Terms of Service</a></li>
            <li>You comply with X's <a href="https://developer.x.com/en/developer-terms/agreement-and-policy" target="_blank" rel="noopener noreferrer">Developer Agreement and Policy</a></li>
            <li>You understand your content distribution is subject to X's policies</li>
            <li>You will not use the Service for surveillance or tracking purposes</li>
          </ul>
        </section>

        <section>
          <h2>5. User Responsibilities</h2>

          <h3>Content Ownership and Rights</h3>
          <p>You represent and warrant that:</p>
          <ul>
            <li>You own or have the right to use all content you upload</li>
            <li>Your content does not infringe on intellectual property rights of others</li>
            <li>You have obtained all necessary licenses, rights, and permissions for music, images, or other third-party content</li>
            <li>You have rights to publish content to all selected platforms</li>
          </ul>

          <h3>Prohibited Content</h3>
          <p>You agree not to upload, post, or transmit content that:</p>
          <ul>
            <li>Violates any law or regulation</li>
            <li>Infringes on intellectual property rights</li>
            <li>Contains malware, viruses, or harmful code</li>
            <li>Is defamatory, obscene, or offensive</li>
            <li>Promotes violence, discrimination, or hate speech</li>
            <li>Involves minors inappropriately</li>
            <li>Contains sensitive personal data without consent</li>
            <li>Violates platform-specific community guidelines</li>
          </ul>

          <h3>Prohibited Uses</h3>
          <p>You agree not to:</p>
          <ul>
            <li>Use the Service for surveillance, tracking, or monitoring of users</li>
            <li>Investigate or track users or their content</li>
            <li>Monitor sensitive events (protests, rallies, community organizing)</li>
            <li>Scrape or harvest data from platforms</li>
            <li>Reverse engineer or attempt to access source code</li>
            <li>Bypass rate limits or security measures</li>
            <li>Resell or redistribute platform data</li>
            <li>Use the Service to spam or send unsolicited messages</li>
            <li>Impersonate others or misrepresent your affiliation</li>
            <li>Interfere with or disrupt the Service or servers</li>
          </ul>
        </section>

        <section>
          <h2>6. Content Publishing</h2>

          <h3>Publishing Process</h3>
          <p>When you publish content through our Service:</p>
          <ol>
            <li>You select platforms for publication</li>
            <li>You review and confirm content and settings</li>
            <li>You provide explicit consent by clicking "Publish"</li>
            <li>We process and publish to selected platforms on your behalf</li>
          </ol>

          <h3>Platform Posting Confirmation</h3>
          <p>Before publishing, you confirm:</p>
          <ul>
            <li><strong>TikTok:</strong> "By posting, you agree to TikTok's Music Usage Confirmation"</li>
            <li><strong>Instagram:</strong> "I have rights to publish this content to Instagram"</li>
            <li><strong>X:</strong> "I agree to X's content distribution terms"</li>
          </ul>
        </section>

        <section>
          <h2>7. Disclaimers and Limitations of Liability</h2>

          <h3>Service "As Is"</h3>
          <p className="uppercase-text">
            THE SERVICE IS PROVIDED "AS IS" AND "AS AVAILABLE" WITHOUT WARRANTIES OF ANY KIND, EXPRESS OR IMPLIED.
          </p>

          <h3>Limitation of Liability</h3>
          <p className="uppercase-text">
            TO THE MAXIMUM EXTENT PERMITTED BY LAW, WE SHALL NOT BE LIABLE FOR ANY INDIRECT, INCIDENTAL,
            SPECIAL, OR CONSEQUENTIAL DAMAGES, LOSS OF PROFITS, DATA, OR GOODWILL.
          </p>
          <p>
            OUR TOTAL LIABILITY SHALL NOT EXCEED THE AMOUNT YOU PAID US IN THE PAST 12 MONTHS, OR $100, WHICHEVER IS GREATER.
          </p>
        </section>

        <section>
          <h2>8. Privacy and Data Protection</h2>
          <p>
            Your privacy is important to us. Our collection and use of personal information is described in our{' '}
            <a href="/privacy-policy">Privacy Policy</a>.
          </p>
          <p>Key points:</p>
          <ul>
            <li>We collect only necessary data to provide the Service</li>
            <li>We comply with GDPR, CCPA, and other privacy laws</li>
            <li>You can request data deletion at any time</li>
            <li>We do not sell your personal information</li>
            <li>We use encryption and security measures to protect your data</li>
          </ul>
        </section>

        <section>
          <h2>9. Termination</h2>
          <p>You may terminate your account at any time by disconnecting all platforms or contacting us.</p>
          <p>We may suspend or terminate your access if you violate these Terms or platform policies.</p>
        </section>

        <section>
          <h2>10. Changes to Terms</h2>
          <p>
            We may update these Terms from time to time. Continued use after changes constitutes acceptance of new Terms.
          </p>
        </section>

        <section>
          <h2>11. Platform Compliance</h2>
          <p>We comply with developer requirements for:</p>

          <p><strong>TikTok:</strong></p>
          <ul>
            <li><a href="https://developers.tiktok.com/doc/our-guidelines-developer-guidelines" target="_blank" rel="noopener noreferrer">Developer Guidelines</a></li>
            <li><a href="https://developers.tiktok.com/doc/content-sharing-guidelines" target="_blank" rel="noopener noreferrer">Content Posting API Guidelines</a></li>
          </ul>

          <p><strong>Instagram (Meta):</strong></p>
          <ul>
            <li><a href="https://developers.facebook.com/terms" target="_blank" rel="noopener noreferrer">Meta Platform Terms</a></li>
            <li>Instagram API requirements</li>
          </ul>

          <p><strong>X (Twitter):</strong></p>
          <ul>
            <li><a href="https://developer.x.com/en/developer-terms/agreement-and-policy" target="_blank" rel="noopener noreferrer">Developer Agreement and Policy</a></li>
            <li><a href="https://developer.twitter.com/en/developer-terms/more-on-restricted-use-cases" target="_blank" rel="noopener noreferrer">Restricted Use Cases</a></li>
          </ul>
        </section>

        <section>
          <h2>12. Contact Information</h2>
          <p>For questions about these Terms, please contact us:</p>
          <ul>
            <li><strong>Email:</strong> [your-email@domain.com]</li>
            <li><strong>Address:</strong> [Your Business Address]</li>
            <li><strong>Legal Inquiries:</strong> [legal@yourdomain.com]</li>
          </ul>
        </section>

        <section>
          <h2>13. Acknowledgment</h2>
          <p>By using our Service, you acknowledge that:</p>
          <ul>
            <li>You have read and understood these Terms of Service</li>
            <li>You have read and understood our Privacy Policy</li>
            <li>You agree to comply with all platform Terms of Service</li>
            <li>You have the authority to enter into this agreement</li>
            <li>You will use the Service in compliance with all applicable laws</li>
          </ul>
        </section>

        <div className="compliance-note">
          <p><strong>Effective Date:</strong> December 27, 2024</p>
          <p><strong>Version:</strong> 1.0</p>
        </div>
      </div>
    </div>
  )
}

export default TermsOfServicePage
