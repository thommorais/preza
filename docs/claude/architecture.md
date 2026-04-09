# Architecture

## Hexagonal architecture (ports and adapters)

Critical rules:

- Domain layer has zero external dependencies.
- Infrastructure types must be mapped to domain models at adapter boundaries.
- Repositories return domain models, never raw API/DB types.
- UI consumes only domain models via port interfaces.

## Result pattern for error handling

All async operations use the `Result<T>` pattern instead of throwing exceptions:

```typescript
type Result<T, E = Error> =
  | { success: true; value: T }
  | { success: false; error: E }

import { ok, err } from '_/lib/result'

const result = await repository.getCardById(id)
if (result.success) {
  // result.value is type T
} else {
  // result.error is type E
}
```

Located in `/src/lib/result.ts`.

## Domain models with branded types

Domain uses branded types for type safety:

```typescript
type CardId = string & { readonly __brand: 'CardId' }
type DeckId = string & { readonly __brand: 'DeckId' }

const id = cardId('abc-123')  // CardId
const deck = deckId('xyz-456') // DeckId
```

## Discriminated union state machines

State is managed with discriminated unions for exhaustive type checking:

```typescript
type SessionState =
  | { readonly type: 'idle' }
  | { readonly type: 'studying'; readonly data: StudyingState }
  | { readonly type: 'answered'; readonly data: AnsweredState }
  | { readonly type: 'completed'; readonly data: CompletedState }
  | { readonly type: 'error'; readonly error: string }

const newState = reduceSessionState(currentState, action)
```

## Dependency injection pattern

Repositories are injected via React Context:

```typescript
<FlashcardRepositoriesProvider repositories={{
  cardRepository: createMockedCardRepository(),
  deckRepository: createMockedDeckRepository(),
}}>
  <App />
</FlashcardRepositoriesProvider>

const { cardRepository } = useFlashcardRepositories()
const result = await cardRepository.getDueCards(deckId)
```
