# Gator RSS Aggregator

`gator` is a small Go command-line RSS aggregator backed by PostgreSQL. It lets you:

- create and switch between users
- register RSS feeds
- follow and unfollow feeds
- continuously scrape feeds into a local database
- browse recent posts from the feeds you follow

The codebase is intentionally compact, but it still has a clear split between:

- CLI command handlers in `internal/commands`
- config loading in `internal/config`
- RSS fetching and XML parsing in `internal/rss`
- SQL schema and queries in `sql/`
- sqlc-generated database access code in `internal/database`

## Prerequisites

You need these installed before running the project:

- PostgreSQL
- Go

PostgreSQL stores users, feeds, follows, and scraped posts.
Go is required both to build the program and to install the CLI with `go install`.

You can verify your environment with:

```bash
psql --version
go version
```

The repository currently declares its Go toolchain in [go.mod](/home/dcbto/dev/github.com/RssAggregator/go.mod:1).

## Install the CLI

If you want to install the CLI from the current checkout:

```bash
go install .
```

If you want to install it directly from GitHub:

```bash
go install github.com/davicbtoliveira/rss_aggregator@latest
```

Important note: because the main package lives at the module root, `go install` will typically produce a binary named `rss_aggregator`, not `gator`.

If you want to run it as `gator`, create a small symlink:

```bash
ln -sf "$(go env GOPATH)/bin/rss_aggregator" "$(go env GOPATH)/bin/gator"
```

After that, the examples below can be run exactly as written. If you skip the symlink, replace `gator` with `rss_aggregator`.

## Running From Source

You can also run the project without installing it:

```bash
go run . <command> [args...]
```

Example:

```bash
go run . users
```

## Database Setup

Create a PostgreSQL database for the app:

```sql
CREATE DATABASE gator;
```

Example connection string:

```text
postgres://postgres:postgres@localhost:5432/gator?sslmode=disable
```

### Apply the schema

This project stores migrations in `sql/schema/`. There is no built-in migration command in the application, so you need to apply the SQL before using the CLI.

One simple approach is to run each schema file with `psql`:

```bash
for file in sql/schema/*.sql; do
  psql "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable" -f "$file"
done
```

The schema creates these tables:

- `users`
- `feeds`
- `feed_follows`
- `posts`

The `feeds` table also tracks `last_fetched_at`, which the aggregator uses to pick the next feed to scrape.

## Config File Setup

Before the CLI can run correctly, create a config file at:

```text
~/.gatorconfig.json
```

This is required. The code reads this file on startup and also writes the current logged-in user back into it. The file must already exist.

Start with:

