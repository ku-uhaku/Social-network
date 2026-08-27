PRAGMA foreign_keys = ON;

-- Replace the payload/actions JSON blobs with the one value they carried.
-- 'actions' only ever held button labels, which the type already determines,
-- and 'payload' only ever held group_id: event_id and target_user_id were
-- written but never read.
ALTER TABLE notifications ADD COLUMN group_id INTEGER REFERENCES groups(id) ON DELETE CASCADE;

UPDATE notifications SET group_id = json_extract(payload, '$.group_id')
    WHERE payload IS NOT NULL;

ALTER TABLE notifications DROP COLUMN payload;
ALTER TABLE notifications DROP COLUMN actions;
