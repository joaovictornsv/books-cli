# books-cli reference

**Binary:** `books` on PATH · `go install ./cmd/books` · Dev: `go build -o books ./cmd/books`

## Global flags

| Flag | Description |
| --- | --- |
| `--json` | JSON output (use for all agent operations) |
| `--help` | Command help |
| `--version` | CLI version |

## Configuration

DB path (first match): `BOOKS_DB` env → `database` in `~/.config/books/config.toml` → `~/.local/share/books/books.db`

## Status enum

| Value | Meaning |
| --- | --- |
| `NOT_STARTED` | Owned, not yet reading |
| `READING` | Currently reading |
| `READ` | Finished |
| `TO_BUY` | Wishlist |
| `ARCHIVED` | Hidden from list/search |

Use `books schema --json` for the canonical list with descriptions.

## Category enum

Nullable on existing books; agent must set on every `add`.

| Value | Meaning |
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

## Commands

### `schema`

```bash
books schema --json
```

Returns `{ "statuses": [...], "categories": [...], "fields": [...] }` with value/type descriptions. No database access.

### `add [title]`

| Flag | CLI default | Notes |
| --- | --- | --- |
| `--author` | empty | Optional |
| `--category` | empty | **Agent must set** |
| `--status` | `NOT_STARTED` | Agent default: `TO_BUY` |
| `--priority` | false | Sets `priority_to_buy = 1` |
| `--eligible-to-sell` | false | |
| `--notes` | empty | |
| `--description` | empty | **Agent must set** — same language as title |

### `check`

Pre-add duplicate detection. Matches **title only** (not description).

| Flag | Notes |
| --- | --- |
| `--title` | Required |
| `--author` | Optional AND filter (author substring) |
| `--exact` | Case-insensitive exact title match (default: title substring) |

JSON: same list shape `{ "books": [], "total": N }`.

### `list`

| Flag | Default |
| --- | --- |
| `--status` | none |
| `--category` | none |
| `--priority` | false (filter) |
| `--eligible-to-sell` | false (filter) |
| `--page` | 1 |
| `--limit` | 20 (max 100) |
| `--sort` | `id` |
| `--order` | `asc` |
| `--fields` | none (requires `--json`) |

Archived excluded unless `--status ARCHIVED`.

Allowed `--fields`: `id`, `title`, `author`, `category`, `status`, `priority_to_buy`, `eligible_to_sell`, `sold`, `notes`, `description`, `added_at`, `started_at`, `finished_at`.

### `search [query]`

| Flag | Default | Notes |
| --- | --- | --- |
| `query` (positional) | none | Substring on title, description, or author |
| `--term` | none | Repeatable; all terms OR'd |
| `--author` | none | AND filter on author substring |
| `--category` | none | Filter by category |
| `--page` / `--limit` | 1 / 20 | Max limit 100 |
| `--sort` | `id` | `id`, `title`, `author`, `status`, `added_at`, `started_at`, `finished_at` |
| `--order` | `asc` | `asc` or `desc` |
| `--fields` | none | Requires `--json` |

**At least one** positional `query` or `--term` required. Each term matches title, description, **or** author (case-insensitive).

```bash
books search "dune"
books search --term hobbit --term "o hobbit"
books search "senhor" --term lord --author "tolkien"
```

For pt-BR/English variants, prefer one command with multiple `--term` flags over separate searches.

### `update [id]`

At least one flag: `--title`, `--author`, `--category`, `--status`, `--notes`, `--description`, `--started-at`, `--finished-at`, `--priority`, `--no-priority`, `--eligible-to-sell`, `--no-eligible-to-sell`, `--sold`, `--no-sold`. Pass `--category ""`, `--started-at ""`, or `--finished-at ""` to clear.

Prefer `--no-priority`, `--no-eligible-to-sell`, `--no-sold` to clear booleans explicitly.

Status changes do not set or clear timestamps automatically. Set `--started-at` / `--finished-at` explicitly (RFC3339).

### `delete [id]`

Destructive — permanently removes the row. **Requires `-y` with `--json`.**

```bash
books delete 42 -y --json
```

### `get [id]`

Provide either positional `id` or `--title`, not both.

| Flag | Notes |
| --- | --- |
| `--title` | Case-insensitive title substring lookup |
| `--author` | AND filter when using `--title` |
| `--exact` | Exact title match when using `--title` |

Fails with ambiguous-title error when multiple books match; narrow with `--author`, `--exact`, or use an ID.

### `config`

No extra flags beyond global `--json`.

## JSON shapes

**Single book** (`add`, `get`, `update`, `delete`):

```json
{
  "id": 1,
  "title": "Dune",
  "author": "Frank Herbert",
  "category": "FICTION",
  "status": "TO_BUY",
  "priority_to_buy": 0,
  "eligible_to_sell": 0,
  "sold": 0,
  "notes": "",
  "description": "Epic science fiction saga set on Arrakis.",
  "added_at": "2024-01-01T00:00:00Z",
  "started_at": null,
  "finished_at": null
}
```

**List/search/check:** `{ "books": [], "total": 45, "page": 1, "limit": 20 }` — `total` is the full filtered count.
