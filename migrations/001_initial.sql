CREATE TABLE books (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  title TEXT NOT NULL,
  author TEXT,
  status TEXT NOT NULL CHECK (status IN ('READ','READING','NOT_STARTED','TO_BUY','ARCHIVED')),
  priority_to_buy INTEGER NOT NULL DEFAULT 0 CHECK (priority_to_buy IN (0, 1)),
  eligible_to_sell INTEGER NOT NULL DEFAULT 0 CHECK (eligible_to_sell IN (0, 1)),
  sold INTEGER NOT NULL DEFAULT 0 CHECK (sold IN (0, 1)),
  notes TEXT,
  added_at TEXT NOT NULL,
  started_at TEXT,
  finished_at TEXT
);

CREATE INDEX idx_books_status ON books(status);
CREATE INDEX idx_books_title ON books(title);
