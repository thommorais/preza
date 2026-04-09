type LocaleLayoutProps = {
	children: React.ReactNode
}

const DashboardLayout = async ({ children }: LocaleLayoutProps) => {
	return <>{children}</>
}

// biome-ignore lint/style/noDefaultExport: layout
export default DashboardLayout
