'use client'

import { animate, spring, utils } from 'animejs'
import 'flag-icons/css/flag-icons.min.css'
import { useI18n } from '_/i18n/config-client'
import { getStoredScore, getStoredSession, persistScore, persistSession } from '_/lib/db/indexedDb'
import { COUNTRY_FLAGS } from '_/lib/flags/data'
import {
	clampRecentHistory,
	computeAccuracy,
	createRound,
	EMPTY_SCORE,
	EMPTY_SESSION,
	type FlagRound,
	type GameSession,
	type RoundResult,
	type ScoreSnapshot,
	summarizeRound,
} from '_/lib/flags/game'
import { tv } from '_/lib/third-party/tv'
import { type ComponentPropsWithoutRef, useCallback, useEffect, useId, useMemo, useRef, useState } from 'react'

const KEY_DEBOUNCE_MS = 140

type SelectionState = {
	choice: string | null
	result: RoundResult | 'idle'
}

const formatDate = (value: string | null): string => {
	if (!value) {
		return '—'
	}

	const parsed = new Date(value)
	if (Number.isNaN(parsed.getTime())) {
		return '—'
	}

	return parsed.toLocaleString()
}

const useStableLatest = <T,>(value: T) => {
	const ref = useRef(value)
	useEffect(() => {
		ref.current = value
	}, [value])
	return ref
}

const clsxClasses = tv({ base: [] })

export const clsx = (...args: (string | undefined)[]) => clsxClasses({ class: [args] })

const Card = ({ className, ...props }: ComponentPropsWithoutRef<'div'>) => {
	return (
		<div
			className={clsx('w-full rounded-xl border border-info-800/50 bg-info-800/40 p-8 shadow-2xl', className)}
			{...props}
		/>
	)
}

