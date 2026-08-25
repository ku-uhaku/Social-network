

# Kuu — Full-Stack Social Network

**Kuu** is a modern, real-time social networking platform. It features a high-performance **Go** backend and a responsive **Next.js 15** frontend, supporting private/public profiles, group management, event scheduling, real-time messaging, and a notification system.

## 🚀 Overview

The project is architected as a decoupled system:
- **Backend**: A RESTful and WebSocket API written in Go, using SQLite for persistence.
- **Frontend**: A React 19 single-page application (SPA) built with Next.js 15 (App Router) and plain CSS.

### Core Features
- **Social Graph**: Follow/Unfollow system with private/public profile toggles.
- **Content**: Posts and comments with image upload support and granular privacy settings (`public`, `almost private`, `private`).
- **Groups**: Create groups, invite members, join requests, and group-specific feeds.
- **Events**: Group-based event creation with attendee tracking ("going" / "not going").
- **Real-time Chat**: 1:1 direct messaging and group chats powered by WebSockets.
- **Notifications**: Real-time alerts for messages, follow requests, and group invitations.

---

## 🛠 Tech Stack

| Layer | Technologies |
| :--- | :--- |
| **Frontend** | Next.js 15 (App Router), React 19, Context API, CSS Modules |
| **Backend** | Go 1.25, `net/http`, Gorilla WebSocket |
| **Database** | SQLite (via `mattn/go-sqlite3`) |
| **Migration** | `golang-migrate` |
| **Auth** | Session-based (HTTP-only cookies, bcrypt) |
| **Deployment** | Docker, Docker Compose |

---

## 📂 Project Structure

```text
.
├── backend/            # Go source code, migrations, and SQLite DB
│   ├── cmd/            # Application entry point
│   ├── internal/       # Business logic (service, repo, handler layers)
│   ├── migrations/     # SQL migration files
│   └── media/          # Uploaded user images
├── frontend/           # Next.js source code
│   ├── app/            # App Router (pages and layouts)
│   ├── components/     # UI components
│   ├── contexts/       # Global React state (Auth, WS, Notifications)
│   └── lib/            # API wrappers and utilities
└── docker-compose.yml  # Full-stack orchestration
```

---

## 🚦 Getting Started

### 🐳 Using Docker (Recommended)
The easiest way to run the full stack (Frontend + Backend + Database) is using Docker Compose:

```bash
docker compose up --build
```
- **Frontend**: `http://localhost:3000`
- **Backend API**: `http://localhost:8000`

### 🛠 Manual Setup

#### 1. Backend
```bash
cd backend
go run .
```
*Note: Migrations run automatically on startup. The database will be created at `backend/internal/database/social.db`.*

#### 2. Frontend
```bash
cd frontend
npm install
npm run dev
```
*By default, the frontend looks for the API at `http://localhost:8000`.*

---

## 🔐 Authentication & Real-time

### Session Management
The system uses **HTTP-only cookies** for security. 
1. The backend sets a `session_token` cookie upon login/registration.
2. The frontend `fetcher.js` includes `credentials: "include"` in every request.
3. The `AuthContext` on the frontend validates the session on page load via `/api/v1/auth/me`.

### WebSocket Communication
Real-time features (Chat, Online Status, Notifications) are handled via a single WebSocket connection:
- **Connection**: `ws://localhost:8000/ws`
- **State**: Managed by `WebSocketContext.jsx` in the frontend, allowing any component to `subscribe` to specific server events or `send` messages.

---

## 📡 Key API Endpoints Summary

| Domain | Key Endpoints |
| :--- | :--- |
| **Auth** | `POST /auth/register`, `POST /auth/login`, `GET /auth/me` |
| **Users** | `GET /user/profile`, `POST /user/follow`, `GET /user/suggestions` |
| **Posts** | `GET /posts/feed`, `POST /posts`, `POST /posts/comments` |
| **Groups** | `GET /groups/all`, `POST /groups/invite`, `GET /groups/events` |
| **Chat** | `GET /chat/conversations`, `GET /chat/direct`, `GET /chat/group` |
| **WS** | `/ws` (Direct Messages, Notifications, Online Status) |

---

## 📝 Development Notes

- **Ports**: Backend runs on `8000`, frontend on `3000` (dev and Docker).
- **Styling**: The frontend uses plain CSS files located under `frontend/css/` to ensure full control over the design without external UI libraries.
- **Image Storage**: Images are stored locally in the `backend/media` folder and served via the `/media/` static route.

For more detailed information, please refer to the individual `README.md` files in the `/frontend` and `/backend` directories.

