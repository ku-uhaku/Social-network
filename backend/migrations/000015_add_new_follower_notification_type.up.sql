PRAGMA foreign_keys = ON;

-- Following a public account is instant, so it needs a notification type with
-- no accept/decline buttons. SQLite cannot alter a CHECK constraint, so we
-- rebuild the table.
ALTER TABLE notifications RENAME TO notifications_old;

CREATE TABLE notifications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    recipient_id INTEGER NOT NULL,
    actor_id INTEGER,
    type TEXT NOT NULL CHECK (type IN ('follow_request', 'new_follower', 'group_invitation', 'group_join_request', 'group_event_created', 'group_message')),
    title TEXT NOT NULL,
    message TEXT NOT NULL,
    is_read INTEGER NOT NULL DEFAULT 0,
    is_expired INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    group_id INTEGER,
    FOREIGN KEY (recipient_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (actor_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE
);

INSERT INTO notifications (id, recipient_id, actor_id, type, title, message, is_read, is_expired, created_at, group_id)
    SELECT id, recipient_id, actor_id, type, title, message, is_read, is_expired, created_at, group_id
    FROM notifications_old;

DROP TABLE notifications_old;

CREATE INDEX IF NOT EXISTS idx_notifications_recipient_created
    ON notifications(recipient_id, created_at DESC);
