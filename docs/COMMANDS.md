# Command reference

Detailed usage for the `books` CLI.

> **For AI agents:** use [`.cursor/skills/books-cli/SKILL.md`](../.cursor/skills/books-cli/SKILL.md) — do not use this file for agent workflows.

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

## Category values

Nullable on existing books. Run `books schema --json` for the canonical list with descriptions.

| Category | Meaning |
| --- | --- |
| `THEOLOGY` | Christian faith, Bible, devotionals, apologetics, pastoral |
| `FICTION` | Novels, short stories, literary and genre fiction |
| `SOFTWARE` | Programming, software engineering, CS |
| `PHILOSOPHY` | Philosophy, ethics, stoicism, political philosophy |
| `HISTORY` | Historical narrative and historiography |
| `PERSONAL_DEVELOPMENT` | Self-help, productivity, habits, popular psychology |
| `FINANCE_BUSINESS` | Money, investing, economics, business |
| `SCIENCE` | Natural sciences, math popularization |
| `POLITICS_CULTURE` | Political/social commentary, cultural criticism |
| `BIOGRAPHY` | Biographies, memoirs, autobiographies |
| `OTHER` | Catch-all when no other category fits |

## Pagination and sorting

Used by `list` and `search`:

| Flag | Default | Description |
| --- | --- | --- |
| `--page` | `1` | Page number (1-based) |
| `--limit` | `20` | Results per page (max 100) |
| `--sort` | `id` | `id`, `title`, `author`, `status`, `added_at`, `started_at`, `finished_at` |
| `--order` | `asc` | `asc` or `desc` |
| `--fields` | _(none)_ | Comma-separated book fields (requires `--json`) |

Allowed `--fields`: `id`, `title`, `author`, `category`, `status`, `priority_to_buy`, `eligible_to_donate`, `donated`, `notes`, `description`, `added_at`, `started_at`, `finished_at`.

With `--fields`, each item in `books` contains only the requested keys. Envelope fields (`total`, `page`, `limit`) are always included.

## JSON shapes

**Single book** (`add`, `get`, `update`, `delete`): returns one book object with `id`, `title`, `author`, `category`, `status`, boolean flags, `notes`, `description`, and timestamps (`added_at`, `started_at`, `finished_at`).

**List envelope** (`list`, `search`, `check`):

```json
{
  "books": [ ... ],
  "total": 45,
  "page": 2,
  "limit": 20
}
```

`total` is the full match count across all pages. `check` omits `page` and `limit`.

**Bulk update** (`update --ids`): `{ "updated": [ ... ], "count": N }`.

**Count:** `{ "total": N }`.

**Stats:** `{ "year": 2025, "by_status": {}, "by_category": {}, "finished_this_year": N, "priority_wishlist": N }`.

**Config:** `{ "database_path": "...", "config_path": "...", "config_exists": false, "source": "default" }` — `source` is `env`, `config_file`, or `default`.

**Export confirmation:** `{ "output": "...", "format": "json", "total": N }`.

**Import:** `{ "created": N, "updated": M, "total": T, "dry_run": false }`.

**Backup:** `{ "source": "...", "output": "..." }`.

**Schema:** `{ "statuses": [...], "categories": [...], "fields": [...] }`.

## `add`

```bash
books add "Dune" --author "Frank Herbert" --category FICTION --status TO_BUY --priority
books add "The Dispossessed" --status NOT_STARTED --description "Anarchist utopia novel."
```

| Flag | Default | Description |
| --- | --- | --- |
| `--author` | _(empty)_ | Author name |
| `--category` | _(empty)_ | Category enum value (see above or `books schema --json`) |
| `--status` | `NOT_STARTED` | One of the status values |
| `--priority` | `false` | Set `priority_to_buy` to `1` |
| `--eligible-to-donate` | `false` | Set `eligible_to_donate` to `1` |
| `--notes` | _(empty)_ | Free-form notes |
| `--description` | _(empty)_ | Book description |

## `get`

Provide either positional `id` or `--title`, not both.

```bash
books get 42 --json
books get --title "Dune" --exact --author "herbert" --json
```

| Flag | Default | Description |
| --- | --- | --- |
| `--title` | _(none)_ | Case-insensitive title substring lookup |
| `--author` | _(none)_ | AND filter on author substring when using `--title` |
| `--exact` | `false` | Case-insensitive exact title match when using `--title` |

Fails with ambiguous-title error when multiple books match; narrow with `--author`, `--exact`, or use an ID.

## `list`

Archived books excluded unless `--status ARCHIVED`.

```bash
books list --status READING
books list --status TO_BUY --category FICTION --priority
books list --status READ --sort finished_at --order desc --limit 10 --json
books list --json --fields id,title,status
```

| Flag | Default | Description |
| --- | --- | --- |
| `--status` | _(none)_ | Filter by status |
| `--category` | _(none)_ | Filter by category |
| `--priority` | `false` | Only books with `priority_to_buy = 1` |
| `--eligible-to-donate` | `false` | Only books with `eligible_to_donate = 1` |

