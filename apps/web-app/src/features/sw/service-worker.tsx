'use client'
import { Button } from '@thom/ui/button'
import { Heading } from '@thom/ui/heading'
import { Input } from '@thom/ui/input'
import { useIsomorphicEffect } from '_/components/hooks/use-isomorphic-effect'
import { ENVS } from '_/constants'
import { logger } from '_/lib/logger'
import { useState } from 'react'

// const VAPID = {
// 	subject: 'mailto: <inboxdothom@gmail.com>',
// 	publicKey: 'BIbmE_3OVGJzeipEKEkKq2_iF6Q9czyybOkMcXn58-qgxFYwBZIElIKpJaZYJZu5YHBBzM3dmxQj62KchFuYHdI',
// 	privateKey: 'WIyn86q5_XlKqm71V9RJCMR1uSVZ9uIFKBv44gEMTYM',
// }

const urlBase64ToUint8Array = (base64String: string) => {
	const padding = '='.repeat((4 - (base64String.length % 4)) % 4)
	const base64 = (base64String + padding).replace(/-/g, '+').replace(/_/g, '/')

	const rawData = window.atob(base64)
	const outputArray = new Uint8Array(rawData.length)

	for (let i = 0; i < rawData.length; ++i) {
		outputArray[i] = rawData.charCodeAt(i)
	}

	return outputArray
}

const swConfig: RegistrationOptions = { scope: '/', updateViaCache: 'none' }

const NOTIFY_ENDPOINT = '/api/notify'

const registerServiceWorker = async () => {
	try {
		return await navigator.serviceWorker.register('/sw.js', swConfig)
	} catch (error) {
		logger.error('Failed to register service worker', error)
		return null
	}
}

const getPushSubscription = async () => {
	try {
		const registration = await navigator.serviceWorker.ready
		return await registration.pushManager.getSubscription()
	} catch (error) {
		logger.error('Failed to read push subscription', error)
		return null
	}
}

const ServiceWorkerRegistrar = () => {
	useIsomorphicEffect(() => {
		if ('serviceWorker' in navigator) {
			registerServiceWorker()
		}
	}, [])

	return null
}

const INACTIVE = true

const PushNotificationManager = () => {
	const [isSupported, setIsSupported] = useState(false)
	const [subscription, setSubscription] = useState<PushSubscription | null>(null)
	const [message, setMessage] = useState('')

	useIsomorphicEffect(() => {
		if ('serviceWorker' in navigator && 'PushManager' in window) {
			;(async () => {
				const sub = await getPushSubscription()
				setSubscription(sub)
			})()
			setIsSupported(true)
		}
	}, [])

	const subscribeToPush = async () => {
		const permission = await Notification.requestPermission()
		if (permission !== 'granted') {
			logger.warn('Push permission not granted', permission)
			return
		}
		const registration = await navigator.serviceWorker.ready
		const sub = await registration.pushManager.subscribe({
			userVisibleOnly: true,
			applicationServerKey: urlBase64ToUint8Array(ENVS.NEXT_PUBLIC_VAPID_PUBLIC_KEY),
		})
		setSubscription(sub)
	}

	const unsubscribeFromPush = async () => {
		await subscription?.unsubscribe()
		setSubscription(null)
	}

	const sendTestNotification = async () => {
		if (subscription) {
			const sub = JSON.parse(JSON.stringify(subscription))
			const res = await fetch(NOTIFY_ENDPOINT, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ message, subscription: sub }),
			})
			if (!res.ok) {
				const errorText = await res.text()
				logger.error('Failed to send notification', { status: res.status, errorText })
			}
			setMessage('')
		}
	}

	if (INACTIVE) {
		return null
	}

	if (!isSupported) {
		return <Heading level={6}>Push notifications are not supported in this browser.</Heading>
	}

	return (
		<div className='rounded-lg bg-b-copy-dark/50 p-8'>
			<Heading>Push Notifications</Heading>
			{subscription ? (
				<>
					<Heading level={6} className='text-sm sm:font-light sm:text-sm'>
						You are subscribed to push notifications.
					</Heading>
					<div className='mt-4 flex items-center gap-6 max-sm:flex-wrap'>
						<Input
							type='text'
							placeholder='Enter notification message'
							value={message}
							onChange={e => setMessage(e.target.value)}
						/>
						<div className='flex shrink-0 gap-2'>
							<Button variant='outline' disabled={!message} onClick={sendTestNotification}>
								Send Test
							</Button>
							<Button variant='ghost' onClick={unsubscribeFromPush}>
								Unsubscribe
							</Button>
						</div>
					</div>
				</>
			) : (
				<div className='flex flex-wrap items-center gap-6'>
					<Heading level={6} className='text-sm sm:font-light sm:text-sm'>
						You are not subscribed to push notifications.
					</Heading>
					<Button variant='outline' onClick={subscribeToPush}>
						Subscribe
					</Button>
				</div>
			)}
		</div>
	)
}

export { PushNotificationManager, ServiceWorkerRegistrar }
