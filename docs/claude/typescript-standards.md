# TypeScript Standards

## Strictly forbidden

- Using `any` type anywhere.
- Type assertions (`as Type`) without proper justification.
- `@ts-ignore` or ignoring TypeScript errors.
- Returning raw infrastructure types from the data layer.

## Required

- All async operations use the `Result<T>` pattern.
- Use domain models for all business logic.
- Infrastructure types are mapped to domain models at boundaries.
- Branded types for domain IDs.