export const FlagsGameClient = () => {
	const roundsInputId = useId()
	const t = useI18n()

	const MAX_ROUNDS = useMemo(() => Math.ceil(COUNTRY_FLAGS.length / 5) * 5, [])

	const [round, setRound] = useState<FlagRound | null>(null)
	const [score, setScore] = useState<ScoreSnapshot>(EMPTY_SCORE)
	const [session, setSession] = useState<GameSession>(EMPTY_SESSION)
	const [selectedRounds, setSelectedRounds] = useState(10)
	const [selection, setSelection] = useState<SelectionState>({ choice: null, result: 'idle' })
	const [loadingScore, setLoadingScore] = useState<boolean>(true)
	const [persistenceWarning, setPersistenceWarning] = useState<string | null>(null)
	const [flagError, setFlagError] = useState(false)
	const [recentCorrect, setRecentCorrect] = useState<string[]>([])
	const [status, setStatus] = useState<'playing' | 'revealed'>('playing')
	const [focusedIndex, setFocusedIndex] = useState(0)

	const feedbackVariants: Record<RoundResult | 'idle', string> = useMemo(
		() => ({
			idle: t('take_best_guess'),
			correct: t('correct_nice_work'),
			incorrect: t('not_quite'),
			skipped: t('round_skipped'),
		}),
		[t],
	)

	const optionRefs = useRef<Array<HTMLButtonElement | null>>([])
	const nextButtonRef = useRef<HTMLButtonElement | null>(null)
	const lastKeyTimestampRef = useRef<number>(0)
	const interactionLockRef = useRef(false)
	const flagCardRef = useRef<HTMLDivElement | null>(null)

	const roundRef = useStableLatest(round)
	const scoreRef = useStableLatest(score)
	const recentRef = useStableLatest(recentCorrect)
	const statusRef = useStableLatest(status)

	const accuracy = useMemo(() => computeAccuracy(score), [score])

	const persistScoreSafely = useCallback(async (nextScore: ScoreSnapshot) => {
		try {
			await persistScore(nextScore)
			setPersistenceWarning(null)
		} catch (_error) {
			setPersistenceWarning('Unable to persist score. Using in-memory backup until storage is available.')
		}
	}, [])

	const persistSessionSafely = useCallback(async (nextSession: GameSession) => {
		try {
			await persistSession(nextSession)
		} catch (_error) {}
	}, [])

	useEffect(() => {
		setRound(createRound(COUNTRY_FLAGS))
	}, [])

	const loadInitialData = useCallback(
		async (cancelled: { value: boolean }) => {
			try {
				const [storedScore, storedSession] = await Promise.all([getStoredScore(), getStoredSession()])

				if (cancelled.value) return

				if (storedScore) {
					setScore(storedScore)
					scoreRef.current = storedScore
				}
				if (storedSession) {
					setSession(storedSession)
				}
			} catch (_error) {
				if (!cancelled.value) {
					setPersistenceWarning('IndexedDB is unavailable. Progress is temporarily stored in memory.')
				}
			} finally {
				if (!cancelled.value) {
					setLoadingScore(false)
				}
			}
		},
		[scoreRef],
	)

	useEffect(() => {
		const cancelled = { value: false }
		void loadInitialData(cancelled)

		return () => {
			cancelled.value = true
		}
	}, [loadInitialData])

	const updateHistory = useCallback(
		(code: string) => {
			setRecentCorrect(prev => {
				const nextHistory = clampRecentHistory(prev, code)
				recentRef.current = nextHistory
				return nextHistory
			})
		},
		[recentRef],
	)

	const focusOption = useCallback(
		(nextIndex: number) => {
			if (!roundRef.current) return
			const safeIndex = (nextIndex + roundRef.current.options.length) % roundRef.current.options.length
			setFocusedIndex(safeIndex)
			queueMicrotask(() => {
				optionRefs.current[safeIndex]?.focus()
			})
		},
		[roundRef],
	)

	const handleNextRound = useCallback(() => {
		if (statusRef.current === 'playing') {
			// Prevent skipping mid-round without using Escape.
			return
		}

		const isCorrect = selection.result === 'correct'
		const updatedSession: GameSession = {
			...session,
			sessionCorrect: session.sessionCorrect + (isCorrect ? 1 : 0),
		}

		// Check if session is complete
		if (session.currentRound >= session.totalRounds) {
			// End session
			const endedSession: GameSession = {
				...updatedSession,
				mode: 'ended',
			}
			setSession(endedSession)
			void persistSessionSafely(endedSession)
			return
		}

		// Continue to next round
		const nextSession: GameSession = {
			...updatedSession,
			currentRound: session.currentRound + 1,
		}

		const next = createRound(COUNTRY_FLAGS, recentRef.current)
		setRound(next)
		setSelection({ choice: null, result: 'idle' })
		setStatus('playing')
		setFlagError(false)
		setFocusedIndex(0)
		setSession(nextSession)
		roundRef.current = next
		void persistSessionSafely(nextSession)
		queueMicrotask(() => {
			optionRefs.current[0]?.focus()
		})
	}, [recentRef, roundRef, statusRef, session, selection, persistSessionSafely])

	const finalizeRound = useCallback(
		(selectedCode: string | null) => {
			if (interactionLockRef.current || statusRef.current === 'revealed' || !roundRef.current) {
				return
			}

			interactionLockRef.current = true
			const now = Date.now()
			const { nextScore, result } = summarizeRound(scoreRef.current, roundRef.current, selectedCode, now)

			setScore(nextScore)
			setSelection({ choice: selectedCode, result })
			setStatus('revealed')
			updateHistory(roundRef.current.correct.code)
			void persistScoreSafely(nextScore)

			scoreRef.current = nextScore
			statusRef.current = 'revealed'

			setTimeout(() => {
				nextButtonRef.current?.focus()
				interactionLockRef.current = false
			}, 60)
		},
		[persistScoreSafely, roundRef, scoreRef, statusRef, updateHistory],
	)

	const handleOptionSelect = useCallback(
		(code: string) => {
			if (statusRef.current === 'revealed' || !roundRef.current) {
				return
			}
			setFocusedIndex(roundRef.current.options.findIndex(option => option.code === code))
			setSelection(current => ({ ...current, choice: code }))
			finalizeRound(code)
		},
		[finalizeRound, roundRef, statusRef],
	)

	const handleSkip = useCallback(() => {
		finalizeRound(null)
	}, [finalizeRound])

	const resetGameState = () => {
		const nextRound = createRound(COUNTRY_FLAGS)
		setRound(nextRound)
		roundRef.current = nextRound
		setSelection({ choice: null, result: 'idle' })
		setStatus('playing')
		statusRef.current = 'playing'
		setFlagError(false)
		setFocusedIndex(0)
		setRecentCorrect([])
		recentRef.current = []
		interactionLockRef.current = false
	}

	const handleArrowNavigation = useCallback(
		(event: KeyboardEvent) => {
			if (event.key === 'ArrowRight' || event.key === 'ArrowDown') {
				event.preventDefault()
				focusOption(focusedIndex + 1)
				return true
			}

			if (event.key === 'ArrowLeft' || event.key === 'ArrowUp') {
				event.preventDefault()
				focusOption(focusedIndex - 1)
				return true
			}

			return false
		},
		[focusOption, focusedIndex],
	)

	const handleEscapeKey = useCallback(
		(event: KeyboardEvent) => {
			if (event.key === 'Escape') {
				event.preventDefault()
				if (statusRef.current === 'playing') {
					handleSkip()
				} else {
					handleNextRound()
				}
				return true
			}
			return false
		},
		[handleSkip, handleNextRound, statusRef],
	)

	useEffect(() => {
		const onKeyDown = (event: KeyboardEvent) => {
			const now = performance.now()
			if (now - lastKeyTimestampRef.current < KEY_DEBOUNCE_MS) {
				return
			}
			lastKeyTimestampRef.current = now

			if (handleArrowNavigation(event)) return
			if (handleEscapeKey(event)) return
		}

		window.addEventListener('keydown', onKeyDown)
		return () => window.removeEventListener('keydown', onKeyDown)
	}, [handleArrowNavigation, handleEscapeKey])

	useEffect(() => {
		queueMicrotask(() => {
			optionRefs.current[0]?.focus()
		})
	}, [])

	useEffect(() => {
		if (!flagCardRef.current) {
			return
		}
		utils.remove(flagCardRef.current)
		animate(flagCardRef.current, {
			scale: [0.97, 1],
			opacity: [0.8, 1],
			duration: 520,
			ease: spring({ mass: 1, stiffness: 90, damping: 12, velocity: 0 }),
		})
	}, [])

	useEffect(() => {
		if (status !== 'revealed') {
			return
		}
		const correctIndex = round?.options?.findIndex(option => option.code === round.correct.code)
		if (!correctIndex) {
			return
		}
		const correctButton = optionRefs.current[correctIndex]
		if (!correctButton) {
			return
		}
		utils.remove(correctButton)
		animate(correctButton, {
			scale: [1, 1.04, 1],
			duration: 420,
			ease: 'inOut(2)',
		})
	}, [round, status])

	const currentFeedback = feedbackVariants[selection.result]

	if (!round) {
		return (
			<div className='flex min-h-dvh items-center justify-center'>
				<p>{t('loading_game')}</p>
			</div>
		)
	}

	// Start Screen
	if (session.mode === 'start') {
		const handleStartGame = () => {
			resetGameState()

			const totalRounds = Math.min(selectedRounds, COUNTRY_FLAGS.length)
			const newSession: GameSession = {
				mode: 'playing',
				totalRounds,
				currentRound: 1,
				sessionCorrect: 0,
			}
			setSession(newSession)
			void persistSessionSafely(newSession)
		}

		return (
			<div className='flex min-h-dvh items-center justify-center px-4'>
				<Card className='max-w-md'>
					<h1 className='mb-2 text-center font-bold text-3xl'>{t('flags_guessing_game')}</h1>
					<p className='mb-6 text-center text-info-200'>{t('test_geography_knowledge')}</p>

					<div className='mb-6'>
						<label htmlFor={roundsInputId} className='mb-2 block font-semibold text-info-200 text-sm'>
							{t('how_many_flags_play')}
						</label>
						<div className='space-y-2'>
							<input
								id={roundsInputId}
								type='range'
								min={5}
								max={MAX_ROUNDS / 10}
								step={5}
								value={selectedRounds}
								onChange={event => setSelectedRounds(Number.parseInt(event.target.value, 10))}
								className='w-full accent-pollen-500'
							/>
							<div className='flex items-center justify-between text-sm'>
								<span className='font-semibold'>
									{Math.min(selectedRounds, COUNTRY_FLAGS.length)} {t('flags')}
								</span>
							</div>
						</div>
					</div>

					<button
						type='button'
						className='w-full rounded-lg bg-pollen-500 px-6 py-3 font-semibold text-pollen-100 shadow-md transition hover:translate-y-[-2px] hover:shadow-lg focus:outline-none focus:ring-2 focus:ring-danger-600/75'
						onClick={handleStartGame}
					>
						{t('start_game')}
					</button>
				</Card>
			</div>
		)
	}

	// End Screen
	if (session.mode === 'ended') {
		const sessionAccuracy =
			session.totalRounds > 0 ? Math.round((session.sessionCorrect / session.totalRounds) * 100) : 0

		const handlePlayAgain = () => {
			resetGameState()

			const newSession: GameSession = {
				mode: 'start',
				totalRounds: 0,
				currentRound: 0,
				sessionCorrect: 0,
			}
			setSession(newSession)
			void persistSessionSafely(newSession)
		}

		return (
			<div className='flex min-h-dvh items-center justify-center px-4'>
				<Card className='w-full max-w-md rounded-xl p-8 shadow-lg'>
					<h1 className='mb-4 text-center font-bold text-3xl'>{t('session_complete')}</h1>

					<div className='mb-6 rounded-lg bg-info-800/60 p-6 text-center text-info-100'>
						<p className='mb-2'>{t('you_got')}</p>
						<p className='mb-2 font-bold text-5xl'>
							{session.sessionCorrect}/{session.totalRounds}
						</p>
						<p className=''>{t('correct')}</p>
						<p className='mt-4 font-semibold text-2xl'>
							{sessionAccuracy}% {t('accuracy')}
						</p>
					</div>

					<div className='mb-6 space-y-2 text-info-100'>
						<div className='flex justify-between rounded-lg bg-info-800/60 px-4 py-3'>
							<span>{t('current_streak')}</span>
							<span className='font-semibold'>{score.currentStreak}</span>
						</div>
						<div className='flex justify-between rounded-lg bg-info-800/60 px-4 py-3'>
							<span>{t('best_streak')}</span>
							<span className='font-semibold'>{score.bestStreak}</span>
						</div>
						<div className='flex justify-between rounded-lg bg-info-800/60 px-4 py-3'>
							<span>{t('overall_accuracy')}</span>
							<span className='font-semibold'>{accuracy.toFixed(1)}%</span>
						</div>
					</div>

					<button
						type='button'
						className='w-full rounded-lg bg-pollen px-6 py-3 font-semibold text-pollen-100 transition hover:translate-y-[-2px] hover:shadow-lg focus:outline-none focus:ring-2 focus:ring-danger-600/75'
						onClick={handlePlayAgain}
					>
						{t('play_again')}
					</button>
				</Card>
			</div>
		)
	}

	return (
		<div className='mx-auto grid max-w-4xl grid-cols-1 place-items-start gap-4 px-4 py-8 lg:grid-cols-12'>
			<Card className='flex flex-col gap-4 p-4 lg:col-span-4 lg:max-w-xs'>
				<header className='flex items-center justify-between'>
					<div>
						<p className='font-semibold text-info-200 text-sm uppercase tracking-wide'>{t('scoreboard')}</p>
					</div>
					<span className='font-semibold text-xs opacity-75'>
						{accuracy.toFixed(1)}% {t('accuracy')}
					</span>
				</header>
				<div className='grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-1'>
					<Stat label={t('total_rounds')} value={score.totalGames} />
					<Stat label={t('correct')} value={score.correctAnswers} />
					<Stat label={t('current_streak')} value={score.currentStreak} />
					<Stat label={t('best_streak')} value={score.bestStreak} />
					<Stat label={t('accuracy')} value={`${accuracy.toFixed(1)}%`} />
					<Stat label={t('last_played')} value={formatDate(score.lastPlayedAt)} />
				</div>

				{persistenceWarning ? (
					<p className='rounded-lg border border-amber-700/60 bg-amber-900/50 px-3 py-2 text-amber-100 text-xs'>
						{persistenceWarning}
					</p>
				) : null}
				{loadingScore ? <p className='text-stone-300 text-xs'>{t('loading_saved_progress')}</p> : null}
			</Card>

			<Card className='flex flex-1 flex-col gap-4 p-4 lg:col-span-8'>
				{session.mode === 'playing' && (
					<div className='rounded-lg border border-info-800/40 p-3 px-3 py-2'>
						<p className='mb-1 font-semibold text-xs uppercase tracking-wide'>{t('session_progress')}</p>
						<div className='flex items-center gap-2'>
							<div className='h-2 flex-1 overflow-hidden rounded-full bg-info-800/50'>
								<div
									className='h-full bg-pollen-500 transition-all duration-300'
									style={{ width: `${(session.currentRound / session.totalRounds) * 100}%` }}
								/>
							</div>
							<p className='font-bold text-sm'>
								{session.currentRound}/{session.totalRounds}
							</p>
						</div>
						<div className='mt-2 flex items-center justify-end gap-3 text-xs'>
							<span className='text-olive-300'>
								✓ {session.sessionCorrect} {t('correct')}
							</span>
							<span className='text-danger-300'>
								✗ {session.currentRound - 1 - session.sessionCorrect} {t('errors')}
							</span>
						</div>
					</div>
				)}
				<header>
					<p className='font-bold text-lg'>{t('which_country_flag')}</p>
				</header>

				<div ref={flagCardRef} className='relative aspect-video w-full overflow-hidden rounded-lg bg-info-900/80'>
					<span
						className={`fi fi-${round.correct.code.toLowerCase()} absolute inset-4 rounded-lg bg-center bg-cover text-transparent`}
						role='img'
						style={{ width: 'calc(100% - 2rem)', height: 'calc(100% - 2rem)' }}
					/>
					{flagError ? (
						<div className='absolute inset-0 flex flex-col items-center justify-center gap-3 bg-white/80 text-center text-stone-700-sm'>
							<p className='font-semibold text-sm'>{t('flag_unavailable')}</p>
							<button
								type='button'
								className='rounded-full border bg-info-700 px-4 py-2 font-semibold text-sm transition hover:shadow-sm focus:outline-none focus:ring-2 focus:ring-danger-600/75'
								onClick={() => {
									setFlagError(false)
								}}
							>
								{t('retry_image')}
							</button>
						</div>
					) : null}
				</div>

				<div className='min-h-18 rounded-lg bg-info-800/40 p-3' aria-live='polite'>
					<p className='font-semibold text-sm'>{currentFeedback}</p>
					{selection.result === 'incorrect' && selection.choice ? (
						<p className='text-sm'>
							{t('correct_answer')}: <span className='font-semibold'>{round.correct.name}</span>
						</p>
					) : null}
					{selection.result === 'skipped' ? (
						<p className='text-sm'>
							{t('skipped_flag_belongs')} <span className='font-semibold'>{round.correct.name}</span>.
						</p>
					) : null}
				</div>

				<div className='grid grid-cols-1 gap-3 sm:grid-cols-2'>
					{round.options.map((option, index) => {
						const isCorrect = option.code === round.correct.code
						const isSelected = selection.choice === option.code
						const isRevealed = status === 'revealed'

						const baseClasses =
							'flex h-full w-full items-center gap-3 rounded-full px-4 py-3 text-left text-sm font-semibold transition focus:outline-none focus:ring-2 focus:ring-danger-600/75'
						const playingClasses = 'border border-info-600/60 text-info-100 hover:-translate-y-0.5 hover:shadow-sm'
						const correctClasses = 'bg-olive-600 text-olive-100'
						const incorrectClasses = 'bg-danger-600 text-danger-100'

						const stateClasses = isRevealed
							? isCorrect
								? correctClasses
								: isSelected
									? incorrectClasses
									: 'opacity-80'
							: playingClasses

						return (
							<button
								key={option.code}
								ref={element => {
									optionRefs.current[index] = element
								}}
								type='button'
								className={`${baseClasses} ${stateClasses}`}
								onClick={() => handleOptionSelect(option.code)}
								disabled={isRevealed}
								aria-pressed={isSelected}
								aria-label={`Option ${index + 1}: ${option.name}`}
							>
								<span className='flex h-8 w-8 items-center justify-center rounded-full bg-info-700/40 font-bold text-xs'>
									{index + 1}
								</span>
								<span className='truncate'>{option.name}</span>
							</button>
						)
					})}
				</div>

				<div className='mt-6 flex w-full flex-col gap-3 sm:flex-row sm:items-center sm:justify-between'>
					<button
						type='button'
						className='rounded-lg px-4 py-2 font-semibold text-sm shadow-sm transition hover:translate-y-[-1px] hover:shadow focus:outline-none focus:ring-2 focus:ring-danger-600/75 disabled:cursor-not-allowed disabled:opacity-60'
						onClick={handleSkip}
						disabled={status === 'revealed'}
					>
						{t('skip_round')}
					</button>
					<button
						ref={nextButtonRef}
						type='button'
						className='rounded-lg px-4 py-2 font-semibold text-sm shadow-sm transition hover:translate-y-[-1px] hover:shadow focus:outline-none focus:ring-2 focus:ring-danger-600/75 disabled:cursor-not-allowed disabled:opacity-60'
						onClick={handleNextRound}
						disabled={status === 'playing'}
					>
						{t('next_round')}
					</button>
				</div>
			</Card>
		</div>
	)
}

type StatProps = {
	label: string
	value: string | number
}

const Stat = ({ label, value }: StatProps): React.ReactElement => {
	return (
		<div className='rounded-lg border border-info-800/40 px-3 py-2 opacity-75'>
			<p className='font-semibold text-info-200/80 text-xs uppercase tracking-wide'>{label}</p>
			<p className='font-bold text-sm'>{value}</p>
		</div>
	)
}
