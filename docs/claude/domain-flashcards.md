# Flashcards Domain (Leitner)

The flashcard feature implements spaced repetition with 7 box levels.

- Correct answer: card moves to next box level with longer review interval.
- Incorrect answer: card moves back to earlier box for more frequent review.
- Card types: `multiple_choice` (discrete options) and `short_answer` (free text).

Key domain logic:

- `LeitnerBox` entity tracks card progression through boxes.
- `getDueCards()` filters cards based on box level and due date.
- `submitAnswer()` updates box level and calculates next review time.

See `/src/features/flashcards/domain/` for business rules.
