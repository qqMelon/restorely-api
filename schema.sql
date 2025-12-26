CREATE TABLE databases (
  id UUID PRIMARY KEY,
  name TEXT NOT NULL,
  type TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE backups (
  id UUID PRIMARY KEY,
  database_id UUID REFERENCES databases(id),
  status TEXT NOT NULL,
  duration_ms BIGINT,
  created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE restore_tests (
  id UUID PRIMARY KEY,
  backup_id UUID REFERENCES backups(id),
  status TEXT NOT NULL,
  duration_ms BIGINT,
  created_at TIMESTAMP DEFAULT now()
);

