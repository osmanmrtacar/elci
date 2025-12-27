import './LegalPage.css'

const PrivacyPolicyPage = () => {
  return (
    <div className="legal-page">
      <div className="legal-container">
        <h1>Privacy Policy</h1>
        <p className="last-updated"><strong>Last Updated: December 27, 2025</strong></p>

        <section>
          <h2>Introduction</h2>
          <p>
            Welcome to our Social Media Publisher application ("Service", "we", "us", or "our").
            This Privacy Policy explains how we collect, use, disclose, and safeguard your information
            when you use our application that integrates with TikTok, Instagram (Meta), and X (formerly Twitter).
          </p>
          <p>
            By using our Service, you agree to the collection and use of information in accordance with this policy.
          </p>
        </section>

        <section>
          <h2>Information We Collect</h2>

          <h3>1. Information You Provide</h3>
          <p><strong>Account Information:</strong></p>
          <ul>
            <li>OAuth tokens from TikTok, Instagram, and X platforms</li>
            <li>Username and display name from connected platforms</li>
            <li>Profile information (avatar, bio) from connected platforms</li>
            <li>Email address (if provided by the platform)</li>
          </ul>

          <p><strong>Content Information:</strong></p>
          <ul>
            <li>Videos and media files you upload</li>
            <li>Captions and text content you create</li>
            <li>Post scheduling preferences</li>
            <li>Platform selection choices</li>
          </ul>

          <h3>2. Information We Collect Automatically</h3>
          <p><strong>Usage Data:</strong></p>
          <ul>
            <li>Log data (IP address, browser type, pages visited)</li>
            <li>Device information (device type, operating system)</li>
            <li>Post publishing history and status</li>
            <li>API interaction logs</li>
          </ul>

          <p><strong>Cookies and Tracking:</strong></p>
          <ul>
            <li>Session cookies for authentication</li>
            <li>Analytics cookies to improve our service</li>
            <li>OAuth state tokens for security</li>
          </ul>
        </section>

        <section>
          <h2>How We Use Your Information</h2>

          <h3>Platform Integration</h3>
          <ul>
            <li><strong>TikTok:</strong> To publish videos to your TikTok account as authorized by you</li>
            <li><strong>Instagram:</strong> To publish media to your Instagram Business or Creator account</li>
            <li><strong>X (Twitter):</strong> To publish posts and media to your X account</li>
          </ul>

          <h3>Service Operations</h3>
          <ul>
            <li>Authenticate your identity across platforms</li>
            <li>Process and publish your content to selected platforms</li>
            <li>Monitor post status and provide publishing confirmations</li>
            <li>Maintain and improve our Service</li>
            <li>Provide customer support</li>
            <li>Detect and prevent fraud or abuse</li>
          </ul>
        </section>

        <section>
          <h2>Data Sharing and Third Parties</h2>

          <h3>Third-Party Platforms</h3>
          <p>We share your information with the following third-party platforms only as necessary to provide our Service:</p>

          <p><strong>TikTok:</strong></p>
          <ul>
            <li>We share video content and captions you authorize us to publish</li>
            <li>User consent is obtained before posting with clear disclosure</li>
            <li>We comply with TikTok's Developer Guidelines</li>
          </ul>

          <p><strong>Instagram (Meta):</strong></p>
          <ul>
            <li>We share media and captions you authorize us to publish</li>
            <li>All data processing complies with Meta's Platform Terms</li>
            <li>We do not share data of minors or sensitive personal data</li>
          </ul>

          <p><strong>X (Twitter):</strong></p>
          <ul>
            <li>We share posts and media you authorize us to publish</li>
            <li>All usage respects users' reasonable expectations of privacy</li>
            <li>We do not use data for surveillance, tracking, or monitoring purposes</li>
          </ul>
        </section>

        <section>
          <h2>Your Rights and Choices</h2>

          <h3>Data Access and Portability</h3>
          <p>You have the right to:</p>
          <ul>
            <li>Access the personal information we hold about you</li>
            <li>Request a copy of your data in a portable format</li>
            <li>Review your posting history and connected platforms</li>
          </ul>

          <h3>Data Deletion</h3>
          <p>You can request deletion of your data by:</p>
          <ul>
            <li>Disconnecting platforms from your account settings</li>
            <li>Contacting us at [your-email@domain.com]</li>
            <li>We will delete your data within 30 days of your request</li>
          </ul>

          <h3>GDPR Rights (EU Users)</h3>
          <ul>
            <li>Right to rectification of inaccurate data</li>
            <li>Right to erasure ("right to be forgotten")</li>
            <li>Right to restrict processing</li>
            <li>Right to data portability</li>
            <li>Right to object to processing</li>
            <li>Right to withdraw consent</li>
          </ul>

          <h3>CCPA Rights (California Users)</h3>
          <ul>
            <li>Know what personal information is collected</li>
            <li>Know whether personal information is sold or disclosed</li>
            <li>Say no to the sale of personal information</li>
            <li>Access your personal information</li>
            <li>Request deletion of personal information</li>
            <li>Equal service and price, even if you exercise your privacy rights</li>
          </ul>
        </section>

        <section>
          <h2>Data Security</h2>
          <p>We implement appropriate technical and organizational measures to protect your information:</p>
          <ul>
            <li>Encryption of data in transit (HTTPS/TLS)</li>
            <li>Encrypted storage of sensitive data</li>
            <li>Access controls and authentication</li>
            <li>Regular security audits and updates</li>
          </ul>
          <p>
            However, no method of transmission over the Internet is 100% secure. While we strive to protect
            your information, we cannot guarantee its absolute security.
          </p>
        </section>

        <section>
          <h2>Platform-Specific Policies</h2>
          <p>Your use of our Service is also governed by the privacy policies of:</p>
          <ul>
            <li><a href="https://www.tiktok.com/legal/page/eea/privacy-policy/en" target="_blank" rel="noopener noreferrer">TikTok Privacy Policy</a></li>
            <li><a href="https://www.facebook.com/privacy/policy/" target="_blank" rel="noopener noreferrer">Meta Privacy Policy</a> (for Instagram)</li>
            <li><a href="https://twitter.com/en/privacy" target="_blank" rel="noopener noreferrer">X Privacy Policy</a></li>
          </ul>
        </section>

        <section>
          <h2>Children's Privacy</h2>
          <p>
            Our Service is not intended for users under 13 years of age. We do not knowingly collect information
            from children under 13. We prohibit sharing or processing data that belongs to or relates to minors,
            in compliance with platform policies.
          </p>
        </section>

        <section>
          <h2>Contact Us</h2>
          <p>If you have questions about this Privacy Policy or wish to exercise your rights, please contact us:</p>
          <ul>
            <li><strong>Email:</strong> [your-email@domain.com]</li>
            <li><strong>Address:</strong> [Your Business Address]</li>
            <li><strong>GDPR Contact:</strong> [gdpr@yourdomain.com]</li>
            <li><strong>CCPA Contact:</strong> [ccpa@yourdomain.com]</li>
          </ul>
        </section>

        <section>
          <h2>Changes to This Privacy Policy</h2>
          <p>
            We may update this Privacy Policy from time to time. We will notify you of any changes by posting
            the new Privacy Policy on this page and updating the "Last Updated" date.
          </p>
        </section>

        <div className="compliance-note">
          <p><strong>This Privacy Policy is designed to comply with:</strong></p>
          <ul>
            <li>General Data Protection Regulation (GDPR)</li>
            <li>California Consumer Privacy Act (CCPA)</li>
            <li>TikTok Developer Guidelines</li>
            <li>Meta Platform Terms</li>
            <li>X Developer Policy</li>
          </ul>
          <p><strong>Version:</strong> 1.0</p>
        </div>
      </div>
    </div>
  )
}

export default PrivacyPolicyPage
