# Architecture

## Hexagonal architecture (ports and adapters)

The flashcards feature exemplifies the architectural pattern for this codebase:

```
/src/features/flashcards/
├── domain/              # Pure business logic (no external dependencies)
│   ├── types.ts         # Domain models with branded types
│   ├── constants.ts     # Named domain constants
│   ├── rules/           # Pure validation functions
│   ├── reducers/        # Pure state machines
│   └── selectors/       # Pure data extraction
├── ports/               # Abstract interfaces (contracts)
│   └── repositories/    # Repository port definitions
├── adapters/            # Concrete implementations
│   └── repositories/    # Repository adapters with domain mapping
└── ui/                  # React presentation layer
    ├── contexts/        # Dependency injection
    ├── hooks/           # Domain integration hooks
    ├── components/      # Feature-specific UI components
    └── pages/           # Full-page components
```

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

See `/src/features/flashcards/domain/reducers/session-reducer.ts`.

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

## Data flow

1. UI calls hook: `useStudySession(deckId)`
2. Hook uses repository port: `cardRepository.getDueCards()`
3. Adapter fetches infrastructure data: API/DB/Mock
4. Adapter maps to domain model: `mapToCard(rawData)`
5. Returns `Result<DomainModel>`: `ok(cards)` or `err(error)`
6. Hook updates state via reducer: `reduceSessionState(state, action)`
7. UI renders from domain state

Never return infrastructure types from adapters.
