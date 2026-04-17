import clsx from 'clsx'
import { Link } from './link'

export function Text({ className, ...props }: React.ComponentPropsWithoutRef<'p'>) {
	return (
		<p
			data-slot='text'
			{...props}
			className={clsx(className, 'text-base/6 text-zinc-500 sm:text-sm/6')}
			data-id='preza-ui'
		/>
	)
}

export function TextLink({ className, ...props }: React.ComponentPropsWithoutRef<typeof Link>) {
	return (
		<Link
			{...props}
			className={clsx(className, 'text-black underline decoration-black/50 data-hover:decoration-black')}
			data-id='preza-ui'
		/>
	)
}

export function Strong({ className, ...props }: React.ComponentPropsWithoutRef<'strong'>) {
	return <strong {...props} className={clsx(className, 'font-medium text-black')} data-id='preza-ui' />
}

export function Code({ className, ...props }: React.ComponentPropsWithoutRef<'code'>) {
	return (
		<code
			{...props}
			className={clsx(
				className,
				'rounded-sm border border-black/10 bg-black/2.5 px-0.5 font-medium text-black text-sm sm:text-[0.8125rem]',
			)}
			data-id='preza-ui'
		/>
	)
}
