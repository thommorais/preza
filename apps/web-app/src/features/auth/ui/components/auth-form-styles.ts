import { tv } from '_/lib/third-party/tv'

export const authFormStyles = tv({
	slots: {
		heading: 'mb-8 text-center font-bold text-3xl text-info-700',
		form: 'w-full space-y-6 px-4',
		field: '',
		label: 'mb-2 block font-medium text-info-700 text-sm',
		required: 'text-red-500',
		input: 'w-full rounded-sm',
		error: 'mt-1 text-red-500 text-sm',
		submitError: 'rounded-md bg-red-50 p-4',
		submitErrorText: 'text-red-800 text-sm',
		button: 'w-full',
	},
	variants: {
		hasError: {
			true: {
				input: 'border-red-500',
			},
			false: {
				input: 'border-info-700/20',
			},
		},
	},
})

export const styles = authFormStyles()
