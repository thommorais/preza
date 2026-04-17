import { tv } from './tv'

type HeadingProps = { level?: 1 | 2 | 3 | 4 | 5 | 6 } & React.ComponentPropsWithoutRef<
	'h1' | 'h2' | 'h3' | 'h4' | 'h5' | 'h6'
>

const headingClasses = tv({
	base: ['font-semibold text-2xl/8 text-black sm:text-xl/8'],
})

export function Heading({ className, level = 1, ...props }: HeadingProps) {
	const Element: `h${typeof level}` = `h${level}`

	return <Element {...props} className={headingClasses({ class: [className] })} />
}

const subHeadingClasses = tv({
	base: ['font-semibold text-base/7 text-black sm:text-sm/6'],
})

export function Subheading({ className, level = 2, ...props }: HeadingProps) {
	const Element: `h${typeof level}` = `h${level}`

	return <Element {...props} className={subHeadingClasses({ class: [className] })} data-id='preza-ui' />
}
