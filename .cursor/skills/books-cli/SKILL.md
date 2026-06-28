---
name: books-cli
description: >-
  Manage the personal books reading list via the books CLI (add, list, search,
  update, wishlist). Use when the user mentions books, reading list,
  livros, lista de leitura, biblioteca, wishlist, to buy, TO_BUY, reading status,
  ler, comprar livro, or books-cli.
---

# books-cli Agent

Operate the `books` CLI in the shell — never simulate database changes.

**Do not** explore `cmd/`, `internal/`, or `docs/COMMANDS.md` to learn usage.

| File | Purpose |
| --- | --- |
| [reference.md](reference.md) | Flags, enums, JSON shapes, search/pagination details |
| [examples.md](examples.md) | User phrase → command mapping |

## Setup

1. `books` on PATH (`go install ./cmd/books` from repo root)
2. Else `./books` after `go build -o books ./cmd/books`

Always append `--json`.

## Add a book

**Title:** Use the user's exact title (pt-BR spelling/accents). Do not invent or translate unless asked.

**Author:** If missing, search the web (or use knowledge when unambiguous). Ask before saving if still unclear.

**Description (required on add):** Fetch a 1–3 sentence synopsis from the web. Language **must match the title** (English title → English; pt-BR title → Brazilian Portuguese). Same rule on `update` when refreshing description.

**Defaults (agent overrides CLI):**

| Field | Agent behavior |
| --- | --- |
| `--status` | Always `TO_BUY` unless user specifies otherwise (CLI default is `NOT_STARTED`) |
| `--category` | Required — pick one enum value; see [reference.md](reference.md#category-enum) |
| `--priority` | Only when user asks |
| `--description` | Required — web synopsis in title language |

**Category tie-breaks:** `SOFTWARE` over `PHILOSOPHY`; `FICTION` for novels; `BIOGRAPHY` for life stories; `THEOLOGY` over `POLITICS_CULTURE` for faith lens; `HISTORY` for broad surveys, `FINANCE_BUSINESS` for money/psychology-of-money books. Prefer specific over `OTHER`.

```bash
books add "<title>" --author "<author>" --category <CATEGORY> --description "<synopsis>" --status TO_BUY [--priority] [--notes "..."] --json
```

Search for duplicates before adding well-known titles (bilingual `--term` variants — see [Search](#search-and-list)).

After success, show id, title, author, category, status, priority (Yes/No), description snippet, and `added_at` in a short table.

## Search and list

Pagination is always on (`page=1`, `limit=20`, max `100`). Never assume one page is the full set — check `total` vs `len(books)`.

The DB mixes **pt-BR and English** titles/descriptions. Search is case-insensitive substring match on title or description; it does not translate. Use multiple `--term` flags (OR) for language variants in one query — details in [reference.md](reference.md#search-query).

**When to use bilingual search:** topic/title lookups and duplicate checks. Not needed for status-only lists (wishlist, currently reading).

**Workflow:** Run `search` or `list` with `--json` → read `total`, `page`, `limit`, `books` → if `total > limit`, show current page and `Page X of Y (N total)`; fetch more pages only when needed.

```bash
books list [--status STATUS] [--priority] [--eligible-to-sell] --page 1 --limit 20 --json
books search [--term "<term>" ...] ["<query>"] [--author "<author>"] --page 1 --limit 20 --json
```

Present list/search results as a table: ID, Title, Author, Category, Status, Priority (`Y` or `-` for booleans).

## Other operations

| Intent | Command |
| --- | --- |
| View one book | `books get <id> --json` |
| Update fields | `books update <id> --status READ [--category FICTION] ... --json` |
| Remove from active lists | `books update <id> --status ARCHIVED --json` |
| Show DB path | `books config --json` |

`update` needs at least one flag. `--priority` without value sets true; confirm intent before clearing priority.

Status/category enums and JSON fields: [reference.md](reference.md). Phrase examples: [examples.md](examples.md).

## Errors

Exit `0` = success; `1` = validation, not found, or DB error. Common: invalid status/category, bad page/limit, missing update flags, unknown ID, no search terms. If `books` missing: `go install ./cmd/books` or build `./books`.
