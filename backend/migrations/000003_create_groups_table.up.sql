PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL UNIQUE,
    description TEXT,
    is_public INTEGER NOT NULL DEFAULT 1,
    creator_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (creator_id) REFERENCES users(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS group_members (
    user_id INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    status TEXT NOT NULL DEFAULT 'accepted',
    joined_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (user_id, group_id),
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_group_members_group_id ON group_members(group_id);

INSERT INTO groups (title, description, creator_id, is_public)
VALUES 
(
    'Hallownest Explorers', 
    'A collaborative group for wanderers, cartographers, and lore-hunters exploring every corner of Hallownest—from the Forgotten Crossroads to the deepest pits of Deepnest.', 
    1, 
    1 -- Public Group
),
(
    'Dreamers & Nightmares', 
    'Delving into the dreams of the Dreamers and the nightmares lurking beneath. Share your encounters, theories, and protective charms with fellow dreamers.', 
    2, 
    1 -- Public Group
),
(
    'Colosseum of Fools', 
    'Private training circle for warriors preparing for the trials of the Colosseum. Share nail techniques, battle strategies, and bragging rights.', 
    4, 
    0 -- Private Group
);

INSERT INTO group_members (user_id, group_id, status, joined_at)
VALUES 
(1, 1, 'accepted', CURRENT_TIMESTAMP),
(2, 1, 'accepted', CURRENT_TIMESTAMP),
(2, 2, 'accepted', CURRENT_TIMESTAMP),
(3, 1, 'accepted', CURRENT_TIMESTAMP),
(4, 3, 'accepted', CURRENT_TIMESTAMP),
(5, 1, 'accepted', CURRENT_TIMESTAMP);