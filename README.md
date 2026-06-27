# books-cli

A personal command-line tool to manage a reading list backed by SQLite.

Replaces a manual spreadsheet workflow with structured storage, fast updates, and scriptable queries. Built in Go as a single binary with no database server required.

## Goals

- Store books and metadata in a local SQLite database
- Manage the list from the terminal with a simple CLI
- Support fetch, filter, search, and update operations
- Keep the data model strict enough to avoid spreadsheet-style inconsistencies
- Leave room for a future MCP layer on top of the same database logic

## Language

**Go**

- Single portable binary
- Strong standard library (`database/sql`, `embed` for migrations)
- Good fit for a small personal CLI tool

## Database schema

Initial `books` table:

| Column | Type | Notes |
| --- | --- | --- |
| `id` | `INTEGER PRIMARY KEY` | Auto-increment |
| `title` | `TEXT NOT NULL` | Book title |
| `author` | `TEXT` | Optional |
| `status` | `TEXT NOT NULL` | One of: `read`, `reading`, `not_started`, `to_buy` |
| `priority_to_buy` | `TEXT` | One of: `low`, `medium`, `high`, or `NULL` |
| `eligible_to_sell` | `INTEGER NOT NULL DEFAULT 0` | Boolean (`0` or `1`) |
| `notes` | `TEXT` | Free-form notes |
| `added_at` | `TEXT NOT NULL` | ISO 8601 timestamp |
| `finished_at` | `TEXT` | ISO 8601 timestamp, set when marked `read` |

### Constraints

- `status` must be one of: `read`, `reading`, `not_started`, `to_buy`
- `priority_to_buy`, when set, must be one of: `low`, `medium`, `high`
- `eligible_to_sell` must be `0` or `1`

### Indexes (initial)

- Index on `status`
- Index on `title` (for search)

## CLI commands

Binary name: `books`

### Core commands

```bash
books add "The Dispossessed" --status to-buy --priority high
books add "Dune" --author "Frank Herbert" --status not_started

books get 42
books list
books list --status reading
books list --status to-buy --priority high
books list --eligible-to-sell

books search "le guin"
books search "dune" --author "herbert"

books update 42 --status read
books update 42 --status to-buy --priority medium --eligible-to-sell
books update 42 --notes "Borrowed from library"

books delete 42
# or
books archive 42
```

### Utility commands

```bash
books stats
books backup
books import spreadsheet.csv
```

### Output formats

- Default: human-readable table
- `--json` for scripting and automation

Example:

```bash
books list --status reading --json
```

## Configuration

Database path should be configurable via:

1. `BOOKS_DB` environment variable
2. Config file at `~/.config/books/config.toml`

Example config:

```toml
database = "/home/user/books.db"
```

Default fallback (if unset): `~/.local/share/books/books.db`

## Project structure (planned)

```text
books-cli/
├── cmd/books/          # CLI entrypoint
├── internal/
│   ├── db/             # SQLite access, migrations, queries
│   ├── models/         # Book model and enums
│   └── output/         # Table and JSON formatters
├── migrations/         # SQL migration files
└── README.md
```

## Practical tips

### 1. Single module for database logic

Keep all SQL in one internal package (`internal/db`). The CLI (and a future MCP server) should call shared functions like `AddBook()`, `ListBooks(filters)`, and `UpdateBook()` instead of duplicating queries.

### 2. Migrations from day one

Use versioned SQL migrations from the start, even if v1 is just one `schema.sql`. A simple `schema_migrations` table avoids painful changes later when new columns are added.

### 3. Backup is trivial

SQLite is a single file. Support:

```bash
books backup
```

And document manual backup:

```bash
cp ~/.local/share/books/books.db ~/.local/share/books/books.db.bak
```

### 4. Import once from the spreadsheet

Provide a one-time import path from the existing spreadsheet (CSV export). After validation, retire the sheet or keep it as a read-only reference.

### 5. Human-friendly output by default

Optimize for quick terminal use:

- Table output for `list` and `search`
- Clear errors for invalid status/priority values
- `--json` only when needed for scripting

### 6. MCP later, not now

Do not build MCP in v1. If conversational management in Cursor becomes useful, add an MCP server that reuses the same `internal/db` package.

## Initial implementation scope

### v0.1

- [ ] SQLite schema and migrations
- [ ] `add`, `get`, `list`, `search`, `update`, `delete`
- [ ] Status and priority validation
- [ ] Table output + `--json`
- [ ] Configurable database path

### v0.2

- [ ] `stats`
- [ ] `backup`
- [ ] CSV import from spreadsheet

### Future

- [ ] MCP server wrapper
- [ ] Optional sync/export to CSV
- [ ] Reading history and change log

## Development

Requirements:

- Go 1.22+

Planned commands:

```bash
go build -o books ./cmd/books
./books list
```

## License

Private personal project.
