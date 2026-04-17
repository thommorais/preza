# apps/base

PocketBase backend for Preza. Runs on port `8090` by default.

## Running

```bash
pnpm dev        # from repo root
go run . serve  # from this directory
```

Admin UI: `http://localhost:8090/_/`

---

## Custom API

PocketBase's built-in REST API (`/api/collections/*`) is available for all collections. The routes below are custom endpoints registered on top of it.

### Response shapes

**Invitation**

```json
{
  "ID": "abc123",
  "EventID": "evt456",
  "CreatedBy": "prm789",
  "GuestID": "gst000",
  "Status": "pending",
  "Token": "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4",
  "PlusOnes": 1,
  "PlusOneNames": ["Mariana Alves"],
  "RSVPAt": "2026-04-10T20:00:00Z",
  "CheckedInAt": null,
  "Notes": "VIP frequente.",
  "CreatedAt": "2026-04-01T10:00:00Z",
  "UpdatedAt": "2026-04-01T10:00:00Z"
}
```

**Event**

```json
{
  "ID": "evt456",
  "VenueID": "ven123",
  "CreatedBy": "prm789",
  "Name": "Réveillon Preza 2027",
  "Date": "2026-12-31T23:00:00Z",
  "MaxCapacity": 500,
  "Status": "draft",
  "CoverImage": "",
  "Description": "A maior virada do ano.",
  "CreatedAt": "2026-04-01T10:00:00Z",
  "UpdatedAt": "2026-04-01T10:00:00Z"
}
```

**Guest**

```json
{
  "ID": "gst000",
  "Name": "Ana Almeida",
  "Phone": "+5511912345678",
  "Email": "ana.almeida@gmail.com",
  "Notes": "Lista A.",
  "CreatedBy": "usr111",
  "CreatedAt": "2026-04-01T10:00:00Z",
  "UpdatedAt": "2026-04-01T10:00:00Z"
}
```

---

### Invitations

#### `POST /api/invitations`

Create an invitation. Event must have status `active`.

**Body**

| Field | Type | Required | Description |
| ----- | ---- | -------- | ----------- |
| `event_id` | string | yes | ID of an active event |
| `created_by` | string | yes | Promoter ID |
| `guest_id` | string | yes | Guest ID |
| `plus_ones` | int | no | Number of additional guests |
| `notes` | string | no | Free-text notes |

**Responses**

| Status | Body |
| ------ | ---- |
| `201` | Invitation object. `Status` is always `"pending"`. `Token` is a 32-char hex string. |
| `400` | Event not found or event is not `active` |

---

#### `GET /api/events/{eventId}/invitations`

List all invitations for an event.

**Responses**

| Status | Body |
| ------ | ---- |
| `200` | Array of invitation objects |
| `404` | Event not found |

---

#### `GET /api/promoters/{promoterId}/invitations`

List all invitations created by a promoter.

**Responses**

| Status | Body |
| ------ | ---- |
| `200` | Array of invitation objects |
| `404` | Promoter not found |

---

#### `POST /api/invitations/{token}/check-in`

Mark an invitation as used at the door. Fails if already checked in.

**Path params**

| Param | Description |
| ----- | ----------- |
| `token` | 32-char hex token from the invitation |

**Responses**

| Status | Body |
| ------ | ---- |
| `200` | Invitation object with `Status: "used"` and `CheckedInAt` set |
| `400` | Token not found or invitation already checked in |

---

### Events

#### `POST /api/events`

Create an event. Created with status `draft`.

**Body**

| Field | Type | Required | Description |
| ----- | ---- | -------- | ----------- |
| `venue_id` | string | yes | ID of an existing venue |
| `created_by` | string | yes | Promoter ID |
| `name` | string | yes | Event name |
| `date` | string | yes | RFC3339 datetime, e.g. `"2026-12-31T23:00:00Z"` |
| `max_capacity` | int | no | Maximum number of guests |
| `description` | string | no | Free-text description |

**Responses**

| Status | Body |
| ------ | ---- |
| `201` | Event object. `Status` is always `"draft"`. |
| `400` | Invalid body, invalid date format, venue not found, or promoter not found |

---

#### `GET /api/promoters/{promoterId}/events`

List all events created by a promoter, ordered by date descending.

**Responses**

| Status | Body |
| ------ | ---- |
| `200` | Array of event objects |
| `404` | Promoter not found |

---

### Guests

#### `GET /api/promoters/{promoterId}/guests`

List all guests on a promoter's guest list.

**Responses**

| Status | Body |
| ------ | ---- |
| `200` | Array of guest objects |
| `404` | Promoter not found |

---

## Domain model

### Invitation statuses

| Status | Meaning |
| ------ | ------- |
| `pending` | Created, guest not yet responded |
| `confirmed` | Guest RSVP'd |
| `declined` | Guest declined |
| `used` | Guest checked in at the door |

### Event statuses

| Status | Meaning |
| ------ | ------- |
| `draft` | Created, not yet published |
| `active` | Published - invitations can be created |
| `closed` | Event ended |

### Promoter types

| Type | Meaning |
| ---- | ------- |
| `independent` | Freelance promoter |
| `venue` | Promoter tied to a venue |

### Collections

| Collection | Description |
| ---------- | ----------- |
| `promoters` | Promoter profiles linked to a PocketBase user |
| `venues` | Venues owned by a promoter |
| `promoter_venues` | Many-to-many: promoters associated with venues, with per-promoter `invite_limit` |
| `events` | Events at a venue, created by a promoter |
| `guests` | Guest records |
| `guest_lists` | Many-to-many: guests on a promoter's list |
| `invitations` | Invitations linking a guest to an event via a unique token |

---

## Architecture

Hexagonal (ports and adapters):

```
internal/
  domain/               # pure Go - no external deps
    promoter/
    venue/
    guest/
    event/
    invitation/         # includes CheckIn() domain method
  application/          # use cases, depend only on domain ports
    invitation/         # CreateInvitation, CheckInGuest, ListByEvent, ListByPromoter
    event/              # CreateEvent, ListByPromoter
    guest/              # ListByPromoter
  infrastructure/
    pocketbase/         # adapters implementing domain repository ports
    http/               # route registration (hooks.go) and handlers per domain
```

`main.go` is the composition root - wires adapters into use cases into handlers.

## Tests

```bash
go test ./internal/...
```

Integration tests clone `pb_data` (seeded with mock promoters, venues, events, guests, and invitations) - no network required. Domain tests are pure unit tests with no dependencies.
