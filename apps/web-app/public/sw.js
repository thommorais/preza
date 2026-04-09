const CACHE_NAME = 'app-cache-v2'
const FILES_TO_CACHE = []

self.addEventListener('install', event => {
	// biome-ignore lint/suspicious/noConsole: allowed here
	console.log('[Service Worker] Install')
	event.waitUntil(
		caches.open(CACHE_NAME).then(cache => {
			if (FILES_TO_CACHE.length === 0) {
				return undefined
			}
			// biome-ignore lint/suspicious/noConsole: allowed here
			console.log('[Service Worker] Pre-caching static assets')
			return cache.addAll(FILES_TO_CACHE)
		}),
	)
	self.skipWaiting()
})

self.addEventListener('fetch', event => {
	const request = event.request
	if (request.method !== 'GET') {
		return
	}
	const url = new URL(request.url)
	if (!url.protocol.startsWith('http')) {
		return
	}

	if (url.pathname.endsWith('.svg') || url.pathname.endsWith('.gif')) {
		event.respondWith(
			caches.match(request).then(cachedResponse => {
				if (cachedResponse) {
					// biome-ignore lint/suspicious/noConsole: allowed here
					console.log(`[Service Worker] Serving cached: ${request.url}`)
					return cachedResponse
				}
				// biome-ignore lint/suspicious/noConsole: allowed here
				console.log(`[Service Worker] Fetching and caching: ${request.url}`)
				return fetch(request).then(response => {
					if (!response || response.status !== 200) {
						return response
					}
					return caches.open(CACHE_NAME).then(cache => {
						cache.put(request, response.clone())
						return response
					})
				})
			}),
		)
	}
})

self.addEventListener('activate', event => {
	// biome-ignore lint/suspicious/noConsole: allowed here
	console.log('[Service Worker] Activate')
	event.waitUntil(
		caches.keys().then(cacheNames => {
			return Promise.all(
				cacheNames.map(cacheName => {
					if (cacheName !== CACHE_NAME) {
						// biome-ignore lint/suspicious/noConsole: allowed here
						console.log(`[Service Worker] Deleting old cache: ${cacheName}`)
						return caches.delete(cacheName)
					}
					return undefined
				}),
			)
		}),
	)
	self.clients.claim()
})

self.addEventListener('push', event => {
	event.waitUntil(
		(async () => {
			let data = { title: 'Notification', body: '' }
			if (event.data) {
				try {
					data = event.data.json()
				} catch {
					try {
						const text = await event.data.text()
						data = { title: 'Notification', body: text }
					} catch {
						// ignore malformed payloads
					}
				}
			}
			const options = {
				body: data.body ?? '',
				icon: '/icons/maskable_icon.png',
				badge: '/icons/icon.svg',
				vibrate: [100, 50, 100],
				data: {
					dateOfArrival: Date.now(),
					primaryKey: '2',
				},
			}
			await self.registration.showNotification(data.title || 'Notification', options)
		})(),
	)
})

self.addEventListener('notificationclick', event => {
	// biome-ignore lint/suspicious/noConsole: allowed here
	console.log('Notification click received.')
	event.notification.close()
	event.waitUntil(
		(async () => {
			const clientList = await clients.matchAll({ type: 'window', includeUncontrolled: true })
			const origin = self.location.origin
			const existingClient = clientList.find(client => client.url.startsWith(origin))
			if (existingClient && 'focus' in existingClient) {
				return existingClient.focus()
			}
			return clients.openWindow(origin)
		})(),
	)
})
