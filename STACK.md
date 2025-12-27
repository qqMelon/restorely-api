## Restorely complet stack

Objectif:
- Une DB Restorely persistente
- Une API Go
- Un worker Go (backup + restore test)
- Orchestration temporaire par docker-compose

### Architecture finale (local)
```scss
┌─────────────┐
│   Frontend  │ (Vue)
│  localhost  │
└──────┬──────┘
       │ HTTP
┌──────▼──────┐
│     API     │ (Go)
│ :8080       │
└──────┬──────┘
       │ SQL
┌──────▼──────┐
│ PostgreSQL  │  ← source de vérité
│ restorely   │
└──────┬──────┘
       │ SQL
┌──────▼──────┐
│   Worker    │ (Go)
│ backups +   │
│ restore     │
└─────────────┘
```
