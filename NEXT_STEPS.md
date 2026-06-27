# Next steps

Roadmap grouped by priority. **Each priority block is intended to ship as one release** — implement everything in a group, then cut a version before starting the next.

Emphasis on agent-friendly workflows: `--json`, minimal round-trips, predictable enums.

---

## Priority 1 — insights & backup

Aggregates and database safety. Small, mostly read-only changes; `count` can reuse existing DB helpers.

| Item | Type | Summary |
| --- | --- | --- |
| `count` | command | Totals with optional filters — no paginating through `list` |
| `stats` | command | Aggregates by status/category, recent activity, wishlist size |
| `backup` | command | Copy the SQLite database file |
| `config --json` | flag | JSON output on `config` for agent consistency |

### `count`

```bash
books count [--status STATUS] [--category CATEGORY] [--priority] [--eligible-to-sell] [--sold] --json
```

```json
{ "total": 42 }
```

### `stats`

```bash
books stats [--year YYYY] --json
```

```json
{
  "by_status": { "READ": 12, "READING": 2, "TO_BUY": 8, "NOT_STARTED": 5 },
  "by_category": { "FICTION": 10, "SOFTWARE": 7 },
  "finished_this_year": 4,
  "priority_wishlist": 3
}
```

### `backup`

```bash
books backup --output /path/to/books-backup.db [--force]
```

Copies the resolved database path from `books config`. Fails if destination exists unless `--force`.

---

## Priority 2 — agent safety & core filters

Commands and flags that make everyday agent workflows safer and more precise.

| Item | Type | Summary |
| --- | --- | --- |
| `schema` | command | Machine-readable enums and field semantics |
| `restore` | command | Undo `archive` with an explicit target status |
| `check` | command | Pre-add duplicate detection by title/author |
| `--category` | flag on `list`, `search` | Filter by category (e.g. fiction wishlist) |
| `--sold` | flag on `list` | Filter sold books |
| Clear booleans | flags on `update` | `--no-priority`, `--no-eligible-to-sell`, `--no-sold` |

### `schema`

```bash
books schema --json
```

Returns status/category values, field types, and brief descriptions.

### `restore`

```bash
books restore <id> [--status STATUS] --json
```

Default `--status`: `NOT_STARTED`. Only valid when current status is `ARCHIVED`.

### `check`

```bash
books check --title "Dune" [--author "Herbert"] [--exact] --json
```

Returns likely duplicates before `add`.

### `--category` on list

```bash
books list --status TO_BUY --category FICTION --json
```

### Clear booleans on update

```bash
books update 42 --no-priority --json
```

### Agent skill updates (ship with Priority 2)

| Item | Action |
| --- | --- |
| Document `--fields` | `list` / `search --fields id,title,status` — smaller JSON payloads |
| Document `delete` | Destructive; confirm with user before running |
| Duplicate check | Point to `check` once shipped; until then, bilingual `search --term` recipe |

---

## Priority 3 — search & list ergonomics

Better discovery and navigation without schema changes.

| Item | Type | Summary |
| --- | --- | --- |
| `--sort` / `--order` | flags on `list`, `search` | Sort by `added_at`, `title`, `finished_at`, etc. |
| Search scope | `search` | Match `--term` against `notes` and `author` (keep `--author` as AND filter) |
| Resolve by title | `get` or `resolve` | Look up a book by title when the user does not give an ID |
| `recent` | command | Shortcut for recently finished or added books |

### `--sort`

```bash
books list --status READ --sort finished_at --order desc --limit 10 --json
```

Default today: `id ASC`.

### Resolve by title

```bash
books get --title "Dune" --json
# or
books resolve "dune" --json
```

### `recent`

```bash
books recent [--field finished_at|added_at] --limit 5 --json
```

---

## Priority 4 — bulk & portability

Multi-book operations and data migration beyond raw DB copy.

| Item | Type | Summary |
| --- | --- | --- |
| Bulk update | `update` | Update several books in one call |
| `export` | command | Dump library to JSON or CSV |
| `import` | command | Load from JSON or CSV (`--dry-run` supported) |

### Bulk update

```bash
books update --ids 1,2,3 --status READ --json
```

```json
{ "updated": [], "count": 3 }
```

### `export` / `import`

```bash
books export --format json --output books.json
books import --input books.json [--dry-run]
```

Complements `backup` (file copy) with structured, human-editable data.

---

## Priority 5 — schema & polish

Larger or optional work; may require DB migrations.

| Item | Type | Summary |
| --- | --- | --- |
| `completion` docs | docs | Shell completion generation for human users |
| ISBN / external ID | schema + CLI | Richer duplicate checks and web lookups |
| Reading progress | schema + CLI | Optional `current_page` / `total_pages` for `READING` books |

---

## Release checklist

When finishing a priority group:

1. Implement all items in that group (including skill/doc updates listed there).
2. Update [CHANGELOG.md](CHANGELOG.md) and [docs/COMMANDS.md](docs/COMMANDS.md).
3. Tag a release.
4. Remove or move completed items from this file into the changelog entry.
