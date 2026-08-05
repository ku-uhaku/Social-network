-- Seed data for testing: users, groups, posts, comments, follows

-- Users
INSERT INTO users (
    id, 
    username, 
    email, 
    first_name, 
    last_name, 
    gender, 
    date_of_birth, 
    is_public, 
    password, 
    avatar, 
    about_me
) VALUES (
    1, 
    'the_knight', 
    'the.knight@hallownest.com', 
    'Ghost', 
    'Vessel', 
    'male', 
    '1950-01-01', 
    1, 
    '$2a$10$OQ8d7gdfOe5uy3tCm1u1leT9RP0EOhsrq8k4HdtuBSUoE.EBCzutK', -- Plaintext: SecurePassword123!
    'https://api.dicebear.com/7.x/avataaars/svg?seed=the_knight', 
    'The little ghost who fell into Hallownest. Armed with a broken nail and a whole lot of determination.'
);

INSERT INTO users (id, username, email, first_name, last_name, gender, date_of_birth, is_public, password, avatar, about_me)
VALUES 
(
    2, 
    'hornet', 
    'hornet@deepnest.com', 
    'Hornet', 
    'Protector', 
    'female', 
    '1965-06-15', 
    1, 
    '$2a$10$OQ8d7gdfOe5uy3tCm1u1leT9RP0EOhsrq8k4HdtuBSUoE.EBCzutK', -- Plaintext: SecurePassword123!
    'https://api.dicebear.com/7.x/avataaars/svg?seed=hornet', 
    'Guardian of Hallownest and spinner of silk. Also known as the Gendered Child. Do not test me.'
),
(
    3, 
    'quirrel', 
    'quirrel@blue.lake.com', 
    'Quirrel', 
    'Scholar', 
    'male', 
    '1958-03-22', 
    1, 
    '$2a$10$OQ8d7gdfOe5uy3tCm1u1leT9RP0EOhsrq8k4HdtuBSUoE.EBCzutK', -- Plaintext: SecurePassword123!
    'https://api.dicebear.com/7.x/avataaars/svg?seed=quirrel', 
    'Wandering scholar with a nail and a thirst for knowledge. Ask me about the Blue Lake.'
),
(
    4, 
    'zote_the_mighty', 
    'zote@colo.fools.com', 
    'Zote', 
    'Mighty', 
    'male', 
    '1952-11-08', 
    1, 
    '$2a$10$OQ8d7gdfOe5uy3tCm1u1leT9RP0EOhsrq8k4HdtuBSUoE.EBCzutK', -- Plaintext: SecurePassword123!
    'https://api.dicebear.com/7.x/avataaars/svg?seed=zote_the_mighty', 
    'The 57th Precept: Do not run to satisfy yourself. Glory awaits those who speak loudly and swing bigger.'
),
(
    5, 
    'myla', 
    'myla@crystal.peak.com', 
    'Myla', 
    'Miner', 
    'female', 
    '1972-09-05', 
    0, -- Private profile example
    '$2a$10$OQ8d7gdfOe5uy3tCm1u1leT9RP0EOhsrq8k4HdtuBSUoE.EBCzutK', -- Plaintext: SecurePassword123!
    'https://api.dicebear.com/7.x/avataaars/svg?seed=myla', 
    'Just a little miner digging for geo in Crystal Peak. La la la! The crystals are singing to me...'
),
(
    6, 
    'sly', 
    'sly@dirtmouth.shop.com', 
    'Sly', 
    'The Great Zote-Like Shopkeeper', 
    'male', 
    '1940-04-30', 
    1, 
    '$2a$10$OQ8d7gdfOe5uy3tCm1u1leT9RP0EOhsrq8k4HdtuBSUoE.EBCzutK', -- Plaintext: SecurePassword123!
    'https://api.dicebear.com/7.x/avataaars/svg?seed=sly', 
    'Former Nailmaster turned humble shopkeeper in Dirtmouth. Everything is for sale—except my secrets.'
);

-- Groups
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

-- Group members
INSERT INTO group_members (user_id, group_id, status, joined_at)
VALUES 
(1, 1, 'accepted', CURRENT_TIMESTAMP),
(2, 1, 'accepted', CURRENT_TIMESTAMP),
(2, 2, 'accepted', CURRENT_TIMESTAMP),
(3, 1, 'accepted', CURRENT_TIMESTAMP),
(4, 3, 'accepted', CURRENT_TIMESTAMP),
(5, 1, 'accepted', CURRENT_TIMESTAMP);

-- Posts
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

-- Comments
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

-- Follows
-- the_knight (1) follows everyone
INSERT INTO follows (follower_id, following_id, status) VALUES
    (1, 2, 'accepted'), -- the_knight -> hornet
    (1, 3, 'accepted'), -- the_knight -> quirrel
    (1, 4, 'accepted'), -- the_knight -> zote_the_mighty
    (1, 5, 'accepted'), -- the_knight -> myla (private, accepted for testing)
    (1, 6, 'accepted'); -- the_knight -> sly

-- Other relationships
INSERT INTO follows (follower_id, following_id, status) VALUES
    (2, 1, 'accepted'), -- hornet -> the_knight
    (2, 3, 'accepted'), -- hornet -> quirrel
    (3, 1, 'accepted'), -- quirrel -> the_knight
    (3, 2, 'accepted'), -- quirrel -> hornet
    (4, 1, 'accepted'), -- zote_the_mighty -> the_knight
    (5, 1, 'accepted'), -- myla -> the_knight
    (6, 1, 'accepted'); -- sly -> the_knight