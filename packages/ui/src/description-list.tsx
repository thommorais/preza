import clsx from 'clsx'

export function DescriptionList({ className, ...props }: React.ComponentPropsWithoutRef<'dl'>) {
	return (
		<dl
			{...props}
			className={clsx(
				className,
				'grid grid-cols-1 text-base/6 sm:grid-cols-[min(50%,calc(var(--spacing)*80))_auto] sm:text-sm/6',
			)}
		/>
	)
}

export function DescriptionTerm({ className, ...props }: React.ComponentPropsWithoutRef<'dt'>) {
	return (
		<dt
			{...props}
			className={clsx(
				className,
				'col-start-1 border-black/5 border-t pt-3 text-zinc-500 first:border-none sm:border-black/5 sm:border-t sm:py-3',
			)}
		/>
	)
}

export function DescriptionDetails({ className, ...props }: React.ComponentPropsWithoutRef<'dd'>) {
	return (
		<dd
			{...props}
			className={clsx(className, 'pt-1 pb-3 text-black sm:border-black/5 sm:border-t sm:nth-2:border-none sm:py-3')}
		/>
	)
}
