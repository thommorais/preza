import { LOCAL_HREFS } from '_/constants'

async function sitemap() {
	return [LOCAL_HREFS].map(route => ({
		url: `${route}`,
		lastModified: new Date().toISOString().split('T')[0],
	}))
}

// biome-ignore lint/style/noDefaultExport: page
export default sitemap
