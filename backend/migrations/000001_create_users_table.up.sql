PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    gender TEXT NOT NULL CHECK(gender IN ('male', 'female')),
    date_of_birth TEXT NOT NULL,  
    is_public INTEGER NOT NULL DEFAULT 1,
    password TEXT NOT NULL,
    
    avatar TEXT,                  
    about_me TEXT,
    
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

PRAGMA foreign_keys = ON;

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