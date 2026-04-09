type Microtask = () => void

/**
 * Creates a microtask scheduler using Promises
 */
const createQueueMicrotaskViaPromises = (): ((microtask: Microtask) => void) => {
	return (microtask: Microtask): void => {
		Promise.resolve().then(() => {
			try {
				microtask()
			} catch (error) {
				// Re-throw in a new microtask to prevent unhandled rejections
				Promise.reject(error)
			}
		})
	}
}

/**
 * Creates a microtask scheduler using MutationObserver
 */
const createQueueMicrotaskViaMutationObserver = (): ((microtask: Microtask) => void) => {
	let counter = 0
	let microtaskQueue: Microtask[] = []

	const node = document.createTextNode('')

	const observer = new MutationObserver(() => {
		const tasks = [...microtaskQueue] // Create a copy to prevent mutation during execution
		microtaskQueue = [] // Clear queue first to allow new tasks to be queued

		for (const microtask of tasks) {
			try {
				microtask()
			} catch (error) {
				setTimeout(() => {
					throw error
				}, 0)
			}
		}
	})

	observer.observe(node, { characterData: true })

	return (microtask: Microtask): void => {
		microtaskQueue.push(microtask)
		node.data = String(++counter % 2)
	}
}

/**
 * Determines if the native Promise implementation is available
 */
const hasNativePromise = (): boolean => {
	return typeof Promise === 'function' && Promise.toString().indexOf('[native code]') > -1
}

/**
 * Creates the appropriate microtask implementation based on browser support
 */
const createQueueMicrotask = (): ((microtask: Microtask) => void) => {
	if (hasNativePromise()) {
		return createQueueMicrotaskViaPromises()
	}
	return createQueueMicrotaskViaMutationObserver()
}

const queueMicrotask = createQueueMicrotask()
export { queueMicrotask, type Microtask }
