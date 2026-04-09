import type { Locale } from '_/i18n/dictionaries/types'
import { setStaticParamsLocale } from 'next-international/server'

type LoginProps = {
	params: Promise<{ locale: string }>
}

const Recovery = async ({ params }: LoginProps) => {
	const resolvedParams = await params
	const locale = resolvedParams.locale as Locale
	setStaticParamsLocale(locale)
	return <h1>Recovery Page</h1>
}

// biome-ignore lint/style/noDefaultExport: page
export default Recovery
