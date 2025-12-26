## PIPELINE Restorely complet + modèle de données

### Cycle complet automatisé
| Config DB -> Backup -> Stockage -> Restore test -> Statut clair

```arduino
Scheduler
   │
   ▼
Backup Job (daily)
   │
   ├─ pg_dump
   ├─ encrypt
   ├─ upload S3
   └─ save status
          │
          ▼
Restore Test Job (weekly)
   │
   ├─ download
   ├─ decrypt
   ├─ docker restore
   └─ save result
```

Deux types de jobs:
- Backup (fréquent)
- Restore test (moins fréquent)

### Modèle de données
Organization
```sql
id
name
created_at
```

databases
```sql
id
organization_id
type            -- postgres
dsn_encrypted   -- chiffré
created_at
```

backups
```sql
id
database_id
status          -- success / failed
size_bytes
duration_ms
s3_path
created_at
```

restore_test
```sql
id
backup_id
status          -- success / failed
duration_ms
error_message
created_at
```

### État possibles
Backup
- pending
- success
- failed

Restore test
- pending
- success
- failed

### Scheduling
Cron interne en GO

Dans le worker:
- Un ticket toutes les X minutes
- Vérifie quoi lancer

Exemple:
- Backup: toute les 24h
- Restore test: tous les 7 jours

#### Pseudo-code Scheduler
```go
for {
    runPendingBackups()
    runPendingRestoreTests()
    time.Sleep(1 * time.Minutes)
}
```

### Flow Backup
Étapes
1. Créer ligne `backup (pending)`
2. Dump PostgresSQL
3. Chiffrer
4. Upload S3
5. Update `backup (success|failed)`

#### Pseudo-code
```go
backupID := createBackup()
dumpPath := dumpPostgres(db)
encPatn := encrypt(dumpPath)
s3Path := upload(encPath)

markBackupSuccess(backupID, s3Path)
```

### Flow RESTORE Test
Étapes
1. Trouver le dernier backup (success)
2. Créer ligne `restore_test (pending)`
3. Download dump
4. Déchiffrer
5. Restore docker
6. Update status

#### Pseudo-code
```go
backup := latestSuccessfulBackup(dbID)

retoreID := createRestoreTest(backup.ID)

dump := download(backup.s3Path)
plain := decrypt(dump)

err := testRestoreDocker(plain)

if err != nil {
    markRestoreFailed(restoreID, err)
} else {
    markRestoreSuccess(restoreID)
}
```

### Affichage UI
- LastBackup: OK/FAIL
- LastRestore: OK/FAIL
- Date
- Durée
