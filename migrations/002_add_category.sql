ALTER TABLE books ADD COLUMN category TEXT CHECK (
  category IS NULL OR category IN (
    'THEOLOGY',
    'FICTION',
    'SOFTWARE',
    'PHILOSOPHY',
    'HISTORY',
    'PERSONAL_DEVELOPMENT',
    'FINANCE_BUSINESS',
    'SCIENCE',
    'POLITICS_CULTURE',
    'BIOGRAPHY',
    'OTHER'
  )
);

CREATE INDEX idx_books_category ON books(category);
