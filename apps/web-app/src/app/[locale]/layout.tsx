import { IconSprites } from '@thom/ui/icon'
import { LocaleProviders } from '_/components/providers/locale-providers'
import { HTML } from '_/components/ui/html'
import { AuthProvider } from '_/features/auth/ui/contexts/auth-provider'
import { ServiceWorkerRegistrar } from '_/features/sw/service-worker'
import { defaultLocale, locales } from '_/i18n/dictionaries/locales'
import type { Locale } from '_/i18n/dictionaries/types'
import type { LayoutProps } from '_/types/pages-layouts'

const isLocale = (value: string): value is Locale => locales.includes(value as Locale)

const LocaleLayout = async ({ children, params }: LayoutProps) => {
	const resolvedParams = await params
	const locale = isLocale(resolvedParams.locale) ? resolvedParams.locale : defaultLocale

	return (
		<HTML locale={locale}>
			<LocaleProviders locale={locale}>
				<AuthProvider>
					<body className='relative min-h-dvh w-full overflow-x-hidden'>
						{children}
						<ServiceWorkerRegistrar />
						<IconSprites />
					</body>
				</AuthProvider>
			</LocaleProviders>
		</HTML>
	)
}
// biome-ignore lint/style/noDefaultExport: layout
export default LocaleLayout
