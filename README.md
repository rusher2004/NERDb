# NERDb

Back at it. More soon..

## Data & Server

Built on data provided by [Zkillboard], [EVE Ref], and [ESI].

Killmail, player and other data is stored in [Postgres]. The schema and migrations managed in [postgres](/postgres/).

Data loading and updating is handled by the [killfeed](/killfeed/) package, written in [Go].

## Web app

Built with [Next].

## Running the project locally.

### Requirements

- [Docker]
- [Go]
- [Node]

> Note: The following instructions use `go run` commands, which is fine for now. Future revisions will have a compiled binary to execute them.
>
> For each command, more info about its usage can be found in its containing directory.

#### Start the database

From project root, run `docker compose up -d`

#### Create database schema

Run `go run postgres/migrate.go`

#### Start loading current Zkillboard data

Run `go run killfeed/main.go zkill`

#### Load historical killmail data from Eve Ref

Run `go run killfeed/main.go everef`

> Expect this step to take hours. It is intentionally not loading days concurrently, so that there is not a heavy load on the network or database.

#### Load character, corporation, and alliance names

Killmail data contains only IDs, and CCP's ESI API does not provide bulk operations to get multiple entity's info. This is a two step process to load bulk data, and then update more current data.

1. Load bulk data from Eve Ref
   1. Download the [latest bulk dataset](https://data.everef.net/characters-corporations-alliances/backfills/)
   2. Unzip the file
   3. Run `go run killfeed/main.go updater --src everef --dir $PATH_TO_UNZIPPED_FILE`
2. Run ESI updater for the remaining values
   1. Run `go run killfeed/main.go updater --src esi`

> TODO:
>
> These processes will be updated to also update corporations and alliances.
>
> The web app will update values on demand as users view entities in the app, respecting ESI's caching.

[Docker]: https://www.docker.com
[ESI]: https://esi.evetech.net/ui/
[Eve Ref]: https://docs.everef.net/datasets/
[Go]: https://go.dev
[Next]: https://nextjs.org "Next.js"
[Node]: https://nodejs.org/en "Node.js"
[Postgres]: https://www.postgresql.org
[Zkillboard]: https://github.com/zKillboard/zKillboard/wiki
