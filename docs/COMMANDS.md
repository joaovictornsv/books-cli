# Command reference

Detailed usage for the `books` CLI.

## Global flags

| Flag | Description |
| --- | --- |
| `--json` | Machine-readable JSON output |
| `--help` | Help for the current command |
| `--version` | Print CLI version |

## Configuration

Database path resolution (first match wins):

1. `BOOKS_DB` environment variable
2. `database` key in `~/.config/books/config.toml`
3. Default: `~/.local/share/books/books.db`

Example config file:

```toml
database = "/home/user/books.db"
```

## Status values

| Status | Meaning |
| --- | --- |
| `NOT_STARTED` | Owned or queued, not yet reading |
| `READING` | Currently reading |
| `READ` | Finished |
| `TO_BUY` | On the wishlist |
| `ARCHIVED` | Logical deletion; hidden from `list` and `search` |

Status values are case-insensitive on input but stored as uppercase.

## `add`

Add a book to the database.

```bash
books add "The Dispossessed" --status TO_BUY --priority
books add "Dune" --author "Frank Herbert" --status NOT_STARTED
```

### Arguments

| Argument | Description |
| --- | --- |
| `title` | Book title (required) |

### Flags

| Flag | Default | Description |
| --- | --- | --- |
| `--author` | _(empty)_ | Author name |
| `--status` | `NOT_STARTED` | One of the status values above |
| `--priority` | `false` | Set `priority_to_buy` to `1` |
| `--eligible-to-sell` | `false` | Set `eligible_to_sell` to `1` |
| `--notes` | _(empty)_ | Free-form notes |
| `--description` | _(empty)_ | Book description (e.g. from the web) |

### JSON output

Returns a single book object.

## `get`

Show one book by ID or by title.

```bash
books get 42
books get 42 --json
books get --title "Dune" --json
books get --title "Dune" --exact --author "herbert" --json
```

Provide either a positional `id` or `--title`, not both.

### Arguments

| Argument | Description |
| --- | --- |
| `id` | Positive integer book ID (optional when using `--title`) |

### Flags

| Flag | Default | Description |
| --- | --- | --- |
| `--title` | _(none)_ | Case-insensitive title substring lookup |
| `--author` | _(none)_ | AND filter on author substring when using `--title` |
| `--exact` | `false` | Case-insensitive exact title match when using `--title` |

When multiple books match `--title`, the command fails with an ambiguous-title error; narrow with `--author`, `--exact`, or use an ID.

### Exit codes

- `0` on success
- `1` if the book is not found, the ID is invalid, or the title is ambiguous

## `list`

List books with optional filters. Archived books are excluded unless you filter explicitly with `--status ARCHIVED`.

```bash
books list
books list --status READING
books list --status TO_BUY --priority
books list --category FICTION
books list --eligible-to-sell
books list --page 2 --limit 20
books list --status READ --sort finished_at --order desc --limit 10
books list --json
books list --json --fields id,title,status
```

### Flags

| Flag | Default | Description |
| --- | --- | --- |
| `--status` | _(none)_ | Filter by status |
| `--category` | _(none)_ | Filter by category |
| `--priority` | `false` | Only books with `priority_to_buy = 1` |
| `--eligible-to-sell` | `false` | Only books with `eligible_to_sell = 1` |
| `--page` | `1` | Page number (1-based); used with `--limit` |
| `--limit` | `20` | Results per page; used with `--page` |
| `--sort` | `id` | Sort field: `id`, `title`, `author`, `status`, `added_at`, `started_at`, `finished_at` |
| `--order` | `asc` | Sort order: `asc` or `desc` |
| `--fields` | _(none)_ | Comma-separated book fields to return (requires `--json`) |

Pagination applies when either `--page` or `--limit` is passed. If only one is set, the other defaults to `page=1` or `limit=20`.

Allowed `--fields` values: `id`, `title`, `author`, `category`, `status`, `priority_to_buy`, `eligible_to_sell`, `sold`, `notes`, `description`, `added_at`, `started_at`, `finished_at`.

### JSON output

Without pagination:

```json
{
  "books": [ ... ],
  "total": 2
}
```

With pagination:

```json
{
  "books": [ ... ],
  "total": 45,
  "page": 2,
  "limit": 20
}
```

`total` is the full match count across all pages.

With `--fields`, each item in `books` contains only the requested keys (in the order you passed them). Envelope fields (`total`, `page`, `limit`) are always included.

## `search`

Search books by title, description, or author substring (case-insensitive). Optionally filter by author substring with `--author` (AND).

