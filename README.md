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
| `status` | `TEXT NOT NULL` | One of: `READ`, `READING`, `NOT_STARTED`, `TO_BUY`, `ARCHIVED` |
| `priority_to_buy` | `INTEGER NOT NULL DEFAULT 0` | Boolean (`0` or `1`) |
| `eligible_to_sell` | `INTEGER NOT NULL DEFAULT 0` | Boolean (`0` or `1`) |
| `sold` | `INTEGER NOT NULL DEFAULT 0` | Boolean (`0` or `1`); whether the book was sold |
| `notes` | `TEXT` | Free-form notes |
| `description` | `TEXT` | Optional book description (searchable) |
| `added_at` | `TEXT NOT NULL` | ISO 8601 timestamp |
| `started_at` | `TEXT` | ISO 8601 timestamp; set when reading begins |
| `finished_at` | `TEXT` | ISO 8601 timestamp; set when marked `READ` |

### Status values

All status values are **uppercase**:

| Status | Meaning |
| --- | --- |
| `NOT_STARTED` | Owned or queued, not yet reading |
| `READING` | Currently reading |
| `READ` | Finished |
| `TO_BUY` | On the wishlist |
| `ARCHIVED` | Logical deletion; hidden from normal lists but retained in the database |

`ARCHIVED` replaces hard delete for day-to-day use. A physical `delete` command may still exist for exceptional cases, but `archive` is the preferred path.

### Constraints

- `status` must be one of: `READ`, `READING`, `NOT_STARTED`, `TO_BUY`, `ARCHIVED`
- `priority_to_buy`, `eligible_to_sell`, and `sold` must be `0` or `1`
- `started_at` should be set when status changes to `READING` (and cleared or preserved per update rules)
- `finished_at` should be set when status changes to `READ`

### Indexes (initial)

- Index on `status`
- Index on `title` (for search)

## CLI commands

Binary name: `books`

Full command reference: [docs/COMMANDS.md](docs/COMMANDS.md)

### Core commands

```bash
books add "The Dispossessed" --status TO_BUY --priority
books add "Dune" --author "Frank Herbert" --status NOT_STARTED

books get 42
books list
books list --status READING
books list --status TO_BUY --priority
books list --eligible-to-sell
books list --page 2 --limit 20

books search "le guin"
books search --term hobbit --term "o hobbit"
books search "dune" --author "herbert" --page 1 --limit 10

books update 42 --status READ
books update 42 --status TO_BUY --priority --eligible-to-sell
books update 42 --notes "Borrowed from library"

books archive 42
books delete 42
```

`list` and `search` support pagination (`--page`, `--limit`) for large libraries.

### Utility commands

```bash
books stats
books backup
books config
```

`config` prints the effective CLI configuration (e.g. resolved database path, config file location, and which source won: env, config file, or default).

### Output formats

- Default: human-readable table
- `--json` for scripting and automation

Example:

```bash
books list --status READING --json
```

### Out of scope

- **No `import` command.** Initial data migration from the spreadsheet will be done manually.

## Configuration

Database path should be configurable via:

1. `BOOKS_DB` environment variable
2. Config file at `~/.config/books/config.toml`

Example config:

```toml
database = "/home/user/books.db"
```

Default fallback (if unset): `~/.local/share/books/books.db`

Inspect resolved settings at any time:

```bash
books config
```

## Releases

- [Semantic Versioning](https://semver.org/) with GitHub releases and git tags (`v0.1.0`, …).
- Each release includes a `books-linux-amd64` binary asset (linux/amd64).
- Changes are tracked in [CHANGELOG.md](CHANGELOG.md).
- Release workflow for agents: [.agents/github-releases/SKILL.md](.agents/github-releases/SKILL.md).

## CI

GitHub Actions runs on **every push and every pull request**:

- Unit tests (`go test ./...`)
- Build verification (`go build ./...`)

CI must pass before merging.

## Project structure (planned)

```text
books-cli/
├── .agents/            # Agent skills (e.g. release workflow)
├── cmd/books/          # CLI entrypoint
├── docs/
│   └── COMMANDS.md     # Detailed command usage reference
├── internal/
│   ├── db/             # SQLite access, migrations, queries
│   ├── models/         # Book model and enums
│   └── output/         # Table and JSON formatters
├── migrations/         # SQL migration files
├── CHANGELOG.md
└── README.md
```

## Initial implementation scope

### v0.1

- [ ] SQLite schema and migrations (uppercase statuses, booleans, date fields)
- [ ] `add`, `get`, `list`, `search`, `update`, `archive`, `delete`
- [ ] Status and boolean field validation
- [ ] Table output + `--json`
- [ ] Configurable database path + `config` command
- [ ] `docs/COMMANDS.md`
- [ ] CI (test + build on push/PR)
- [ ] `CHANGELOG.md` maintained per release
- [ ] Pagination for `list` and `search`

### v0.2

- [ ] `stats`
- [ ] `backup`

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
