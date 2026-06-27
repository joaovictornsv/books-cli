# books-cli examples

Phrase → command mapping. Always append `--json`. Resolve binary per [SKILL.md](SKILL.md#binary-setup).

## Add

| User says | Action |
| --- | --- |
| "Add Dune to my wishlist" | Resolve author + English description from the web, classify as `FICTION`, then `books add "Dune" --author "Frank Herbert" --category FICTION --description "Epic science fiction saga set on the desert planet Arrakis." --status TO_BUY --json` |
| "Adiciona O Hobbit" | Keep pt-BR title; fetch pt-BR description: `books add "O Hobbit" --author "J.R.R. Tolkien" --category FICTION --description "Bilbo Bolseiro parte em uma aventura inesperada com treze anões e um mago." --status TO_BUY --json` |
| "Add 1984, priority" | `books add "1984" --author "George Orwell" --category FICTION --description "In a totalitarian future, Winston Smith rebels against omnipresent surveillance." --status TO_BUY --priority --json` |
| "Add Sapiens as reading" | User override status: `books add "Sapiens" --author "Yuval Noah Harari" --category HISTORY --description "A brief history of humankind from the Cognitive Revolution to the present." --status READING --json` |
| "Add Elon Musk biography" | `books add "Elon Musk" --author "Walter Isaacson" --category BIOGRAPHY --description "Biography of the entrepreneur behind Tesla, SpaceX, and other ventures." --status TO_BUY --json` |

Before adding, search for duplicates when the title is well-known. Try **both languages** when the book has common pt-BR and English editions:

```bash
books search "dune" --page 1 --limit 20 --json
books search "hobbit" --page 1 --limit 20 --json
books search "o hobbit" --page 1 --limit 20 --json
```

Merge by `id` before concluding the book is not in the library.

## Search

Search matches title and description substrings. The DB mixes **pt-BR and English** — run separate queries per language variant and merge by `id`.

| User says | Command |
| --- | --- |
| "Do I already have Dune?" | `books search "dune" --page 1 --limit 20 --json` |
| "Do I have The Hobbit?" | `books search "hobbit" --page 1 --limit 20 --json` **and** `books search "o hobbit" --page 1 --limit 20 --json` — merge by `id` |
| "Tenho O Senhor dos Anéis?" | `books search "senhor" --page 1 --limit 20 --json` **and** `books search "lord" --author "tolkien" --page 1 --limit 20 --json` |
| "Any book with 'senhor' in the title?" | `books search "senhor" --page 1 --limit 20 --json` |
| "Books about Arrakis" | `books search "arrakis" --page 1 --limit 20 --json` (matches description) |
| "Find books by Le Guin" | `books search "guin" --page 1 --limit 20 --json` (empty query is rejected) |
| "Search 1984 page 2" | `books search "1984" --page 2 --limit 20 --json` |
| "Dune by Herbert" | `books search "dune" --author "herbert" --page 1 --limit 20 --json` |

## List

Status filters are language-agnostic — no bilingual variants needed.

| User says | Command |
| --- | --- |
| "What am I reading?" | `books list --status READING --page 1 --limit 20 --json` |
| "Show my wishlist" | `books list --status TO_BUY --page 1 --limit 20 --json` |
| "Priority books to buy" | `books list --status TO_BUY --priority --page 1 --limit 20 --json` |
| "Books I haven't started" | `books list --status NOT_STARTED --page 1 --limit 20 --json` |
| "All my books" | `books list --page 1 --limit 20 --json` — paginate if `total > limit` |
| "Fiction on my wishlist about war" | `books list --status TO_BUY --page 1 --limit 20 --json` **plus** `books search "war" --page 1 --limit 20 --json` and `books search "guerra" --page 1 --limit 20 --json` — intersect with `TO_BUY` / `FICTION` in the merged set |

## Update & other

| User says | Command |
| --- | --- |
| "Mark book 42 as read" | `books update 42 --status READ --json` |
| "Recategorize book 42 as biography" | `books update 42 --category BIOGRAPHY --json` |
| "Refresh description for book 42" | Re-fetch synopsis in the book's title language, then `books update 42 --description "..." --json` |
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
