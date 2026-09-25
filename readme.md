# Kuu — Full-Stack Social Network

**Kuu** is a real-time social networking platform with a **Go** backend and a **Next.js 15** frontend: private/public profiles, followers with request approval, posts with granular privacy, groups with events, direct & group chat over WebSockets, and notifications.

## Overview

- **Backend**: REST + WebSocket API written in Go, SQLite for persistence.
- **Frontend**: React 19 / Next.js 15 (App Router), plain CSS, Context API.

### Core features

- **Social graph**: follow/unfollow; private profiles require a follow request (public profiles accept instantly).
- **Profiles**: all register fields, activity feed, followers/following lists, public/private toggle.
- **Posts & comments**: image/GIF attachments, privacy `public` / `almost private` / `private` (chosen viewers).
- **Groups**: create, invite (accept/decline), join requests, member posts/comments, events with Going / Not going votes, shared group chat.
- **Chat**: 1:1 direct messages (only between users who follow each other) delivered in real time over WebSockets, emoji support.
- **Notifications**: follow requests, group invitations, join requests, events, and more — visible on every page.

## Tech stack

| Layer | Technologies |
| :--- | :--- |
| **Frontend** | Next.js 15 (App Router), React 19, Context API, plain CSS |
| **Backend** | Go 1.25, `net/http`, Gorilla WebSocket |
| **Database** | SQLite (`mattn/go-sqlite3`) |
| **Migrations** | `golang-migrate`, applied automatically on boot |
| **Auth** | Session cookie (HTTP-only) + bcrypt |
| **Deployment** | Docker, Docker Compose |

## Project structure

```text
.
├── backend/             # Go API
│   ├── main.go          # entry point → internal/cmd
│   ├── internal/        # config, database, handler, service, repository,
│   │                    # routes, middleware, websocket, models, requests, helper
│   ├── migrations/      # SQL migrations (golang-migrate)
│   └── media/           # uploaded images
├── frontend/            # Next.js app
│   ├── app/             # App Router pages
│   ├── components/      # UI components
│   ├── contexts/        # global state providers
│   ├── css/             # plain CSS (Hollow Knight theme)
│   └── lib/             # API layer + utils
└── docker-compose.yml   # both services
```

## Getting started

### Docker (recommended)

```bash
docker compose up --build
```

- **Frontend**: `http://localhost:3000`
- **Backend API**: `http://localhost:8080`

### Manual setup

#### 1. Backend

```bash
cd backend
go run .
```

*Migrations run automatically on startup; the database is created at `backend/internal/database/social.db`.*

#### 2. Frontend

```bash
cd frontend
npm install
npm run dev
```

*The frontend talks to `http://localhost:8080` by default (`NEXT_PUBLIC_API_URL` overrides it).*

## Authentication & real-time

### Session management

1. Login/register sets an HTTP-only `session_token` cookie.
2. Every frontend request goes through `lib/api/fetcher.js` with `credentials: "include"`.
3. `AuthContext` restores the session on page load via `GET /api/v1/auth/me`; logout clears it.

### WebSocket communication

- Single connection per user: `ws://localhost:8080/ws` (session cookie required), managed by `WebSocketContext.jsx`.
- Components call `send(type, payload)` to emit events and `subscribe(type, callback)` to listen (chat, notifications).

| Direction | Event types |
| :--- | :--- |
| Client → server | `send_direct_message`, `send_group_message` |
| Server → client | `new_direct_message`, `new_group_message`, `new_notification`, `notification_read`, `notification_expired` |

## Backend internals

Request path: `routes` → `handler` → `service` → `repository` → SQLite. Middleware: CORS (allows `http://localhost:3000` with credentials), `AllowMethods`, `RequireAuth` (reads the `session_token` cookie, injects the user). A per-minute rate limiter wraps the auth endpoints (20/min) and the whole API (300/min), skipping `OPTIONS`.

Every JSON endpoint returns the same envelope:

```json
{
  "success": true,
  "message": "optional human message",
  "data": { },
  "errors": [ ]
}
```

Validation failures: HTTP `400` with an `errors` array. Auth failures: HTTP `401`. Uploaded images are stored in `backend/media/` and served at `/media/<filename>`.

## Frontend internals

- `lib/api/*.js` — one module per domain over `apiFetch` (prepends `API_BASE`, sends `credentials: "include"`, throws on non-2xx).
- Providers mounted in `app/layout.jsx`: `Auth`, `WebSocket`, `Notification`, `GroupChat`, `Audio`, `Particles`. The `(main)` layout adds the authenticated shell and redirects logged-out users to `/login`.
- On load `AuthContext` restores the session via `GET /api/v1/auth/me`; `WebSocketContext` opens the socket with auto-reconnect.

| Route | Page |
| :--- | :--- |
| `/login`, `/register` | `app/login`, `app/register` |
| `/` | home feed — `app/(main)/page.jsx` |
| `/posts/create`, `/posts/:postId` | create form, post detail |
| `/profile/:username` | profile |
| `/group`, `/group/:id`, `/group/:id/events` | group list, group detail, events |
| `/error`, 404 | `app/error/page.jsx`, `app/not-found.jsx` |

## Database schema

Tables (see `backend/migrations/*.up.sql`):

| Table | Purpose |
| :--- | :--- |
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
| `notifications` | typed notifications (JSON `payload`/`actions`, optional `group_id`) |
| `post_viewers` | allowed viewers for `private` posts |
| `group_events` | scheduled events (`upcoming`/`cancelled`/`expired`) |
| `event_responses` | per-user `going` / `not_going` |

## API summary

| Domain | Key endpoints |
| :--- | :--- |
| **Auth** | `POST /api/v1/auth/register`, `POST /api/v1/auth/login`, `GET /api/v1/auth/me` |
| **Users** | `GET /api/v1/user/profile`, `POST /api/v1/user/follow`, `GET /api/v1/user/followers` |
| **Posts** | `GET /api/v1/posts/feed`, `POST /api/v1/posts`, `POST /api/v1/posts/comments` |
| **Groups** | `GET /api/v1/groups/all`, `POST /api/v1/groups/invite`, `GET /api/v1/groups/events` |
| **Chat** | `GET /api/v1/chat/conversations`, `GET /api/v1/chat/direct`, `GET /api/v1/chat/group` |
| **Notifications** | `GET /api/v1/notifications`, `POST /api/v1/notifications/read` |
| **WebSocket** | `/ws` (direct messages, group messages, notifications) |

## Development notes

- **Ports**: backend `8080`, frontend `3000` (dev and Docker).
- **Styling**: plain CSS under `frontend/css/` — no UI library, Hollow Knight theme.
- **Images**: stored in `backend/media/`, served via `/media/<filename>`.

## Authors

- **Mbelhouss**
- **mbarrah**
