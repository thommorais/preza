import { ENVS } from '_/constants'
import { cIC, rIC } from '_/lib/idle-callback-polyfills'
import { now } from '_/lib/now'
import { queueMicrotask } from '_/lib/queue-microtask'

// Constants
const DEFAULT_MIN_TASK_TIME = 0

const isSafari = ENVS.IS_CLIENT && /^((?!chrome|android).)*safari/i.test(window.navigator.userAgent)

// Types
type TaskState = {
	time: number
	visibilityState: DocumentVisibilityState
}

type QueuedTask = {
	state: TaskState
	task: (state: TaskState) => void
	minTaskTime: number
}

type IdleQueueOptions = {
	ensureTasksRun?: boolean
	defaultMinTaskTime?: number
}

type IdleQueueState = {
	taskQueue: QueuedTask[]
	idleCallbackHandle: number | null
	isProcessing: boolean
	currentState: TaskState | null
}

// Helper Functions
const shouldYield = (deadline: IdleDeadline | undefined, minTaskTime: number): boolean => {
	return !!(deadline && deadline.timeRemaining() <= minTaskTime)
}

const createTaskState = (): TaskState => ({
	time: now(),
	visibilityState: document.visibilityState,
})

const createInitialState = (): IdleQueueState => ({
	taskQueue: [],
	idleCallbackHandle: null,
	isProcessing: false,
	currentState: null,
})

// Main Queue Operations
const createIdleQueue = (options: IdleQueueOptions = {}) => {
	const { ensureTasksRun = false, defaultMinTaskTime = DEFAULT_MIN_TASK_TIME } = options

	let state = createInitialState()

	const cancelScheduledRun = () => {
		if (state.idleCallbackHandle !== null) {
			cIC(state.idleCallbackHandle)
			state.idleCallbackHandle = null
		}
	}

	const scheduleTasksToRun = () => {
		if (ensureTasksRun && document.visibilityState === 'hidden') {
			queueMicrotask(runTasks)
		} else if (!state.idleCallbackHandle) {
			state.idleCallbackHandle = rIC(runTasks)
		}
	}

	const addTask = (task: (state: TaskState) => void, options: { minTaskTime?: number } = {}) => {
		const { minTaskTime = defaultMinTaskTime } = options
		const taskState = createTaskState()

		state.taskQueue.push({
			state: taskState,
			task,
			minTaskTime,
		})

		scheduleTasksToRun()
	}

	const runTasks = (deadline?: IdleDeadline) => {
		cancelScheduledRun()

		if (!state.isProcessing) {
			state.isProcessing = true

			while (state.taskQueue.length > 0 && !shouldYield(deadline, state.taskQueue[0]?.minTaskTime ?? 0)) {
				const nextTask = state.taskQueue.shift()
				if (nextTask) {
					state.currentState = nextTask.state
					nextTask.task(nextTask.state)
					state.currentState = null
				}
			}

			state.isProcessing = false

			if (state.taskQueue.length > 0) {
				scheduleTasksToRun()
			}
		}
	}

	const runTasksImmediately = () => {
		runTasks()
	}

	const onVisibilityChange = () => {
		if (document.visibilityState === 'hidden') {
			runTasksImmediately()
		}
	}

	// Setup event listeners
	if (ensureTasksRun) {
		addEventListener('visibilitychange', onVisibilityChange, true)
		if (isSafari) {
			addEventListener('beforeunload', runTasksImmediately, true)
		}
	}

	// Public API
	return {
		pushTask: (task: (state: TaskState) => void, options?: { minTaskTime?: number }) => {
			addTask(task, options)
		},

		unshiftTask: (task: (state: TaskState) => void, options?: { minTaskTime?: number }) => {
			const { minTaskTime = defaultMinTaskTime } = options || {}
			const taskState = createTaskState()

			state.taskQueue.unshift({
				state: taskState,
				task,
				minTaskTime,
			})

			scheduleTasksToRun()
		},

		runTasksImmediately,

		hasPendingTasks: () => state.taskQueue.length > 0,

		clearPendingTasks: () => {
			state.taskQueue = []
			cancelScheduledRun()
		},

		getState: () => state.currentState,

		destroy: () => {
			state.taskQueue = []
			cancelScheduledRun()

			if (ensureTasksRun) {
				removeEventListener('visibilitychange', onVisibilityChange, true)
				if (isSafari) {
					removeEventListener('beforeunload', runTasksImmediately, true)
				}
			}

			state = createInitialState()
		},
	}
}

export type IdleQueue = ReturnType<typeof createIdleQueue>
export { createIdleQueue }