Plus [pagination and sorting](#pagination-and-sorting) flags.

## `search`

Search by title, description, or author substring (case-insensitive). Archived books excluded.

Pass a positional `query` and/or repeatable `--term` flags. Multiple terms are OR'd.

```bash
books search "le guin"
books search --term hobbit --term "o hobbit"
books search "dune" --author "herbert" --category FICTION --json
books search "dune" --json --fields id,title,author
```

| Flag | Default | Description |
| --- | --- | --- |
| `--term` | _(none)_ | Search term (repeatable; terms OR'd across title/description/author) |
| `--author` | _(none)_ | AND filter on author substring |
| `--category` | _(none)_ | Filter by category |

Plus [pagination and sorting](#pagination-and-sorting) flags. At least one positional `query` or `--term` required.

## `check`

Pre-add duplicate detection by title (not description). Archived books excluded.

```bash
books check --title "Dune" --author "Herbert" --json
books check --title "Dune" --exact --json
```

| Flag | Default | Description |
| --- | --- | --- |
| `--title` | _(required)_ | Title to check |
| `--author` | _(none)_ | Author substring filter |
| `--exact` | `false` | Exact title match (default: substring) |

## `update`

Only flags you pass are changed. To hide from `list`/`search`, set `--status ARCHIVED`.

```bash
books update 42 --status READ --finished-at "2025-06-01T00:00:00Z"
books update --ids 1,2,3 --status READ --json
books update 42 --no-priority
books update 42 --status ARCHIVED
```

Provide either positional `id` or `--ids` (comma-separated), not both.

| Flag | Description |
| --- | --- |
| `--ids` | Comma-separated IDs for bulk update |
| `--title`, `--author`, `--category` | New values; pass `--category ""` to clear |
| `--status` | New status |
| `--notes`, `--description` | New text |
| `--started-at`, `--finished-at` | RFC3339 timestamps; pass `""` to clear |
| `--priority` / `--no-priority` | Set or clear `priority_to_buy` |
| `--eligible-to-donate` / `--no-eligible-to-donate` | Set or clear `eligible_to_donate` |
| `--donated` / `--no-donated` | Set or clear `donated` |

Status changes do not modify timestamps — set `--started-at` / `--finished-at` explicitly.

## `delete`

Permanently removes the row (unlike `--status ARCHIVED`). **Requires `-y` with `--json`.** Interactive mode prompts on a TTY unless `-y` is passed.

```bash
books delete 42 -y --json
```

| Flag | Description |
| --- | --- |
| `--yes`, `-y` | Confirm without prompting |

## `count`

Count books matching filters without paginating. Archived excluded (same as `list`).

```bash
books count --status READ --json
books count --status TO_BUY --category FICTION --priority --json
```

| Flag | Default | Description |
| --- | --- | --- |
| `--status` | _(none)_ | Filter by status |
| `--category` | _(none)_ | Filter by category |
| `--priority` | `false` | Only priority-to-buy books |
| `--eligible-to-donate` | `false` | Only eligible-to-donate books |

## `stats`

Library aggregates. Archived excluded from `by_status` and `by_category`.

```bash
books stats --year 2025 --json
```

| Flag | Default | Description |
| --- | --- | --- |
| `--year` | current year | Year for `finished_this_year` |

`priority_wishlist` counts `TO_BUY` books with `priority_to_buy = 1`.

## `backup`

Consistent SQLite copy via `VACUUM INTO`. Uses resolved database path from `books config`.

```bash
books backup --output /path/to/books-backup.db --force --json
```

| Flag | Default | Description |
| --- | --- | --- |
| `--output` | _(required)_ | Destination file path |
| `--force` | `false` | Overwrite if destination exists |

## `export`

Dump library to JSON or CSV (human-editable; complements `backup`).

```bash
books export --format json --output books.json
books export --format csv --output - --json
```

| Flag | Default | Description |
| --- | --- | --- |
| `--format` | _(required)_ | `json` or `csv` |
| `--output` | _(required)_ | File path or `-` for stdout |
| `--include-archived` | `false` | Include archived books |

## `import`

Load from JSON or CSV. Upserts by `id`. Use `--dry-run` to validate without writing.

```bash
books import --input books.json --dry-run --json
```

| Flag | Default | Description |
| --- | --- | --- |
| `--input` | _(required)_ | Source file (`.json` or `.csv`) |
| `--dry-run` | `false` | Validate without modifying the database |

## `schema`

Machine-readable enums and field semantics. No database access.

```bash
books schema --json
```

## `config`

Print effective configuration.

```bash
books config --json
```

## `version`

Show CLI version and build metadata.

```bash
books version
books version --json
```

Example JSON output:

```json
{
  "version": "0.5.0",
  "commit": "unknown",
  "go_version": "go1.25.0"
}
```

Release binaries embed the git commit via build flags; local `go build` / `go install` builds show `commit: "unknown"` unless you pass matching `-ldflags`.

## Exit codes

- `0` — success
- `1` — validation error, not found, or other failure

Errors are written to stderr.
