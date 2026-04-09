# UI Forms and i18n

## Form state management

Forms use pure reducers for state management:

- `reduceDeckForm()` - Create/edit deck with validation
- `reduceCardForm()` - Multi-step card creation

Validation happens via pure functions in `/domain/rules/` and executes on form submission, storing errors in form state.

## Internationalization

Uses `next-international` with dynamic locale routing:

- Routes are `app/[locale]/...`
- Static locale params exported from i18n config
- All text must use translation keys
