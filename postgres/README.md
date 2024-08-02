# Postgres

This contains the Postgres Docker Compose definition and Schema migrations.

## Schema Migration

Migrations are created and run using [golang-migrate](https://github.com/golang-migrate/migrate).

As of now, `migrate.go` is pretty inflexible, and is provided as a convenience to not require installing and using the migrate CLI. It just runs the equivalent of `migrate -database YOUR_DATABASE_URL -path PATH_TO_YOUR_MIGRATIONS up`. If a migration fails or needs to be cleaned up for any reason, it will need manual intervention.
