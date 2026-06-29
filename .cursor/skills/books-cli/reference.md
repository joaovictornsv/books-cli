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

### `list`

| Flag | Default |
| --- | --- |
| `--status` | none |
| `--priority` | false (filter) |
| `--eligible-to-sell` | false (filter) |
| `--page` | 1 |
| `--limit` | 20 (max 100) |

Archived excluded unless `--status ARCHIVED`.

### `search [query]`

| Flag | Default | Notes |
| --- | --- | --- |
| `query` (positional) | none | Substring on title or description |
| `--term` | none | Repeatable; all terms OR'd |
| `--author` | none | AND filter on author substring |
| `--page` / `--limit` | 1 / 20 | Max limit 100 |

**At least one** positional `query` or `--term` required. Each term matches title **or** description (case-insensitive).

```bash
books search "dune"
books search --term hobbit --term "o hobbit"
books search "senhor" --term lord --author "tolkien"
```

For pt-BR/English variants, prefer one command with multiple `--term` flags over separate searches.

### `update [id]`

At least one flag: `--title`, `--author`, `--category`, `--status`, `--notes`, `--description`, `--started-at`, `--finished-at`, `--priority`, `--eligible-to-sell`, `--sold`. Pass `--category ""`, `--started-at ""`, or `--finished-at ""` to clear.

Status changes do not set or clear timestamps automatically. Set `--started-at` / `--finished-at` explicitly (RFC3339).

### `get [id]` · `config`

No extra flags beyond global `--json`.

## JSON shapes

**Single book** (`add`, `get`, `update`):

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

**List/search:** `{ "books": [], "total": 45, "page": 1, "limit": 20 }` — `total` is the full filtered count.
