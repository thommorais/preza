'use client'

import * as Headless from '@headlessui/react'
import clsx from 'clsx'
import { LayoutGroup, motion } from 'motion/react'
import type React from 'react'
import { useId } from 'react'
import { TouchTarget } from './button'
import { Link } from './link'

export function Sidebar({ className, ...props }: React.ComponentPropsWithoutRef<'nav'>) {
	return <nav {...props} className={clsx(className, 'flex h-full min-h-0 flex-col')} />
}

export function SidebarHeader({ className, ...props }: React.ComponentPropsWithoutRef<'div'>) {
	return (
		<div
			{...props}
			className={clsx(
				className,
				'flex flex-col border-primary-200/50 border-b p-4 [&>[data-slot=section]+[data-slot=section]]:mt-2.5',
			)}
		/>
	)
}

export function SidebarBody({ className, ...props }: React.ComponentPropsWithoutRef<'div'>) {
	return (
		<div
			{...props}
			className={clsx(
				className,
				'flex flex-1 flex-col overflow-y-auto p-4 [&>[data-slot=section]+[data-slot=section]]:mt-8',
			)}
		/>
	)
}

export function SidebarFooter({ className, ...props }: React.ComponentPropsWithoutRef<'div'>) {
	return (
		<div
			{...props}
			className={clsx(
				className,
				'flex flex-col border-primary-200/50 border-t p-4 [&>[data-slot=section]+[data-slot=section]]:mt-2.5',
			)}
		/>
	)
}

export function SidebarSection({ className, ...props }: React.ComponentPropsWithoutRef<'div'>) {
	const id = useId()

	return (
		<LayoutGroup id={id}>
			<div {...props} data-slot='section' className={clsx(className, 'flex flex-col gap-0.5')} />
		</LayoutGroup>
	)
}

export function SidebarDivider({ className, ...props }: React.ComponentPropsWithoutRef<'hr'>) {
	return <hr {...props} className={clsx(className, 'my-4 border-primary-200/50 border-t lg:-mx-4')} />
}

export function SidebarSpacer({ className, ...props }: React.ComponentPropsWithoutRef<'div'>) {
	return <div aria-hidden='true' {...props} className={clsx(className, 'mt-8 flex-1')} />
}

export function SidebarHeading({ className, ...props }: React.ComponentPropsWithoutRef<'h3'>) {
	return (
		<h3
			{...props}
			className={clsx(className, 'mb-1 px-2 font-medium text-primary-700 text-xs/6 uppercase tracking-wider')}
		/>
	)
}

export const SidebarItem = ({
	current,
	className,
	children,
	LinkComponent = Link,
	ref,
	...props
}: {
	current?: boolean
	className?: string
	children: React.ReactNode
	// biome-ignore lint/suspicious/noExplicitAny: yeap
	LinkComponent?: React.ComponentType<any>
	ref?: React.ForwardedRef<HTMLAnchorElement | HTMLButtonElement>
} & (
	| ({ href?: never } & Omit<Headless.ButtonProps, 'as' | 'className'>)
	| ({ href: string } & Omit<Headless.ButtonProps<typeof Link>, 'as' | 'className'>)
)) => {
	const classes = clsx(
		// Base
		'flex w-full items-center gap-3 rounded-lg px-2 py-2.5 text-left font-medium text-base/6 text-primary-900 lowercase sm:py-2 sm:text-sm/5',
		// Leading icon/icon-only
		'*:data-[slot=icon]:size-6 *:data-[slot=icon]:shrink-0 *:data-[slot=icon]:stroke-primary-600 sm:*:data-[slot=icon]:size-5',
		// Trailing icon (down chevron or similar)
		'*:last:data-[slot=icon]:ml-auto *:last:data-[slot=icon]:size-5 sm:*:last:data-[slot=icon]:size-4',
		// Avatar
		'*:data-[slot=avatar]:-m-0.5 *:data-[slot=avatar]:size-7 sm:*:data-[slot=avatar]:size-6',
		// Hover
		'data-hover:bg-primary-100/50 data-hover:*:data-[slot=icon]:stroke-primary-700',
		// Active
		'data-active:bg-primary-200/50 data-active:*:data-[slot=icon]:stroke-primary-800',
		// Current
		'data-current:bg-primary-100 data-current:*:data-[slot=icon]:stroke-primary-800',
	)

	return (
		<span className={clsx(className, 'relative')}>
			{current && (
				<motion.span
					layoutId='current-indicator'
					className='absolute inset-y-2 -left-4 w-0.5 rounded-full bg-primary-600'
				/>
			)}
			{typeof props.href === 'string' ? (
				<Headless.CloseButton
					as={LinkComponent}
					{...props}
					className={classes}
					data-current={current ? 'true' : undefined}
					ref={ref}
				>
					{children}
				</Headless.CloseButton>
			) : (
				<Headless.Button
					{...props}
					className={clsx('cursor-default', classes)}
					data-current={current ? 'true' : undefined}
					ref={ref}
				>
					<TouchTarget>{children}</TouchTarget>
				</Headless.Button>
			)}
		</span>
	)
}

export const SidebarLabel = ({ className, ...props }: React.ComponentPropsWithoutRef<'span'>) => {
	return <span {...props} className={clsx(className, 'truncate')} />
}
