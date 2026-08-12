export const prerender = true;

const siteUrl = 'https://elci.io';
const lastmod = new Date().toISOString().split('T')[0];

export function GET(): Response {
	const urls = [
		{ loc: `${siteUrl}/`, changefreq: 'weekly', priority: '1.0' },
		{ loc: `${siteUrl}/privacy-policy`, changefreq: 'monthly', priority: '0.3' },
		{ loc: `${siteUrl}/terms-of-service`, changefreq: 'monthly', priority: '0.3' },
		{ loc: `${siteUrl}/login`, changefreq: 'monthly', priority: '0.5' }
	];

	const xml = `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
${urls.map((u) => `  <url><loc>${u.loc}</loc><lastmod>${lastmod}</lastmod><changefreq>${u.changefreq}</changefreq><priority>${u.priority}</priority></url>`).join('\n')}
</urlset>`;

	return new Response(xml, {
		headers: { 'Content-Type': 'application/xml' }
	});
}
