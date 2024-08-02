# killfeed

Command line tools to fetch ESI and killmail data.

## Data Sources

### Killmails

Live updates from zKill are fetched using the [listener](/killfeed/listener) package.

Historical data is loaded from EveRef with the [EveRef Engine](/killfeed/everef/engine.go).

### ESI Data

Bulk character, corporation, and alliance data is loaded from [EveRef's dataset](https://data.everef.net/characters-corporations-alliances/) using the engine's RunPlayerUpdater method. To fill the gaps for more recent data, we load directly from ESI with the [updater package](/killfeed/updater/)'s UpdateCharacter method.
