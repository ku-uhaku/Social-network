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
    'kuuhaku', 
    'kuuhaku@email.com', 
    'Sora', 
    'Shiro', 
    'male', 
    '2000-01-01', 
    1, 
    '$2a$10$OQ8d7gdfOe5uy3tCm1u1leT9RP0EOhsrq8k4HdtuBSUoE.EBCzutK', -- Plaintext: SecurePassword123!
    'https://api.dicebear.com/7.x/avataaars/svg?seed=kuuhaku', 
    'Master of all games. Built with Go and SQLite.'
);

INSERT INTO users (username, email, first_name, last_name, gender, date_of_birth, is_public, password, avatar, about_me)
VALUES 
(
    'yassine_b', 
    'yassine.bennani@email.com', 
    'Yassine', 
    'Bennani', 
    'male', 
    '1994-04-12', 
    1, 
    '$2a$10$OQ8d7gdfOe5uy3tCm1u1leT9RP0EOhsrq8k4HdtuBSUoE.EBCzutK', -- Plaintext: SecurePassword123!
    'https://api.dicebear.com/7.x/avataaars/svg?seed=Yassine', 
    'Software developer from Casablanca, passionate about backend engineering and sports.'
),
(
    'amina_alami', 
    'amina.alami@email.com', 
    'Amina', 
    'Alami', 
    'female', 
    '1998-11-23', 
    1, 
    '$2a$10$OQ8d7gdfOe5uy3tCm1u1leT9RP0EOhsrq8k4HdtuBSUoE.EBCzutK', -- Plaintext: SecurePassword123!
    'https://api.dicebear.com/7.x/avataaars/svg?seed=Amina', 
    'Graphic designer based in Rabat. Love art, photography, and Moroccan architecture.'
),
(
    'mehdi_tazi', 
    'mehdi.tazi@email.com', 
    'Mehdi', 
    'Tazi', 
    'male', 
    '1991-07-05', 
    1, 
    '$2a$10$OQ8d7gdfOe5uy3tCm1u1leT9RP0EOhsrq8k4HdtuBSUoE.EBCzutK', -- Plaintext: SecurePassword123!
    NULL, 
    'Project manager from Fez. Big fan of traditional Andalusian music and history.'
),
(
    'meriem_idrissi', 
    'm.idrissi@email.com', 
    'Meriem', 
    'Idrissi', 
    'female', 
    '2001-02-18', 
    1, 
    '$2a$10$OQ8d7gdfOe5uy3tCm1u1leT9RP0EOhsrq8k4HdtuBSUoE.EBCzutK', -- Plaintext: SecurePassword123!
    'https://api.dicebear.com/7.x/avataaars/svg?seed=Meriem', 
    'Data science student at EHTP. Exploring machine learning models.'
),
(
    'amza_mansouri', 
    'hamza.mansouri@email.com', 
    'Hamza', 
    'Mansouri', 
    'male', 
    '1996-09-30', 
    0, -- Private profile example
    '$2a$10$OQ8d7gdfOe5uy3tCm1u1leT9RP0EOhsrq8k4HdtuBSUoE.EBCzutK', -- Plaintext: SecurePassword123!
    NULL, 
    'UI/UX researcher located in Marrakech. Cyclist and travel enthusiast.'
);