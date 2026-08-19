# Backend — `kuu`

A REST + WebSocket API for a social network (posts, comments, groups, events, follows, direct/group chat, notifications). Written in **Go**, backed by **SQLite**.

## Tech stack

| Concern   | Choice |
|-----------|--------|
| Language  | Go (module `kuu`, `go 1.25`) |
| HTTP      | Standard library `net/http` + `http.ServeMux` |
| Database  | SQLite via `mattn/go-sqlite3` |
| Migrations| `golang-migrate` (applied automatically on boot) |
| Real-time | WebSockets via `gorilla/websocket` |
| Auth      | Session tokens in an HTTP-only cookie (no JWT) |
| Passwords | bcrypt (`golang.org/x/crypto`) |
| IDs       | `google/uuid` (sessions) |

## Running it

```bash
cd backend
go run .                 # listens on :8080
```

- The database file is `./internal/database/social.db` (path set in `internal/config/config.go`).
- Migrations run automatically from `migrations/` on startup.
- Uploaded images are stored in `./media/` and served at `/media/<filename>`.
- Port and DB path are hard-coded in `internal/config/config.go` (no `.env` currently loaded).

Or run the whole stack with Docker from the repo root:

```bash
docker compose up --build
```

## Layered architecture

The code is split into clean layers. Dependencies point one way: **routes → handler → service → repository → database**.

```
main.go
  └─ cmd/server.go          # wires everything together (composition root)
       ├─ config            # hard-coded port / DB path
       ├─ database          # opens SQLite + runs migrations
       ├─ repository        # raw SQL queries (one file per domain)
       ├─ service           # business logic + validation glue
       ├─ handler           # HTTP handlers (parse request → call service → write JSON)
       ├─ middleware        # CORS, method guards, session auth
       ├─ routes            # registers every endpoint on the mux
       └─ websocket         # Hub (in-memory fan-out for real-time events)
```

### Wiring order (in `cmd/server.go`)

1. Load config, open DB, run migrations.
2. Create the WebSocket `Hub` and start its `Run()` loop in a goroutine.
3. `repository.New(db)` → `service.New(repo)` → `handler.New(svc, hub)` → `middleware.New(svc)`.
4. `routes.Register(h, m)` mounts all routes; `middleware.CORS` wraps the mux.

### Key packages

- **`models/`** — structs that double as the JSON contract (see `json:` tags) and DB rows.
- **`requests/`** — request parsing + field validation (multipart form and JSON).
- **`handler/`** — thin HTTP layer; one file per domain (`auth`, `user`, `post`, `group`, `chat`, `notification`, `webSocket`).
- **`service/`** — business logic and DB orchestration.
- **`repository/`** — SQL. All queries live here, not in handlers.
- **`middleware/`** — `CORS`, `AllowMethods`, and `RequireAuth` (reads the `session_token` cookie, injects the user into the request context).
- **`websocket/`** — the `Hub` keeps every open connection and routes typed `Event`s to users/groups.
- **`helper/`** — JSON response writers (`Success` / `Error` / `ValidationErrorResponse`) and image upload helpers.

## Auth model

- **Register** creates a user (bcrypt-hashed password).
- **Login** verifies credentials and creates a `sessions` row; the session ID is set as an HTTP-only `session_token` cookie.
- **Logout** deletes the session and clears the cookie.
- **`RequireAuth`** middleware checks the cookie on every protected route, resolves the user, and stores it in the request context (`middleware.UserContextKey`).

## Response envelope

Every JSON endpoint returns the same shape:

```json
{
  "success": true,
  "message": "optional human message",
  "data": { },        // present on success
  "errors": [ ]       // present on validation failure (array of {field, message} or strings)
}
```

- Validation failures: HTTP `400` with `"message": "Validation failed"` and an `errors` array.
- Auth failures: HTTP `401`.
- CORS allows only `http://localhost:3000` with credentials.

## API reference

Base URL: `http://localhost:8080` · Everything under `/api/v1/*` (except `/ws` and `/media/`) requires the `session_token` cookie unless marked **public**.

### Auth

| Method | Path | Auth | Body / Query |
|--------|------|------|--------------|
| POST | `/api/v1/auth/register` | public | `multipart/form-data`: `username`, `email`, `first_name`, `last_name`, `gender` (`male`/`female`), `date_of_birth` (`YYYY-MM-DD`, 18+), `password` (≥8), optional `about_me`, optional `avatar` file |
| POST | `/api/v1/auth/login` | public | JSON `{ "login": "username-or-email", "password": "..." }` |
| GET | `/api/v1/auth/me` | ✔ | returns the current user |
| POST | `/api/v1/auth/logout` | ✔ | clears session |

### Users

| Method | Path | Body / Query |
|--------|------|--------------|
| GET | `/api/v1/user/profile` | `?username=` |
| GET | `/api/v1/user/posts` | `?username=` |
| PUT | `/api/v1/user/profile/update` | JSON `{ "is_public": 0\|1 }` |
| POST | `/api/v1/user/follow` | `{ "target_user_id": <id> }` |
| POST | `/api/v1/user/unfollow` | `{ "target_user_id": <id> }` |
| POST | `/api/v1/user/follow/accept` | `{ "target_user_id": <id> }` |
| POST | `/api/v1/user/follow/decline` | `{ "target_user_id": <id> }` |
| GET | `/api/v1/user/follow/requests` | incoming follow requests |
| GET | `/api/v1/user/followers` | `?id=<user_id>` |
| GET | `/api/v1/user/following` | `?id=<user_id>` |
| GET | `/api/v1/user/all` | list all users |
| GET | `/api/v1/user/suggestions` | `?limit=<n>` — suggested users to follow |

