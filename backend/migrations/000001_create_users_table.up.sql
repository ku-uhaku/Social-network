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