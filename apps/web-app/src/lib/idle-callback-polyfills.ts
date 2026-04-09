import { now } from '_/lib/now'

// Type definitions
interface IdleDeadline {
	didTimeout: boolean
	timeRemaining: () => number
}

type IdleRequestCallback = (deadline: IdleDeadline) => void

interface IdleRequestOptions {
	timeout?: number
}

const supportsRequestIdleCallback = typeof requestIdleCallback === 'function'

const createIdleDeadline = (initTime: number): IdleDeadline => ({
	didTimeout: false,
	timeRemaining: () => Math.max(0, 50 - (now() - initTime)),
})

const requestIdleCallbackShim = (callback: IdleRequestCallback, _options?: IdleRequestOptions): number => {
	const deadline = createIdleDeadline(now())
	return setTimeout(() => callback(deadline), 0) as unknown as number
}

const cancelIdleCallbackShim = (handle: number): void => {
	clearTimeout(handle)
}

const rIC = (callback: IdleRequestCallback, options?: IdleRequestOptions): number => {
	return supportsRequestIdleCallback
		? window.requestIdleCallback(callback, options)
		: requestIdleCallbackShim(callback, options)
}

const cIC = (handle: number): void => {
	if (supportsRequestIdleCallback) {
		window.cancelIdleCallback(handle as number)
	} else {
		cancelIdleCallbackShim(handle)
	}
}

export { cIC, rIC }