```json
{
  "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

Fields:

- `db_url`: PostgreSQL connection string used on startup
- `current_user_name`: the active user for commands that act on behalf of a logged-in user

## Quick Start

Once PostgreSQL is running, the schema is applied, and `~/.gatorconfig.json` exists:

```bash
gator register alice
gator addfeed "Boot.dev Blog" "https://blog.boot.dev/index.xml"
gator agg 30s
```

In another terminal, browse posts:

```bash
gator browse 10
```

## Commands

The app registers commands in [main.go](/home/dcbto/dev/github.com/RssAggregator/main.go:18). These are the commands available today.

### `register <username>`

Creates a new user in the database and sets that user as the current user in `~/.gatorconfig.json`.

Example:

```bash
gator register alice
```

Behavior:

- inserts a row in `users`
- updates `current_user_name` in the config file
- fails if the username already exists

### `login <username>`

Switches the active user in `~/.gatorconfig.json`.

Example:

```bash
gator login alice
```

Behavior:

- checks that the user exists
- updates `current_user_name` in the config file

### `users`

Lists all users in the database and marks the current one.

Example:

```bash
gator users
```

Typical output:

```text
* alice (current)
* bob
```

### `reset`

Deletes all users from the database.

Example:

```bash
gator reset
```

Because the `feeds` and `feed_follows` tables use foreign keys with `ON DELETE CASCADE`, removing users will also remove their related data.

### `addfeed <name> <url>`

Creates a feed and automatically follows it as the current user.

Example:

```bash
gator addfeed "Hacker News" "https://hnrss.org/frontpage"
```

Behavior:

- requires a current logged-in user
- inserts a row into `feeds`
- immediately creates a row in `feed_follows`

This is a useful shortcut because you usually want to follow a feed as soon as you add it.

### `feeds`

Lists all feeds currently stored in the database.

Example:

```bash
gator feeds
```

For each feed, the command prints:

- feed name
- feed URL
- owning user name

### `follow <feed_url>`

Follows an existing feed by URL.

Example:

```bash
gator follow "https://hnrss.org/frontpage"
```

Behavior:

- requires a current logged-in user
- looks up the feed by URL
- creates a `feed_follows` row

Important: `follow` expects the feed URL, not the feed name.

### `following`

Lists the names of feeds followed by the current user.

Example:

```bash
gator following
```

### `unfollow <feed_url>`

Stops following a feed by URL.

Example:

```bash
gator unfollow "https://hnrss.org/frontpage"
```

Behavior:

- requires a current logged-in user
- deletes the matching row from `feed_follows`

### `agg <duration>`

Starts the long-running feed scraper.

Example:

```bash
gator agg 30s
gator agg 1m
gator agg 5m
```

Behavior:

- runs forever on a ticker
- selects the next feed ordered by `last_fetched_at`
- fetches the RSS XML over HTTP
- parses items from the feed
- inserts posts into the `posts` table
- skips duplicate posts when the same post URL already exists

This is the command that actually populates your database with content. Until `agg` runs, `browse` will have little or nothing to show.

### `browse [limit]`

Shows recent posts for the current user from followed feeds.

Examples:

```bash
gator browse
gator browse 5
gator browse 20
```

Behavior:

- defaults to `2` posts if no limit is provided
- joins `posts` with `feed_follows`
- only shows posts from feeds followed by the current user
- orders by `published_at DESC`

## Typical Workflow

The project is easiest to understand as a loop:

1. Create a user with `register`.
2. Add or follow feeds with `addfeed` or `follow`.
3. Run `agg <duration>` to continuously scrape RSS items into PostgreSQL.
4. Run `browse [limit]` to read the newest posts from followed feeds.

Example session:

```bash
gator register alice
gator addfeed "Boot.dev Blog" "https://blog.boot.dev/index.xml"
gator addfeed "Hacker News" "https://hnrss.org/frontpage"
gator following
gator agg 30s
```

In another terminal:

```bash
gator browse 10
```

## How It Works Internally

### Startup flow

On startup, [main.go](/home/dcbto/dev/github.com/RssAggregator/main.go:14):

1. reads `~/.gatorconfig.json`
2. opens a PostgreSQL connection
3. creates a sqlc query wrapper with `database.New(db)`
4. builds a shared `State` containing the config and query layer
5. registers command handlers
6. dispatches based on `os.Args`

### Command dispatch

[internal/commands/base_command.go](/home/dcbto/dev/github.com/RssAggregator/internal/commands/base_command.go:7) stores handlers in a `map[string]func(*State, Command) error`.

This keeps the CLI simple:

- parse the command name
- look up the handler
- execute it

### Logged-in middleware

Some commands need a real current user. Those commands are wrapped by [internal/middleware/loggedin.go](/home/dcbto/dev/github.com/RssAggregator/internal/middleware/loggedin.go:10), which:

- reads `current_user_name` from the config
- fetches the corresponding `users` row
- passes the resolved user into the real handler

This is used for:

- `addfeed`
- `follow`
- `unfollow`

### RSS fetching

[internal/rss/fetchfeed.go](/home/dcbto/dev/github.com/RssAggregator/internal/rss/fetchfeed.go:10) is responsible for downloading and parsing RSS XML.

Key details:

- sends an HTTP request with the `User-Agent` header set to `gator`
- unmarshals XML into `RSSFeed`
- unescapes HTML entities in channel and item text

The parser currently maps:

- channel title
- channel link
- channel description
- item title
- item link
- item description
- item publication date

### Aggregation loop

[internal/commands/agg.go](/home/dcbto/dev/github.com/RssAggregator/internal/commands/agg.go:15) drives scraping.

On each cycle it:

1. picks the next feed to fetch
2. marks it as fetched
3. downloads the RSS feed
4. loops through each item
5. parses the publication time when possible
6. inserts a post row
7. ignores duplicate-post errors caused by the unique constraint on `posts.url`

### Data model

The schema in `sql/schema/` defines four core tables:

- `users`: application users
- `feeds`: RSS feeds added by users
- `feed_follows`: many-to-many relationship between users and feeds
- `posts`: scraped RSS items

Relationships:

- each feed belongs to one user
- users can follow many feeds
- feeds can be followed by many users
- each post belongs to one feed

## Project Layout

```text
.
├── main.go
├── internal/
│   ├── commands/      # CLI handlers
│   ├── config/        # ~/.gatorconfig.json read/write
│   ├── database/      # sqlc-generated query code
│   ├── middleware/    # logged-in user middleware
│   └── rss/           # RSS fetch + XML types
├── sql/
│   ├── queries/       # sqlc query definitions
│   └── schema/        # database schema / migrations
└── sqlc.yaml          # sqlc generation config
```

## Development Notes

- SQL query definitions live in `sql/queries/`.
- Generated Go query code lives in `internal/database/`.
- The project uses PostgreSQL through `github.com/lib/pq`.
- UUIDs are generated in application code with `github.com/google/uuid`.

If you change the schema or query files, regenerate the database layer with `sqlc`.

## Troubleshooting

### `~/.gatorconfig.json not found`

Create the file manually before running the CLI. The app expects it to exist.

### `pq: database ... does not exist`

Create the database first in PostgreSQL, then update `db_url` in `~/.gatorconfig.json`.

### `sql: no rows in result set`

This usually means one of these:

- you tried to `login` as a user that does not exist
- you tried to `follow` or `unfollow` a feed URL that has not been added yet
- you tried to run a user-specific command before setting `current_user_name`

### `browse` shows nothing

Check all of the following:

- you are following at least one feed
- `agg` has been run long enough to fetch posts
- the RSS feeds you follow actually contain items

## Summary

To run this project successfully:

1. Install PostgreSQL and Go.
2. Install the CLI with `go install`.
3. Create the database.
4. Apply the SQL schema in `sql/schema/`.
5. Create `~/.gatorconfig.json`.
6. Register a user, add some feeds, run `agg`, and browse posts.
