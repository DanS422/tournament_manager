# README

## Troubleshooting

### Dirty migrations
It happens normally when migrations within one file are run partially successfully. To check the version
```
sqlite3 tournaments.db "SELECT * FROM schema_migrations;"
```

Result may be something like `4|1`, which means version 4 with dirty migration 1. To resolve it, update the schema by running 
```
sqlite3 tournaments.db "UPDATE schema_migrations SET version = 4, dirty = 0;"
```

And then run the first cmd again to verify the result.
