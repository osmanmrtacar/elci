# Legal Documents Setup Guide

## Overview

This guide explains how to set up and customize the Privacy Policy and Terms of Service for your social media publishing application.

## Documents Created

1. **PRIVACY_POLICY.md** - Comprehensive privacy policy compliant with:
   - GDPR (EU)
   - CCPA (California)
   - TikTok Developer Guidelines
   - Meta Platform Terms
   - X Developer Policy

2. **TERMS_OF_SERVICE.md** - Complete terms of service covering:
   - User responsibilities
   - Platform integration requirements
   - Content publishing rules
   - Liability limitations
   - Compliance with all three platforms

## Required Customizations

Before deploying, you MUST update the following placeholders:

### In Both Documents

1. **Contact Information**
   ```markdown
   Email: [your-email@domain.com]
   Address: [Your Business Address]
   Company Name: [Your Company Name]
   ```

2. **Domain URLs**
   ```markdown
   https://yourdomain.com
   ```

3. **Data Protection Officers** (if applicable)
   ```markdown
   DPO Email: [dpo@yourdomain.com]
   GDPR Contact: [gdpr@yourdomain.com]
   CCPA Contact: [ccpa@yourdomain.com]
   ```

4. **Jurisdiction** (Terms of Service)
   ```markdown
   Governed by the laws of [Your Jurisdiction]
   ```

## Hosting Options

### Option 1: GitHub Pages (Free)

1. Create a new repository or use existing one
2. Add privacy-policy.html and terms.html to repo
3. Enable GitHub Pages in repository settings
4. Access at: `https://yourusername.github.io/repo-name/privacy-policy.html`

### Option 2: Your Main Domain

Host on your application domain:
- `https://yourdomain.com/privacy-policy`
- `https://yourdomain.com/terms-of-service`

**Implementation:**
1. Convert markdown to HTML (see below)
2. Add routes in your application
3. Or use static hosting (nginx, Vercel, etc.)

### Option 3: Vercel/Netlify (Recommended)

1. Create a simple static site with the documents
2. Deploy to Vercel or Netlify
3. Use custom domain or free subdomain

## Converting to HTML

### Using Pandoc (Recommended)

```bash
# Install pandoc
brew install pandoc  # macOS
# or
sudo apt install pandoc  # Linux

# Convert to HTML
pandoc PRIVACY_POLICY.md -o privacy-policy.html -s --metadata title="Privacy Policy"
pandoc TERMS_OF_SERVICE.md -o terms.html -s --metadata title="Terms of Service"

# With custom CSS
pandoc PRIVACY_POLICY.md -o privacy-policy.html -s -c style.css
```

### Simple HTML Template

Create files in `frontend/public/` for static hosting:

**public/privacy-policy.html**
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Privacy Policy - Your App Name</title>
    <style>
        body {
            max-width: 800px;
            margin: 40px auto;
            padding: 0 20px;
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
            line-height: 1.6;
        }
        h1 { border-bottom: 2px solid #333; padding-bottom: 10px; }
        h2 { margin-top: 30px; color: #333; }
        a { color: #0066cc; }
    </style>
</head>
<body>
    <!-- Paste converted HTML content here -->
</body>
</html>
```

## Platform Requirements

### TikTok Developer Portal

1. Go to [TikTok Developer Portal](https://developers.tiktok.com/)
2. Navigate to your app settings
3. Add Privacy Policy URL: `https://yourdomain.com/privacy-policy`
4. Add Terms of Service URL: `https://yourdomain.com/terms`

### Meta for Developers (Instagram)

1. Go to [Meta for Developers](https://developers.facebook.com/)
2. Select your app
3. Settings → Basic
4. Add Privacy Policy URL: `https://yourdomain.com/privacy-policy`
5. Add Terms of Service URL: `https://yourdomain.com/terms`

### X Developer Portal

1. Go to [X Developer Portal](https://developer.x.com/)
2. Navigate to your app
3. Add Privacy Policy URL in app settings
4. Add Terms of Service URL

## App Store Requirements (If Applicable)

### Apple App Store

- Privacy Policy URL is **required**
- Must be accessible without authentication
- Must be in the same language as the app

### Google Play Store

- Privacy Policy URL is **required**
- Must be hosted on a publicly accessible URL
- Must be up-to-date with current practices

## Compliance Checklist

Before going live, ensure:

- [ ] All placeholder text is replaced with your information
- [ ] Contact email addresses are valid and monitored
- [ ] URLs are publicly accessible (test in incognito mode)
- [ ] Documents are linked in your app settings
- [ ] Documents are linked in platform developer portals
- [ ] "Last Updated" date is current
- [ ] Documents are linked in your app footer/settings
- [ ] GDPR and CCPA contact methods are functional
- [ ] You can actually fulfill data deletion requests
- [ ] You have a process for handling privacy inquiries

## Linking in Your Application

### Frontend Footer

Add links to your React app footer:

```typescript
// src/components/Footer.tsx
export const Footer = () => (
  <footer>
    <a href="/privacy-policy" target="_blank">Privacy Policy</a>
    <span> | </span>
    <a href="/terms-of-service" target="_blank">Terms of Service</a>
  </footer>
)
```

### During OAuth Flow

Show links before platform connection:

```typescript
<div className="legal-notice">
  By connecting your account, you agree to our{' '}
  <a href="/terms" target="_blank">Terms of Service</a> and{' '}
  <a href="/privacy" target="_blank">Privacy Policy</a>
</div>
```

### Before Publishing

Show confirmation with links:

```typescript
<Checkbox>
  I agree to the{' '}
  <a href="/terms" target="_blank">Terms of Service</a>
  {' '}and confirm I have rights to publish this content
</Checkbox>
```

## Updating Documents

### When to Update

Update documents when you:
- Add new platforms or features
- Change data collection practices
- Modify data retention policies
- Change service providers
- Receive legal advice to do so

### How to Update

1. Update the markdown files
2. Update "Last Updated" date
3. Increment version number
4. Convert to HTML and redeploy
5. Notify users of material changes (email)
6. Keep old versions for compliance (archive)

## Legal Disclaimer

These documents are provided as templates and may not be suitable for all situations. You should:

- **Consult a lawyer** specializing in internet/privacy law
- **Review with your legal team** before deployment
- **Customize** based on your specific jurisdiction and business model
- **Stay updated** on changing laws (GDPR, CCPA, etc.)
- **Monitor** platform policy changes

We provide these templates to help you get started, but they do not constitute legal advice.

## Resources

### Platform Documentation

- [TikTok Developer Guidelines](https://developers.tiktok.com/doc/our-guidelines-developer-guidelines)
- [Meta Platform Terms](https://developers.facebook.com/terms)
- [X Developer Agreement](https://developer.x.com/en/developer-terms/agreement-and-policy)

### Privacy Law Resources

- [GDPR Official Text](https://gdpr-info.eu/)
- [CCPA Official Website](https://oag.ca.gov/privacy/ccpa)
- [Privacy Policy Generators](https://www.termsfeed.com/)

### Tools

- [Pandoc Documentation](https://pandoc.org/MANUAL.html)
- [Markdown to HTML Converters](https://markdowntohtml.com/)
- [Privacy Policy Templates](https://www.freeprivacypolicy.com/)

## Support

If you need help customizing these documents:
1. Consult with a qualified attorney
2. Use professional legal document services
3. Contact platform support for specific requirements

---

**Remember**: These are living documents that should be reviewed and updated regularly to reflect changes in your service, applicable laws, and platform requirements.
