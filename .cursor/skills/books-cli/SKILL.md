---
name: books-cli
description: >-
  Manage the personal books reading list via the books CLI (add, list, search,
  update, archive, wishlist). Use when the user mentions books, reading list,
  livros, lista de leitura, biblioteca, wishlist, to buy, TO_BUY, reading status,
  ler, comprar livro, or books-cli.
---

# books-cli Agent

Operate the `books` CLI on behalf of the user. Always run commands in the shell — never simulate database changes.

**Do not** explore `cmd/`, `internal/`, or `docs/COMMANDS.md` to learn how the CLI works.

| File | Purpose |
| --- | --- |
| [reference.md](reference.md) | Flags, status values, JSON shapes, pagination limits |
| [examples.md](examples.md) | User phrase → command mapping |

## Binary setup

Resolve the CLI in this order:

1. **`books`** on PATH (`go install ./cmd/books` from repo root)
2. **`./books`** from repo root if PATH binary is missing: `go build -o books ./cmd/books`

Verify with `books --version` or `./books --version`.

Append `--json` on every command for reliable parsing.

## Add a book

### Required from user

The user **must** provide the **title**. Do not invent or translate it.

- If the user gives the title in **pt-BR**, use that exact title (same spelling, accents, casing).
- Only normalize when the user explicitly asks (e.g. English edition name).

### Author resolution

If the user did not provide an author:

1. Search the web for the book title (+ edition/language if pt-BR).
2. Fall back to your knowledge base when the match is unambiguous.
3. If still ambiguous, ask the user before saving.

### Defaults (agent overrides CLI defaults)

| Field | Default | Override |
| --- | --- | --- |
| `--status` | `TO_BUY` | Only when user specifies another status |
| `--priority` | omit (false) | Pass `--priority` only when user asks for priority |

The CLI default for status is `NOT_STARTED` — **always pass `--status TO_BUY`** unless the user says otherwise.

### Command

```bash
books add "<title>" --author "<author>" --status TO_BUY [--priority] [--notes "..."] --json
```

Search for duplicates before adding well-known titles. See [examples.md](examples.md).

### Validation output

After a successful add, present saved data in this format so the user can confirm:

```markdown
## Book saved

| Field | Value |
| --- | --- |
| **ID** | 42 |
| **Title** | O Pequeno Príncipe |
| **Author** | Antoine de Saint-Exupéry |
| **Status** | TO_BUY |
| **Priority** | No |
| **Added** | 2024-06-27T12:00:00Z |

Use `books get <id>` to view details or `books update <id>` to change fields.
```

Map JSON fields: `priority_to_buy` → Yes/No, `eligible_to_sell` → Yes/No, `sold` → Yes/No.

On failure, show stderr and suggest a fix (duplicate title is not enforced — focus on validation errors).

## List and search (pagination)

Pagination is **always active** (default `page=1`, `limit=20`, max `limit=100`). Handle it explicitly — never assume one page is the full result set.

See [reference.md](reference.md) for flags and [examples.md](examples.md) for phrase → command mapping.

### Workflow

1. Run the query with `--json`.
2. Read `total`, `page`, `limit`, and `books` from the response.
3. If `total == 0`, say no matches.
4. If `total <= limit`, show all results.
5. If `total > limit`:
   - Show the current page in a table.
   - State: `Page {page} of {ceil(total/limit)} ({total} total)`.
   - Fetch remaining pages only when the user needs the full list; otherwise offer the next page.

### Commands

```bash
books list [--status STATUS] [--priority] [--eligible-to-sell] --page 1 --limit 20 --json
books search "<query>" [--author "<author>"] --page 1 --limit 20 --json
```

### List/search table format

```markdown
## Books (page 1 of 3 — 45 total)

| ID | Title | Author | Status | Priority |
| --- | --- | --- | --- | --- |
| 1 | Dune | Frank Herbert | READING | - |
| 2 | O Hobbit | J.R.R. Tolkien | TO_BUY | Y |
```

Boolean columns: `priority_to_buy` / `eligible_to_sell` / `sold` → `Y` or `-`.

### Pagination rules

- Start at `--page 1` unless the user asks for a specific page.
- Use `--limit 20` unless the user requests another size (1–100).
- When iterating all pages: increment `page` until `page * limit >= total` or `books` is empty.
- Never report "these are all your books" when `total > len(books)`.

## Other operations

| Intent | Command |
| --- | --- |
| View one book | `books get <id> --json` |
| Update fields | `books update <id> --status READ [--priority] ... --json` |
| Remove from active lists | `books archive <id> --json` |
| Show DB path | `books config --json` |

**Status values:** `NOT_STARTED`, `READING`, `READ`, `TO_BUY`, `ARCHIVED` (case-insensitive input). Full enum in [reference.md](reference.md).

**Update:** at least one flag required. Setting `--priority` without a value sets priority to true; to clear priority on update, the user must say so — confirm intent if unclear.

## Error handling

- Exit code `0` = success; `1` = validation, not found, or DB error.
- Parse stderr for invalid status, bad page/limit, missing update flags, or unknown ID.
- If `books` is missing, run `go install ./cmd/books` from the repo root; if that fails, try `go build -o books ./cmd/books` and use `./books`.
