import type { Locale } from '_/i18n/dictionaries/types'
import type { Metadata } from 'next'
import { setStaticParamsLocale } from 'next-international/server'
import { FlagsGameClient } from './flags-client'

export const metadata: Metadata = {
	title: 'Flags Guessing Game',
	description: 'Guess the country based on the flag, track your streak, and improve your geography skills.',
}

type FlagsGamePageProps = {
	params: Promise<{ locale: string }>
}

const FlagsGamePage = async ({ params }: FlagsGamePageProps): Promise<React.ReactNode> => {
	const resolvedParams = await params
	const locale = resolvedParams.locale as Locale
	setStaticParamsLocale(locale)

	return (
		<main className='min-h-dvh bg-info-900 text-info-100'>
			<FlagsGameClient />
		</main>
	)
}

// biome-ignore lint/style/noDefaultExport: page
export default FlagsGamePage
