package repository

import (
	"context"

	"kuu/internal/models"
)

// SaveDirectMessage persists a 1-on-1 message
func (r *Repository) SaveDirectMessage(ctx context.Context, senderID, receiverID int64, content string) (*models.DirectMessage, error) {
	query := `
INSERT INTO direct_messages (sender_id, receiver_id, content)
VALUES ($1, $2, $3)
RETURNING id, sender_id, receiver_id, content, created_at
`
	var msg models.DirectMessage
	err := r.DB.Database.QueryRowContext(ctx, query, senderID, receiverID, content).
		Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

// SaveGroupMessage persists a group message and returns it joined with the
// sender's username/avatar so realtime recipients can render it immediately.
func (r *Repository) SaveGroupMessage(ctx context.Context, senderID, groupID int64, content string) (*models.GroupMessage, error) {
	query := ` INSERT INTO group_messages (group_id, sender_id, content)
VALUES ($1, $2, $3)
RETURNING id, group_id, sender_id, content, created_at,
          (SELECT username FROM users WHERE id = sender_id),
          (SELECT avatar FROM users WHERE id = sender_id)
`
	var msg models.GroupMessage
	err := r.DB.Database.QueryRowContext(ctx, query, groupID, senderID, content).
		Scan(&msg.ID, &msg.GroupID, &msg.SenderID, &msg.Content, &msg.CreatedAt, &msg.Username, &msg.Avatar)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

// GetDirectHistory returns a page of DMs between two users, newest first.
func (r *Repository) GetDirectHistory(ctx context.Context, userA, userB int64, limit, offset int) ([]models.DirectMessage, error) {
	query := `
SELECT id, sender_id, receiver_id, content, created_at
FROM direct_messages
WHERE (sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1)
ORDER BY id DESC
LIMIT $3 OFFSET $4
`
	rows, err := r.DB.Database.QueryContext(ctx, query, userA, userB, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []models.DirectMessage
	for rows.Next() {
		var m models.DirectMessage
		if err := rows.Scan(&m.ID, &m.SenderID, &m.ReceiverID, &m.Content, &m.CreatedAt); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}
	return msgs, rows.Err()
}

// MarkChatRead records that the viewer has read the conversation up to its latest message.
func (r *Repository) MarkChatRead(ctx context.Context, viewerID, otherUserID int64) error {
	query := `
INSERT INTO chat_reads (user_id, other_user_id, last_read_message_id)
VALUES ($1, $2, COALESCE((
SELECT MAX(id) FROM direct_messages
WHERE (sender_id = $1 AND receiver_id = $2)
   OR (sender_id = $2 AND receiver_id = $1)
), 0))
ON CONFLICT(user_id, other_user_id) DO UPDATE SET
last_read_message_id = excluded.last_read_message_id,
updated_at = CURRENT_TIMESTAMP
`
	_, err := r.DB.Database.ExecContext(ctx, query, viewerID, otherUserID)
	return err
}

// ListConversations returns every chat-able user for viewerID (accepted follow in
// either direction) with the timestamp of their latest DM and their unread message
// count, ordered by most recent message first and alphabetically for conversations
// with no messages yet.
func (r *Repository) ListConversations(ctx context.Context, viewerID int64) ([]models.ConversationMetadata, error) {
	query := `
SELECT u.id, u.username, u.avatar, dm.created_at,
       (SELECT COUNT(*)
        FROM direct_messages du
        WHERE du.sender_id = u.id
          AND du.receiver_id = $1
          AND du.id > COALESCE(cr.last_read_message_id, 0))
FROM users u
LEFT JOIN direct_messages dm ON dm.id = (
SELECT dm2.id FROM direct_messages dm2
WHERE (dm2.sender_id = $1 AND dm2.receiver_id = u.id)
   OR (dm2.sender_id = u.id AND dm2.receiver_id = $1)
ORDER BY dm2.created_at DESC, dm2.id DESC
LIMIT 1
)
LEFT JOIN chat_reads cr ON cr.user_id = $1 AND cr.other_user_id = u.id
WHERE u.id != $1
  AND EXISTS (
SELECT 1 FROM follows f
WHERE f.status = 'accepted'
  AND ((f.follower_id = $1 AND f.following_id = u.id)
OR (f.follower_id = u.id AND f.following_id = $1))
  )
ORDER BY (dm.created_at IS NULL) ASC, dm.created_at DESC, u.username ASC
`
	rows, err := r.DB.Database.QueryContext(ctx, query, viewerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	conversations := []models.ConversationMetadata{}
	for rows.Next() {
		var c models.ConversationMetadata
		if err := rows.Scan(&c.UserID, &c.Username, &c.Avatar, &c.LastMessageAt, &c.UnreadCount); err != nil {
			return nil, err
		}
		conversations = append(conversations, c)
	}
	return conversations, rows.Err()
}

// GetGroupHistory retrieves paginated chat history for a group
func (r *Repository) GetGroupHistory(ctx context.Context, groupID int64, limit, offset int) ([]models.GroupMessage, error) {
	query := `
SELECT gm.id, gm.group_id, gm.sender_id, u.username, u.avatar, gm.content, gm.created_at
FROM group_messages gm
JOIN users u ON u.id = gm.sender_id
WHERE gm.group_id = $1
ORDER BY gm.created_at DESC
LIMIT $2 OFFSET $3
`
	rows, err := r.DB.Database.QueryContext(ctx, query, groupID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []models.GroupMessage
	for rows.Next() {
		var m models.GroupMessage
		if err := rows.Scan(&m.ID, &m.GroupID, &m.SenderID, &m.Username, &m.Avatar, &m.Content, &m.CreatedAt); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}
	return msgs, nil
}

// GetGroupMemberIDs returns all active user IDs in a group for event routing
func (r *Repository) GetGroupMemberIDs(ctx context.Context, groupID int64) ([]int64, error) {
	query := `SELECT user_id FROM group_members WHERE group_id = $1 AND status = 'accepted'`
	rows, err := r.DB.Database.QueryContext(ctx, query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
