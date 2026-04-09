import { SidebarLayoutWithAuth } from '_/app/[locale]/(authorized)/sidebar-layout'

type AuthorizedLayoutProps = {
	children: React.ReactNode
}

const AuthorizedLayout = ({ children }: AuthorizedLayoutProps) => {
	return <SidebarLayoutWithAuth>{children}</SidebarLayoutWithAuth>
}

// biome-ignore lint/style/noDefaultExport: Next.js layout requirement
export default AuthorizedLayout
