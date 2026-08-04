PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    image_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

INSERT INTO comments (post_id, user_id, title, content, image_url)
VALUES 
(
    1, 
    3, 
    'A familiar descent', 
    'I remember my first fall from those cliffs. Hallownest keeps its secrets well—but it also keeps its wonders. Take your time, little ghost, and listen to the echoes.', 
    NULL
),
(
    1, 
    6, 
    'Welcome to Dirtmouth', 
    'New face! And a paying one, I hope. When you come back up from the well, stop by my shop. I''ll even offer a first-time discount. Barely.', 
    NULL
),
(
    2, 
    3, 
    'A tip for geo farming', 
    'The husks near the hot springs drop more geo than the crossing husks. If you have a decent nail and some patience, you can also take down the armored sentries near the old temple entrance.', 
    NULL
),
(
    2, 
    6, 
    'Or you could buy a map', 
    'I sell maps too, you know. Slightly used, very slightly outdated, but I''m the only one who has them. Geo invested in knowledge is geo well spent.', 
    NULL
),
(
    3, 
    1, 
    'I hear the skittering', 
    'I know that sound. I''ve been through the Weaver''s Den once. I have no intention of a second visit. Your warning is heeded.', 
    NULL
),
(
    3, 
    5, 
    'Deepnest gives me the shivers', 
    'The Crystal Peak is scary enough for me! Please stay safe down there. The singing in Deepnest is... not the nice kind of singing.', 
    NULL
),
(
    4, 
    1, 
    'I found this place too', 
    'I sat by the shore for a long time. The stillness was almost overwhelming. Thank you for putting into words what I could not.', 
    NULL
),
(
    4, 
    6, 
    'Beautiful spot, but no shops nearby', 
    'Scenic, I admit. Terrible for business, though. The view doesn''t buy geo.', 
    'https://placehold.co/600x400?text=Sly+Approves'
),
(
    5, 
    2, 
    'The Watcher''s records', 
    'The archives you speak of are watched. The Watcher''s Spire is not just a monument—it is a prison. Be careful which old stones you turn, scholar.', 
    NULL
),
(
    6, 
    1, 
    'Blinded by my own magnificence', 
    'Zote, the mirror in the City of Tears already cracked when you walked past it. I think you have nothing to fear on that front.', 
    NULL
),
(
    6, 
    3, 
    'I have seen the reflection', 
    'I have met many warriors in my travels, but none quite like you, Zote. I mean that in the most literal sense.', 
    NULL
),
(
    7, 
    1, 
    'You tripped over the floor?', 
    'The floor was right there, Zote. It has always been right there.', 
    NULL
),
(
    8, 
    1, 
    'The humming is growing louder', 
    'Myla, the song you hear—are you sure it is the crystals? Sometimes the lights in the Peak whisper. If it changes tune, please come back to Dirtmouth.', 
    NULL
),
(
    8, 
    3, 
    'I have heard of such songs', 
    'In my studies I read of minerals that sing. The Scholars of the Soul Sanctum believed it was a form of memory. Beautiful, but I fear it for you, Myla.', 
    NULL
),
(
    9, 
    4, 
    'Everything? Even your nail?', 
    'Sly, you would sell your own nail if the price were right. I might buy it just to feel what a Nailmaster''s sword is like.', 
    NULL
),
(
    10, 
    3, 
    'The Old Nail''s history', 
    'It is not merely rusted—it belonged to a Soul Master who fell to vanity. Some weapons carry their wielder''s fate. Yours, I trust, only carries profit.', 
    NULL
);