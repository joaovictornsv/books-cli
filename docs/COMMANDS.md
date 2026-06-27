# Command reference

Detailed usage for the `books` CLI.

## Global flags

| Flag | Description |
| --- | --- |
| `--json` | Machine-readable JSON output (not supported by `config`) |
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

Show one book by ID.

```bash
books get 42
books get 42 --json
```

### Arguments

| Argument | Description |
| --- | --- |
| `id` | Positive integer book ID |

### Exit codes

- `0` on success
- `1` if the book is not found or the ID is invalid

## `list`

List books with optional filters. Archived books are excluded unless you filter explicitly with `--status ARCHIVED`.

```bash
books list
books list --status READING
books list --status TO_BUY --priority
books list --eligible-to-sell
books list --page 2 --limit 20
books list --json
```

### Flags

| Flag | Default | Description |
| --- | --- | --- |
| `--status` | _(none)_ | Filter by status |
| `--priority` | `false` | Only books with `priority_to_buy = 1` |
| `--eligible-to-sell` | `false` | Only books with `eligible_to_sell = 1` |
| `--page` | `1` | Page number (1-based); used with `--limit` |
| `--limit` | `20` | Results per page; used with `--page` |

Pagination applies when either `--page` or `--limit` is passed. If only one is set, the other defaults to `page=1` or `limit=20`.

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

## `search`

Search books by title or description substring (case-insensitive). Optionally filter by author substring.

```bash
books search "le guin"
books search "dune" --author "herbert"
books search "dune" --author "herbert" --page 1 --limit 10
books search "dune" --author "herbert" --json
```

### Arguments

| Argument | Description |
| --- | --- |
| `query` | Substring to match against title or description |

### Flags

| Flag | Default | Description |
| --- | --- | --- |
| `--author` | _(none)_ | Substring to match against author (case-insensitive) |
| `--page` | `1` | Page number (1-based); used with `--limit` |
| `--limit` | `20` | Results per page; used with `--page` |

Archived books are excluded from results.

### JSON output

Same shape as `list` (including optional `page` and `limit` when paginating).

## `update`

Update one or more fields on an existing book. Only flags you pass are changed.

```bash
books update 42 --status READ
books update 42 --status TO_BUY --priority --eligible-to-sell
books update 42 --notes "Borrowed from library"
books update 42 --description "Epic science fiction saga set on Arrakis."
books update 42 --title "Dune" --author "Frank Herbert"
```

### Arguments

| Argument | Description |
| --- | --- |
| `id` | Positive integer book ID |

### Flags

| Flag | Description |
| --- | --- |
| `--title` | New title |
| `--author` | New author |
| `--status` | New status |
| `--notes` | New notes |
| `--description` | New description |
| `--priority` | Set `priority_to_buy` (`true` → `1`, `false` → `0`) |
| `--eligible-to-sell` | Set `eligible_to_sell` |
| `--sold` | Set `sold` |

### Status side-effects

- Changing status to `READING` sets `started_at` if it is not already set.
- Changing status to `READ` sets `finished_at` if it is not already set.
- Existing timestamps are not overwritten when re-entering `READING` or `READ`.

## `archive`

Logically delete a book by setting its status to `ARCHIVED`.

```bash
books archive 42
```

Archived books are hidden from `list` and `search` by default.

## `delete`

Permanently remove a book from the database. Unlike `archive`, the row is deleted and cannot be recovered from the CLI.

```bash
books delete 42
```

### Arguments

| Argument | Description |
| --- | --- |
| `id` | Positive integer book ID |

### Exit codes

- `0` on success
- `1` if the book is not found or the ID is invalid

Returns the deleted book (same JSON shape as `get`).

## `config`

Print the effective CLI configuration. Always human-readable (ignores `--json`).

```bash
books config
```

Example output:

```text
database_path: /home/user/.local/share/books/books.db
config_path: /home/user/.config/books/config.toml
config_exists: false
source: default
```

`source` is one of: `env`, `config_file`, `default`.

## Planned for v0.2

| Command | Summary |
| --- | --- |
| `stats` | Reading statistics |
| `backup` | Copy the database file |

## Exit codes

- `0` — success
- `1` — validation error, not found, or other failure

Errors are written to stderr.
