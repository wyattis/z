CREATE TABLE comment (
  id INTEGER PRIMARY KEY,
  user_id INTEGER,
  content TEXT NOT NULL,
  created_at TEXT NOT NULL
);