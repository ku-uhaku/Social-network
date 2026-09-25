## A. Backend + infra

- [X] Backend Dockerfile → `golang:1.25.0-alpine`, `apk add --no-cache gcc musl-dev`, `CGO_ENABLED=1 go build` (sqlite needs CGO); no `echo/ls`/`cat` debug layers.
- [X] `backend/.dockerignore` → `media/`, `*.db`, `.git`; `frontend/.dockerignore` → `node_modules/`, `.next/`, `out/`, `build/`, `dist/`, `.git`.
- [X] `docker-compose.yml` → contexts `./backend` + `./frontend`, ports 8080/3000, `NEXT_PUBLIC_API_URL=http://localhost:8080`, db volume; one frontend CMD (`npm ci` + build + `npm start`).
- [X] `config/config.go` → port must be `8080` (compose + frontend expect it).
- [X] Re-implement the rate limiter in `helper/ratelimiter.go` (per-IP window, auth 20/min, api 300/min, skip `OPTIONS`) and wrap register/login/logout in `routes/auth.go` + the whole mux in `routes/route.go`.
- [X] Re-implement `IsValidImage` in `helper/helper.go` (decode config → jpeg/png/gif only → 8000px cap → full decode) + `imageExtensions`; make `SaveUploadedImage` validate the uploaded file (not a path string).
- [X] `requests/post.go` → title required/≤40, content required/≤200 (rune counts), comments same, privacy enum incl. `group`, private posts need `visible_to`, fix the privacy error wording.
- [X] `requests/group.go` → group title 3–20, description ≤200, event title ≤20, event description ≤100, with readable messages (drop “moore long”).
- [X] `requests/auth.go` → parse DOB, enforce one age threshold (same number front and back), keep nickname optional (regex only when provided).
- [X] `repository/user.go` → auto nickname `user_<8 hex>` when empty, `IsPrivate`.
- [X] `helper/helper.go` `DisplayName` + `lib/utils.js isAutoUsername` → auto nicknames never shown as `@user_xxx` in the UI.
- [X] `service/post.go checkPostVisibility` → owner ok, group posts members-only, private posts viewer-list-only, posts by a private user need an accepted follow; rename `is_private` → `isPrivate`.
- [X] `repository/post.go` feed query → a private user's public posts only for accepted followers.
- [X] `handler/user.go GetUserProfile` → return the profile, gate posts/details to owner/public/accepted follower, everything else = “This account is private” notice.
- [X] Profile page → details + posts only when `isOwner || is_public === 1 || follow_status === "accepted"`; no inverted `is_public === 0` logic.

## B. Features back

- [X] Post-in-group: allow `privacy='group'` in `requests/post.go` + the `000004` CHECK; `posts/create/page.jsx` reads `group_id` via `useSearchParams`, forces `privacy=group`, hides the privacy select, fixes the undefined `err` and the `console.log`.
- [X] Group chat backend: `repository/chat.go SaveGroupMessage` returning sender username/avatar, `GetGroupHistory`, membership gate in `service/chat.go`, handler + route `/api/v1/chat/group`, websocket `send_group_message` → `new_group_message`.
- [X] Migration `000013` group_message notification type (up + down) re-authored.
- [X] Group chat UI: `components/chat/GroupChat.jsx`, `lib/api/chat.js getGroupHistory`, group styles in `chat.css`, Chat button on the group page.
- [X] Unread badges: `contexts/GroupChatContext.jsx`, provider in `app/layout.jsx`, badges in `group/page.jsx`, `group/[id]/page.jsx`, `GroupCard.jsx` + `.groupChatBadge` CSS.
- [X] Notification list: `contexts/NotificationContext.jsx` filters out `group_message`.
- [X] Chat sounds: send/receive SFX dropped — no `playSfx` calls in `Chat.jsx` / `GroupChat.jsx`.
- [X] Mobile header menu: `app/(main)/layout.jsx` burger (`menuOpen`, 3 bars) + `.headerMenuToggle` / `.headerRight.open` rules in `responsive.css` + burger styles in `home.css`.
- [X] Register form: nickname marked optional, age guard with `isOldEnough`.
- [X] Error page: `app/error/page.jsx` + `ErrorDisplay.jsx` cleanup — no `console.log`/commented music; `/error?message=` handling dropped (zero callers, broken on Next 15).
- [X] `repository/group.go RemoveMember` → transaction: remove membership, delete the group when the creator leaves.


## C. Wrap up

- [X] Docs: single `readme.md` written from the actual code (backend + frontend docs folded in).
- [X] `npm install` to refresh `package-lock.json` (already in sync — reported "up to date").
- [X] Theme decisions: restored master's state — `globals.css` heading rule + `font-family: inherit` lines in `home.css` / `groups.css` commented out.
- [X] Final checks: `go build ./... && go vet ./...`, `npm run build` green; authors = Mbelhouss + mbarrah; junk greps clean; `.clinerules` audit reviewed. Docker run skipped per user.
