## Choix techniques clefs

### Pourquoi pg_dump et pas une lib GO

- pg_dump = battle-tested
- Gère versions PostgreSQL
- Zéro surprise

### Pourquoi pas d'agent ?
- Connexion TCP directe
- Zéro friction
- Persona Alex dit oui !

### Pourquoi pas de streaming ?
- MVP = simplicité (Dans un premier temps)
- Fichier temporaire OK
- Optimisation plus tard

