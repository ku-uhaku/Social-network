PRAGMA foreign_keys = ON;

ALTER TABLE notifications ADD COLUMN payload TEXT;
ALTER TABLE notifications ADD COLUMN actions TEXT;

UPDATE notifications SET payload = json_object('group_id', group_id)
    WHERE group_id IS NOT NULL;

UPDATE notifications SET actions = json_object(
        'buttons', json_array(
            json_object('action', 'accept', 'label', 'Accept'),
            json_object('action', 'decline', 'label', 'Decline')
        )
    )
    WHERE type IN ('follow_request', 'group_invitation', 'group_join_request');

UPDATE notifications SET actions = json_object(
        'buttons', json_array(json_object('action', 'view', 'label', 'View event'))
    )
    WHERE type = 'group_event_created';

ALTER TABLE notifications DROP COLUMN group_id;
