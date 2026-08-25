# Frontend

The web client for the social network. A **Next.js 15 (App Router) + React 19** single-page app that talks to the Go backend over REST and WebSocket.

## Tech stack

| Concern | Choice |
|---------|--------|
| Framework | Next.js 15 (App Router) |
| UI | React 19, plain CSS (no component library) |
| Data fetching | native `fetch` wrapped in `lib/api/fetcher.js` |
| Real-time | WebSocket (`contexts/WebSocketContext`) |
| State | React Context (no Redux) |
| Styling | custom CSS files under `css/` |

## Running it

```bash
cd frontend
npm install
npm run dev       # http://localhost:3000
```

The API base URL comes from `NEXT_PUBLIC_API_URL` (defaults to `http://localhost:8000`, see `lib/utils.js`). Set it when the backend is elsewhere:

```bash
NEXT_PUBLIC_API_URL=http://localhost:8000 npm run dev
```

Production build: `npm run build` then `npm start`. Or run the full stack with `docker compose up --build` from the repo root.

## Architecture

Next.js App Router file-based routing under `app/`, with a `lib/` API layer, `contexts/` for global state, and `components/` for UI.

```
app/
  layout.jsx                  # root layout + global providers
  login/page.jsx              # /login
  register/page.jsx           # /register
  (main)/                     # authenticated area (route group)
    layout.jsx                # header, sidebar, notifications, chat, settings
    page.jsx                  # /  (home feed)
    posts/create/page.jsx     # /posts/create
    posts/[postId]/page.jsx   # /posts/:id
    profile/[username]/page.jsx  # /profile/:username
    group/page.jsx            # /group (list)
    group/[id]/page.jsx       # /group/:id (detail + feed)
    group/[id]/events/page.jsx   # /group/:id/events

lib/
  api/                        # one file per domain, thin wrappers over apiFetch
  fetcher.js                  # fetch + cookie + error parsing
  utils.js                    # API_BASE, media URL resolution, date formatting

contexts/                     # global React state providers
  AuthContext.jsx             # session user, login/register/logout
  WebSocketContext.jsx        # ws connection, send/subscribe
  NotificationContext.jsx     # notifications + unread count
  GroupChatContext.jsx        # group chat state
  AudioContext.jsx            # music / sfx
  ParticlesContext.jsx        # mouse particles
  ToastContext.jsx            # toast notifications

components/                   # presentational + feature components
  chat/  groups/  notifications/  posts/  profile/  shared/  sidebar/
```

### Data flow

1. **API layer** — `lib/api/*.js` export functions that call `apiFetch(path, { method, body })`. No component touches `fetch` directly.
2. **`fetcher.js`** — prepends `API_BASE`, sends `credentials: "include"` (so the session cookie flows), JSON-encodes bodies, and throws a friendly `Error` on non-2xx responses (extracting `data.errors`).
3. **Contexts** — hold global state (`AuthContext`, `NotificationContext`, …) and expose hooks (`useAuth`, `useNotifications`, …). Providers are mounted in `app/layout.jsx`.
4. **Pages/components** — call the API layer and read/write context.

### Auth flow

- On load, `AuthContext` calls `GET /api/v1/auth/me` to restore the session.
- The `(main)` route group's layout redirects unauthenticated users to `/login`.
- Login/register hit `/api/v1/auth/*`; the backend sets the `session_token` cookie, which `fetch` carries automatically.

### Real-time (WebSocket)

`WebSocketContext` opens `ws://<API_BASE>/ws` when a user is logged in, with auto-reconnect. It exposes:

- `send(type, payload)` — send a client event.
- `subscribe(type, callback)` — register a listener for a server event (returns an unsubscribe fn).

Components subscribe to events like `new_direct_message`, `new_group_message`, `new_notification`, and `user_online`/`user_offline`.

## Routing (pages)

| Route | Page |
|-------|------|
| `/login` | `app/login/page.jsx` |
| `/register` | `app/register/page.jsx` |
| `/` | `app/(main)/page.jsx` — home feed |
| `/posts/create` | `app/(main)/posts/create/page.jsx` |
| `/posts/:postId` | `app/(main)/posts/[postId]/page.jsx` |
| `/profile/:username` | `app/(main)/profile/[username]/page.jsx` |
| `/group` | `app/(main)/group/page.jsx` |
| `/group/:id` | `app/(main)/group/[id]/page.jsx` |
| `/group/:id/events` | `app/(main)/group/[id]/events/page.jsx` |

Error/404 handling lives in `app/error.jsx`, `app/error/page.jsx`, and `app/not-found.jsx`.

## Backend endpoints

See [`../backend/README.md`](../backend/README.md) for the full API reference. The frontend `lib/api/*.js` files map 1:1 to those endpoints:

| Frontend module | Backend endpoints |
|-----------------|-------------------|
| `lib/api/auth.js` | `/api/v1/auth/*` |
| `lib/api/user.js` | `/api/v1/user/*` |
| `lib/api/posts.js` | `/api/v1/posts/*` |
| `lib/api/groups.js` | `/api/v1/groups/*` |
| `lib/api/chat.js` | `/api/v1/chat/*` |
| `lib/api/notifications.js` | `/api/v1/notifications/*` |