### Posts

| Method | Path | Body / Query |
|--------|------|--------------|
| GET | `/api/v1/posts/feed` | `?limit=<n>&cursor=<c>` — paginated feed |
| GET | `/api/v1/posts` | `?post_id=<id>` — single post |
| POST | `/api/v1/posts` | `multipart/form-data` or JSON: `title`, `content`, `privacy` (`public`/`almost private`/`private`), optional `group_id`, optional `visible_to[]` (required for `private`), optional `image` file |
| GET | `/api/v1/posts/comments` | `?post_id=<id>` |
| POST | `/api/v1/posts/comments` | `multipart/form-data` or JSON: `post_id`, `title`, `content`, optional `image` file |

### Groups

| Method | Path | Body / Query |
|--------|------|--------------|
| POST | `/api/v1/groups` | `{ "title", "description", "is_public": 0\|1 }` |
| GET | `/api/v1/groups/all` | list all groups |
| GET | `/api/v1/groups/detail` | `?id=` |
| GET | `/api/v1/groups/feed` | `?id=&limit=&cursor=` |
| GET | `/api/v1/groups/members` | `?id=` |
| PUT | `/api/v1/groups/update` | `?id=` + `{ "title", "description", "is_public" }` |
| DELETE | `/api/v1/groups/delete` | `?id=` |
| POST | `/api/v1/groups/invite` | `{ "group_id", "target_user_ids": [ ... ] }` |
| GET | `/api/v1/groups/invitations` | my invitations |
| POST | `/api/v1/groups/invitations/accept` | `{ "group_id" }` |
| POST | `/api/v1/groups/invitations/decline` | `{ "group_id" }` |
| POST | `/api/v1/groups/join` | `{ "group_id" }` |
| POST | `/api/v1/groups/join/accept` | `{ "group_id", "target_user_id" }` |
| POST | `/api/v1/groups/join/decline` | `{ "group_id", "target_user_id" }` |
| POST | `/api/v1/groups/leave` | `{ "group_id" }` |

### Group events

| Method | Path | Body / Query |
|--------|------|--------------|
| GET | `/api/v1/groups/events` | `?id=<group_id>` |
| POST | `/api/v1/groups/events/create` | `{ "group_id", "title", "description", "event_time" }` |
| POST | `/api/v1/groups/events/cancel` | `{ "event_id" }` |
| POST | `/api/v1/groups/events/respond` | `{ "event_id", "status": "going"\|"not_going" }` |

### Chat (REST history; real-time goes over WebSocket)

| Method | Path | Body / Query |
|--------|------|--------------|
| GET | `/api/v1/chat/conversations` | list DM conversations + unread counts |
| GET | `/api/v1/chat/direct` | `?user_id=&before_id=&limit=` — DM history |
| POST | `/api/v1/chat/read` | `{ "user_id" }` — mark DMs from a user as read |
| GET | `/api/v1/chat/group` | `?group_id=&page=` — group message history |

### Notifications

| Method | Path | Body / Query |
|--------|------|--------------|
| GET | `/api/v1/notifications` | `?limit=&last_id=` — paged list + unread count |
| POST | `/api/v1/notifications/read` | `{ "notification_id": <id> }` **or** `{ "all": true }` |
| POST | `/api/v1/notifications/expire` | `{ "notification_id": <id> }` |

### WebSocket — `/ws`

Connect with `GET /ws` (session cookie required). The frontend connects to `ws://localhost:8080/ws`.

Events are typed JSON objects: `{ "type": "...", "payload": { ... } }`.

**Client → server:**

| Type | Payload |
|------|---------|
| `send_direct_message` | `{ "receiver_id", "content" }` |
| `send_group_message` | `{ "group_id", "content" }` |

**Server → client:**

| Type | Payload |
|------|---------|
| `new_direct_message` | the saved `DirectMessage` |
| `new_group_message` | the saved `GroupMessage` |
| `new_notification` | a `Notification` |
| `notification_read` | read state |
| `notification_expired` | expired state |
| `user_online` / `user_offline` | `{ "user_id" }` |
| `online_users_list` | `{ "online_user_ids": [...] }` (sent on connect) |

### Static

| Path | Description |
|------|-------------|
| `/media/<filename>` | uploaded images (stored under `backend/media/`) |

## Database schema

Tables (see `migrations/*.up.sql`):

| Table | Purpose |
|-------|---------|
| `users` | accounts |
| `sessions` | login sessions |
| `groups` | groups |
| `group_members` | membership (`pending` / `accepted` / `declined`) |
| `posts` | posts (nullable `group_id`, `privacy`) |
| `comments` | post comments |
| `follows` | follow relationships (`pending` / `accepted`) |
| `direct_messages` | 1:1 chat |
| `chat_reads` | per-user DM read cursor |
| `group_messages` | group chat |
| `notifications` | notification stack (typed + JSON `payload`/`actions`) |
| `post_viewers` | allowed viewers for `private` posts |
| `group_events` | scheduled events (`upcoming`/`cancelled`/`expired`) |
| `event_responses` | per-user `going` / `not_going` |

A seed migration (`000009_seed.up.sql`) inserts demo users, groups, posts, comments, and follows. All seed users share the password `SecurePassword123!`.
