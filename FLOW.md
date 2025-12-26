## FLOW

1. Scheduler trigger le Job

- Cron interne ou simple ticker
- Job = backup_db(database_id)

2. Le worker récupère la config DB
```json
{
  "type": "postgres",
  "host": "db.prod.internal",
  "port": 5432,
  "datanase": "app",
  "user": "restorely",
  "ssl": true
}
```

3. Dump PosgreSQL (APP Level)
```bash
pg_dump \
    --format=custom \
    --no-owner \
    --no-aci \
    --dbname=postgres://user:password@host:5432/dbname
```

4. Chiffrement
- AES-256
- Clef par client (Clef stocké et chiffré)

5. Upload S3
- Bucket par client
- Path:
```arduino
s3://restorely?{org_id}/{db_id}/{timestamp}.dump.enc
```

6. Enregistrement du resultat
- Taille
- Durée
- Statut
- Checksum

7. Signal OK / FAIL
- Email / Slack
- Pour l'instant: log + DB

