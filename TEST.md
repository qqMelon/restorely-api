## DATA TEST

Création de donnée de test réaliste
```sql
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  email TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT now()
);

INSERT INTO users (email)
SELECT 'user' || g || '@example.com'
FROM generate_series(1, 100) g;

CREATE TABLE orders (
  id SERIAL PRIMARY KEY,
  user_id INT REFERENCES users(id),
  amount INT,
  created_at TIMESTAMP DEFAULT now()
);

INSERT INTO orders (user_id, amount)
SELECT (random() * 100)::int + 1, (random() * 500)::int
FROM generate_series(1, 200);
```

Créer un dump:
```bash
docker exec source-base \
  pg_dump \
  -U postgres \
  --format=custom \
  --no-owner \
  --no-acl \
  appdb > backup.dump
```

Copier le dump dans le container de test de restauration
```bash
docker cp backup.dump restore-base:/backup.dump
```

Lancer la restauration
```bash
docker exec -it restore-base \
  pg_restore \
  -U postgres \
  --no-owner \
  --no-acl \
  -d restoredb \
  /backup.dump
```
