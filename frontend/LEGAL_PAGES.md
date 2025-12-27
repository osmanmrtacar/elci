# Legal Pages - Frontend Integration

## What Was Added

### 1. New Pages
- **`/privacy-policy`** - Complete Privacy Policy page
- **`/terms-of-service`** - Complete Terms of Service page

### 2. New Components
- **`PrivacyPolicyPage.tsx`** - Privacy Policy page component
- **`TermsOfServicePage.tsx`** - Terms of Service page component
- **`Footer.tsx`** - Footer component with legal links
- **`LegalPage.css`** - Shared styling for legal pages
- **`Footer.css`** - Footer styling

### 3. Updated Files
- **`App.tsx`** - Added routes for legal pages
- **`HomePage.tsx`** - Added footer and updated text to "Social Media Publisher"
- **`DashboardPage.tsx`** - Added footer

## Accessing the Pages

Visit these URLs in your app:
- Privacy Policy: `http://localhost:3000/privacy-policy`
- Terms of Service: `http://localhost:3000/terms-of-service`

## Footer Links

The footer appears on:
- Home page (before login)
- Dashboard page (after login)

The footer includes:
- Link to Privacy Policy
- Link to Terms of Service
- Copyright notice

## Customization Required

Before deploying to production, you MUST update these placeholders in both legal pages:

### Contact Information
```tsx
[your-email@domain.com] → support@yourdomain.com
[Your Business Address] → Your actual business address
[gdpr@yourdomain.com] → Your GDPR contact email
[ccpa@yourdomain.com] → Your CCPA contact email
[legal@yourdomain.com] → Your legal inquiries email
```

### Company Name
```tsx
© 2024 [Your Company Name] → © 2024 YourCompany Inc.
```

### How to Update

1. Open the files:
   - `src/pages/PrivacyPolicyPage.tsx`
   - `src/pages/TermsOfServicePage.tsx`

2. Find and replace all placeholders (marked with `[brackets]`)

3. Update the company name in `src/components/Footer.tsx`:
   ```tsx
   © {currentYear} Social Media Publisher. All rights reserved.
   ```

## Styling

The legal pages use a clean, professional design:
- White content container on gradient background
- Responsive layout (mobile-friendly)
- Easy-to-read typography
- Linked platform policies open in new tabs

### Customizing Colors

Edit `src/pages/LegalPage.css`:
```css
/* Main gradient background */
background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);

/* Link color */
color: #667eea;

/* Hover color */
color: #764ba2;
```

## Platform Developer Portal Setup

After customizing and deploying, add these URLs to your developer portals:

### TikTok
1. Go to https://developers.tiktok.com/
2. Navigate to your app settings
3. Add:
   - Privacy Policy URL: `https://yourdomain.com/privacy-policy`
   - Terms URL: `https://yourdomain.com/terms-of-service`

### Instagram (Meta)
1. Go to https://developers.facebook.com/
2. Select your app → Settings → Basic
3. Add:
   - Privacy Policy URL: `https://yourdomain.com/privacy-policy`
   - Terms of Service URL: `https://yourdomain.com/terms-of-service`

### X (Twitter)
1. Go to https://developer.x.com/
2. Navigate to your app
3. Add:
   - Privacy Policy URL: `https://yourdomain.com/privacy-policy`
   - Terms of Service URL: `https://yourdomain.com/terms-of-service`

## Testing

Before going live:

1. ✅ Test all links work correctly
2. ✅ Pages are accessible without authentication
3. ✅ All placeholders are replaced
4. ✅ Footer appears on all intended pages
5. ✅ Mobile responsive design works
6. ✅ External platform links open in new tabs

## Production Deployment

When deploying with Vercel (or any static host):

1. The pages will be automatically included in the build
2. Routes will work with React Router
3. Pages are accessible at:
   - `https://yourdomain.com/privacy-policy`
   - `https://yourdomain.com/terms-of-service`

No additional configuration needed!

## Legal Compliance Note

These pages are templates based on platform requirements. Before using in production:

- **Consult with a lawyer** to ensure compliance with your jurisdiction
- **Review platform requirements** regularly (they may change)
- **Keep documents updated** when you change data practices
- **Maintain version history** for compliance records

See the main `LEGAL_SETUP.md` in the root directory for more details.
