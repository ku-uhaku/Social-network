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
    'Morocco Tech Hub', 
    'A collaborative group for developers, designers, and tech innovators across Morocco to share ideas, code, and networking opportunities.', 
    1, 
    1 -- Public Group
),
(
    'Marrakech Travel & Photo', 
    'Exploring hidden gems, historical riads, and landscape photography around the red city and the Atlas mountains.', 
    1, 
    1 -- Public Group
),
(
    'Casablanca Startups', 
    'Private mastermind group tracking new digital startup operations, venture capitals, and business optimization strategies in Casablanca.', 
    1, 
    0 -- Private Group
);

INSERT INTO group_members (user_id, group_id, status, joined_at)
VALUES 
(1, 1, 'accepted', CURRENT_TIMESTAMP),
(1, 2, 'accepted', CURRENT_TIMESTAMP),
(1, 3, 'accepted', CURRENT_TIMESTAMP);