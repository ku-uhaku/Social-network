PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    group_id INTEGER, -- NULL for user posts, set for group posts
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    privacy TEXT NOT NULL DEFAULT 'public' CHECK (privacy IN ('public', 'almost private', 'private')),
    image_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(group_id) REFERENCES groups(id) ON DELETE CASCADE
);

INSERT INTO posts (user_id, group_id, title, content, privacy, image_url)
VALUES 
(
    1, 
    NULL, 
    'Just fell into Hallownest from the Howling Cliffs', 
    'The air down here is heavy and the light is strange. I found an old town called Dirtmouth—everyone is either resting or gone. Also, an old bug at the well told me not to go below. I''m going below. Wish me luck.', 
    'public', 
    NULL
),
(
    1, 
    1, 
    'Found the Forgotten Crossroads—need geo for a map', 
    'Cornifer''s map is invaluable but a bit pricey. Anyone know the fastest way to farm geo? Also, beware the husks near the hot springs—they hit harder than they look.', 
    'public', 
    NULL
),
(
    2, 
    NULL, 
    'To those who wander Deepnest', 
    'Heed this warning, little ones: tread softly in my mother''s domain. The Weaver''s songs are not for all ears, and the Stalking Devouts are not forgiving. If you hear skittering, it is already too late.', 
    'public', 
    NULL
),
(
    3, 
    NULL, 
    'A rest at the Blue Lake', 
    'There is a stillness up here that Hallownest below has forgotten. The water mirrors the void above. I have set down my nail for a moment just to sit. Sometimes the greatest discovery is a quiet place to think.', 
    'public', 
    'https://placehold.co/600x400?text=Blue+Lake'
),
(
    3, 
    1, 
    'Notes on the ancient civilisation beneath the City of Tears', 
    'I have spent many days pouring over the stone tablets beneath the City. The Great Knights left records of their oath—and of their failures. If you seek the truth of Hallownest''s fall, start with the old archives near the Watcher''s Spire.', 
    'public', 
    NULL
),
(
    4, 
    NULL, 
    'My 56th Precept: Do not play with mirrors', 
    'The 56th Precept is essential for any warrior: never play with mirrors. For one, the reflection may not be you. For two, you may accidentally discover you are more glorious than you realised, and be blinded by your own magnificence.', 
    'public', 
    NULL
),
(
    4, 
    3, 
    'The Colosseum of Fools cannot contain my might', 
    'I, Zote the Mighty, have vanquished every warrior the Colosseum has thrown at me. Or I will, once I finish writing this post. The Trial of the Fool is merely a rehearsal for the day I face the true final boss: the floor after a misstep.', 
    'public', 
    NULL
),
(
    5, 
    NULL, 
    'La la la! The crystals are singing today!', 
    'I found a new vein of crystal in the upper reaches of the Peak. It glows so bright and hums so sweetly. I keep singing along. My pickaxe has a nice rhythm today. La la la...', 
    'almost private', 
    NULL
),
(
    6, 
    NULL, 
    'For sale: everything in Dirtmouth', 
    'Grand reopening of my shop! Nails, charms, maps, dubious relics—all for the right price. And for the brave: I have a secret stash of goods in my old home beneath the town. Knock thrice and mention Sly sent you.', 
    'public', 
    'https://placehold.co/600x400?text=Sly%27s+Shop'
),
(
    6, 
    1, 
    'A warning about the nail that they call "Old Nail"', 
    'I have appraised many nails in my day. The Old Nail you find in the Soul Sanctum is rusted and brittle—better to save your geo and buy a proper nail from me. Trust the shopkeeper who survived the Nailsmith''s apprenticeship.', 
    'public', 
    NULL
);