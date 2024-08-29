CREATE TABLE IF NOT EXISTS subjects (
    id INT PRIMARY KEY,
    subject TEXT NOT NULL,
    words TEXT NOT NULL,
    used INT DEFAULT 0,
    correct INT DEFAULT 0
)
