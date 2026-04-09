type LogLevel = 'error' | 'warn' | 'info' | 'log' | 'debug'
type LogMethod = (...args: unknown[]) => void

type Logger = {
	error: LogMethod
	warn: LogMethod
	info: LogMethod
	log: LogMethod
}

const createLogger = (): Logger => {
	const shouldLog = (level: LogLevel) => {
		if (process.env.NODE_ENV === 'production') {
			return level === 'error'
		}
		return true
	}

	const error: LogMethod = (...args) => {
		if (shouldLog('error')) {
			// biome-ignore lint/suspicious/noConsole: Centralized logging utility - all console usage goes through here
			console.error(...args)
		}
	}

	const warn: LogMethod = (...args) => {
		if (shouldLog('warn')) {
			// biome-ignore lint/suspicious/noConsole: Centralized logging utility - all console usage goes through here
			console.warn(...args)
		}
	}

	const info: LogMethod = (...args) => {
		if (shouldLog('info')) {
			// biome-ignore lint/suspicious/noConsole: Centralized logging utility - all console usage goes through here
			console.info(...args)
		}
	}

	const log: LogMethod = (...args) => {
		if (shouldLog('log')) {
			// biome-ignore lint/suspicious/noConsole: Centralized logging utility - all console usage goes through here
			console.log(...args)
		}
	}

	return {
		error,
		warn,
		info,
		log,
	}
}

export const logger = createLogger()