Pass a positional `query` and/or repeatable `--term` flags. Multiple terms are OR'd — a book matches if any term appears in title, description, or author.

```bash
books search "le guin"
books search --term hobbit --term "o hobbit"
books search --term herbert
books search "dune" --author "herbert"
books search --term senhor --term lord --author "tolkien" --page 1 --limit 10
books search "dune" --sort title --order asc --json
books search "dune" --author "herbert" --category FICTION --json
books search "dune" --json --fields id,title,author
```

### Arguments

| Argument | Description |
| --- | --- |
| `query` | Optional substring to match against title, description, or author |

### Flags

| Flag | Default | Description |
| --- | --- | --- |
| `--term` | _(none)_ | Search term substring (repeatable; terms are OR'd across title/description/author) |
| `--author` | _(none)_ | AND filter on author substring (case-insensitive) |
| `--category` | _(none)_ | Filter by category |
| `--page` | `1` | Page number (1-based); used with `--limit` |
| `--limit` | `20` | Results per page; used with `--page` |
| `--sort` | `id` | Sort field: `id`, `title`, `author`, `status`, `added_at`, `started_at`, `finished_at` |
| `--order` | `asc` | Sort order: `asc` or `desc` |
| `--fields` | _(none)_ | Comma-separated book fields to return (requires `--json`) |

Archived books are excluded from results.

### JSON output

Same shape as `list` (including optional `page` and `limit` when paginating).

## `update`

Update one or more fields on an existing book. Only flags you pass are changed.

```bash
books update 42 --status READ
books update --ids 1,2,3 --status READ --json
books update 42 --status TO_BUY --priority --eligible-to-sell
books update 42 --notes "Borrowed from library"
books update 42 --description "Epic science fiction saga set on Arrakis."
books update 42 --title "Dune" --author "Frank Herbert"
books update 42 --status ARCHIVED
books update 42 --no-priority
```

To hide a book from `list` and `search`, set `--status ARCHIVED`. Use `list --status ARCHIVED` to view archived books.

Prefer `--no-priority`, `--no-eligible-to-sell`, and `--no-sold` to clear boolean fields explicitly (especially for agents).

### Arguments

| Argument | Description |
| --- | --- |
| `id` | Positive integer book ID (use this **or** `--ids`, not both) |

### Flags

| Flag | Description |
| --- | --- |
| `--ids` | Comma-separated book IDs for bulk update (e.g. `1,2,3`) |
| `--title` | New title |
| `--author` | New author |
| `--status` | New status |
| `--notes` | New notes |
| `--description` | New description |
| `--started-at` | Reading start timestamp (RFC3339); pass `""` to clear |
| `--finished-at` | Reading finish timestamp (RFC3339); pass `""` to clear |
| `--priority` | Set `priority_to_buy` (`true` → `1`, `false` → `0`) |
| `--no-priority` | Clear `priority_to_buy` (set to `0`) |
| `--eligible-to-sell` | Set `eligible_to_sell` |
| `--no-eligible-to-sell` | Clear `eligible_to_sell` |
| `--sold` | Set `sold` |
| `--no-sold` | Clear `sold` |

Status changes do not modify `started_at` or `finished_at`. Set those timestamps explicitly with the flags above.

### Bulk update JSON output

When updating multiple books with `--ids` and `--json`:

```json
{
  "updated": [ { "id": 1, "title": "...", "status": "READ" } ],
  "count": 3
}
```

## `delete`

Permanently remove a book from the database. Unlike setting `--status ARCHIVED`, the row is deleted and cannot be recovered from the CLI.

```bash
books delete 42
books delete 42 -y
books delete 42 --json -y
```

When using `--json`, you must pass `--yes` / `-y` to confirm. In interactive mode on a TTY, the CLI prompts unless `-y` is passed.

### Arguments

| Argument | Description |
| --- | --- |
| `id` | Positive integer book ID |

### Flags

| Flag | Description |
| --- | --- |
| `--yes`, `-y` | Confirm deletion without prompting |

### Exit codes

- `0` on success
- `1` if the book is not found, the ID is invalid, or confirmation was not given

Returns the deleted book (same JSON shape as `get`).

## `check`

Find likely duplicate books by title before `add`. Matches **title only** (not description). Archived books are excluded.

```bash
books check --title "Dune"
books check --title "Dune" --author "Herbert" --json
books check --title "Dune" --exact --json
```

### Flags

| Flag | Default | Description |
| --- | --- | --- |
| `--title` | _(required)_ | Title to check for duplicates |
| `--author` | _(none)_ | Substring to match against author (case-insensitive) |
| `--exact` | `false` | Case-insensitive exact title match (default: substring on title) |

### JSON output

Same shape as `list` (without pagination): `{ "books": [...], "total": N }`.

## `schema`

Show machine-readable enums and book field semantics. No database access.

```bash
books schema --json
```

### JSON output

```json
{
  "statuses": [{ "value": "READ", "description": "Finished" }],
  "categories": [{ "value": "FICTION", "description": "..." }],
  "fields": [{ "name": "id", "type": "integer", "optional": false, "description": "..." }]
}
```

## `config`

Print the effective CLI configuration.

```bash
books config
books config --json
```

Example human output:

```text
database_path: /home/user/.local/share/books/books.db
config_path: /home/user/.config/books/config.toml
config_exists: false
source: default
```

Example JSON output:

```json
{
  "database_path": "/home/user/.local/share/books/books.db",
  "config_path": "/home/user/.config/books/config.toml",
  "config_exists": false,
  "source": "default"
}
```

`source` is one of: `env`, `config_file`, `default`.

## `count`

Count books matching optional filters without paginating through `list`.

```bash
books count
books count --status READ --json
books count --status TO_BUY --category FICTION --priority --json
```

### Flags

| Flag | Default | Description |
| --- | --- | --- |
| `--status` | _(none)_ | Filter by status |
| `--category` | _(none)_ | Filter by category |
| `--priority` | `false` | Only priority-to-buy books |
| `--eligible-to-sell` | `false` | Only eligible-to-sell books |

Archived books are excluded (same as `list`).

### JSON output

```json
{ "total": 42 }
```

## `stats`

Show library aggregates: counts by status and category, books finished in a year, and priority wishlist size.

```bash
books stats
books stats --year 2025 --json
```

### Flags

| Flag | Default | Description |
| --- | --- | --- |
| `--year` | current year | Year used for `finished_this_year` |

### JSON output

```json
{
  "year": 2025,
  "by_status": { "READ": 12, "READING": 2, "TO_BUY": 8, "NOT_STARTED": 5 },
  "by_category": { "FICTION": 10, "SOFTWARE": 7 },
  "finished_this_year": 4,
  "priority_wishlist": 3
}
```

Archived books are excluded from `by_status` and `by_category`. `priority_wishlist` counts `TO_BUY` books with `priority_to_buy = 1`.

## `backup`

Create a consistent copy of the SQLite database at another path using SQLite's `VACUUM INTO` (safe while the database is open).

```bash
books backup --output /path/to/books-backup.db
books backup --output /path/to/books-backup.db --force --json
```

### Flags

| Flag | Default | Description |
| --- | --- | --- |
| `--output` | _(required)_ | Destination file path |
| `--force` | `false` | Overwrite destination if it already exists |

Uses the resolved database path from `books config`. Fails if the source database does not exist.

### JSON output

```json
{
  "source": "/home/user/.local/share/books/books.db",
  "output": "/path/to/books-backup.db"
}
```

## `export`

Dump the library to a structured JSON or CSV file. Complements `backup` (raw database copy) with human-editable data.

```bash
books export --format json --output books.json
books export --format csv --output books.csv
books export --format json --output - --json
```

### Flags

| Flag | Default | Description |
| --- | --- | --- |
| `--format` | _(required)_ | `json` or `csv` |
| `--output` | _(required)_ | Destination file path, or `-` for stdout |
| `--include-archived` | `false` | Include archived books |

### JSON output

When `--output` is a file path (not `-`), confirmation JSON:

```json
{
  "output": "books.json",
  "format": "json",
  "total": 42
}
```

Exported JSON uses the same book shape as `list` / `get`, wrapped as `{ "books": [], "total": N }`.

## `import`

Load books from a JSON or CSV file. Upserts by `id`: updates existing rows, inserts new ones. Use `--dry-run` to validate without writing.

```bash
books import --input books.json
books import --input books.csv --dry-run --json
```

### Flags

| Flag | Default | Description |
| --- | --- | --- |
| `--input` | _(required)_ | Source file path (`.json` or `.csv`) |
| `--dry-run` | `false` | Validate and report counts without modifying the database |

### JSON output

```json
{
  "created": 2,
  "updated": 5,
  "total": 7,
  "dry_run": false
}
```

## Exit codes

- `0` — success
- `1` — validation error, not found, or other failure

Errors are written to stderr.
