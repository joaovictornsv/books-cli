# books-cli reference

**Binary:** `books` on PATH · Install: `go install ./cmd/books` · Dev fallback: `go build -o books ./cmd/books`

## Global flags

| Flag | Description |
| --- | --- |
| `--json` | JSON output (use for all agent operations) |
| `--help` | Command help |
| `--version` | CLI version |

## Configuration

Database path (first match):

1. `BOOKS_DB` env var
2. `database` in `~/.config/books/config.toml`
3. Default: `~/.local/share/books/books.db`

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
| `--category` | empty | **Agent must set on add** |
| `--status` | `NOT_STARTED` | Agent default: `TO_BUY` |
| `--priority` | false | Sets `priority_to_buy = 1` |
| `--eligible-to-sell` | false | |
| `--notes` | empty | |
| `--description` | empty | **Agent must set on add** — web synopsis; same language as title |

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

Title or description substring (case-insensitive). Optional `--author` substring filter. Same pagination flags as `list`.

| Flag | Default | Notes |
| --- | --- | --- |
| `query` (positional) | none | Optional; substring match on title or description |
| `--term` | none | Repeatable; each value is OR'd with other terms |
| `--author` | none | AND filter on author substring |
| `--page` | 1 | |
| `--limit` | 20 (max 100) | |

**At least one** positional `query` or `--term` is required. Positional query and `--term` can be combined — all terms are OR'd. Within each term, a book matches if the substring appears in **title or description**.

```bash
books search "dune"
books search --term hobbit --term "o hobbit"
books search "senhor" --term lord --author "tolkien"
```

**Bilingual collection:** titles and descriptions may be in **pt-BR or English**. Search is literal substring match — it does not translate. For alternate-language variants, use one command with multiple `--term` flags (e.g. `--term hobbit --term "o hobbit"`) instead of separate searches. See [SKILL.md](SKILL.md#bilingual-search-pt-br-and-english).

### `update [id]`

At least one flag: `--title`, `--author`, `--category`, `--status`, `--notes`, `--description`, `--priority`, `--eligible-to-sell`, `--sold`.

Pass `--category ""` to clear an existing category.

### `get [id]` · `archive [id]` · `config`

No extra flags beyond global `--json`.

## JSON shapes

**Single book** (`add`, `get`, `update`, `archive`):

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

**List/search:**

```json
{
  "books": [],
  "total": 45,
  "page": 1,
  "limit": 20
}
```

`total` is the full filtered count across all pages.

## Status side-effects (update)

- → `READING`: sets `started_at` if unset
- → `READ`: sets `finished_at` if unset
- Leaves `READING`/`READ`: may clear respective timestamps
