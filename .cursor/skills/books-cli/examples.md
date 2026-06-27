# books-cli examples

Phrase → command mapping. Always append `--json`. Resolve binary per [SKILL.md](SKILL.md#binary-setup).

## Add

| User says | Action |
| --- | --- |
| "Add Dune to my wishlist" | Search author if missing, then `books add "Dune" --author "Frank Herbert" --status TO_BUY --json` |
| "Adiciona O Hobbit" | Keep pt-BR title: `books add "O Hobbit" --author "J.R.R. Tolkien" --status TO_BUY --json` |
| "Add 1984, priority" | `books add "1984" --author "George Orwell" --status TO_BUY --priority --json` |
| "Add Sapiens as reading" | User override: `books add "Sapiens" --author "Yuval Noah Harari" --status READING --json` |

Before adding, search for duplicates when the title is well-known:

```bash
books search "dune" --page 1 --limit 20 --json
```

## Search

| User says | Command |
| --- | --- |
| "Do I already have Dune?" | `books search "dune" --page 1 --limit 20 --json` |
| "Any book with 'senhor' in the title?" | `books search "senhor" --page 1 --limit 20 --json` |
| "Find books by Le Guin" | `books search "guin" --page 1 --limit 20 --json` (empty query is rejected) |
| "Search 1984 page 2" | `books search "1984" --page 2 --limit 20 --json` |
| "Dune by Herbert" | `books search "dune" --author "herbert" --page 1 --limit 20 --json` |

## List

| User says | Command |
| --- | --- |
| "What am I reading?" | `books list --status READING --page 1 --limit 20 --json` |
| "Show my wishlist" | `books list --status TO_BUY --page 1 --limit 20 --json` |
| "Priority books to buy" | `books list --status TO_BUY --priority --page 1 --limit 20 --json` |
| "Books I haven't started" | `books list --status NOT_STARTED --page 1 --limit 20 --json` |
| "All my books" | `books list --page 1 --limit 20 --json` — paginate if `total > limit` |

## Update & other

| User says | Command |
| --- | --- |
| "Mark book 42 as read" | `books update 42 --status READ --json` |
| "Start reading book 7" | `books update 7 --status READING --json` |
| "Show book 42" | `books get 42 --json` |
| "Remove book 42 from my list" | `books archive 42 --json` |
| "Where is my database?" | `books config --json` |

## Pagination response

When `total > limit`, always report:

```text
Page 1 of 3 (45 total)
```

Fetch page 2 only when the user asks or needs the full set.
