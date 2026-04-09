import { SignupPage } from '_/features/auth/ui/pages/signup-page'
import type { Locale } from '_/i18n/dictionaries/types'
import { setStaticParamsLocale } from 'next-international/server'

type PageProps = {
	params: Promise<{ locale: string }>
}

const Welcome = async ({ params }: PageProps) => {
	const resolvedParams = await params
	const locale = resolvedParams.locale as Locale
	setStaticParamsLocale(locale)
	return <SignupPage />
}

// biome-ignore lint/style/noDefaultExport: page
export default Welcome
