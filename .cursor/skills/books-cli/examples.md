# books-cli examples

Phrase → command. Always append `--json`. See [SKILL.md](SKILL.md) for agent rules (description, category, defaults).

## Add

| User says | Command |
| --- | --- |
| "Add Dune to my wishlist" | `books add "Dune" --author "Frank Herbert" --category FICTION --description "..." --status TO_BUY --json` |
| "Adiciona O Hobbit" | `books add "O Hobbit" --author "J.R.R. Tolkien" --category FICTION --description "..." --status TO_BUY --json` (pt-BR description) |
| "Add 1984, priority" | `... --status TO_BUY --priority --json` |
| "Add Sapiens as reading" | `... --category HISTORY --status READING --json` |
| "Add Elon Musk biography" | `... --category BIOGRAPHY --json` |

Duplicate check before add:

```bash
books check --title "Dune" --author "Herbert" --json
books check --title "O Hobbit" --json
```

Cross-language fallback (when titles differ by language):

```bash
books search --term hobbit --term "o hobbit" --json
```

## Search

| User says | Command |
| --- | --- |
| "Do I already have Dune?" | `books check --title "Dune" --json` |
| "Do I have The Hobbit?" | `books check --title "Hobbit" --json` or bilingual `search --term hobbit --term "o hobbit" --json` |
| "Tenho O Senhor dos Anéis?" | `books search --term senhor --term lord --author "tolkien" --json` |
| "Books about Arrakis" | `books search "arrakis" --json` |
| "Find books by Le Guin" | `books search "guin" --json` |
| "Find books by Herbert" | `books search --term herbert --json` |
| "Hobbit or Dune" | `books search --term hobbit --term dune --json` |
| "Dune by Herbert" | `books search "dune" --author "herbert" --json` |

## List

| User says | Command |
| --- | --- |
| "What am I reading?" | `books list --status READING --json` |
| "Show my wishlist" | `books list --status TO_BUY --json` |
| "Fiction on my wishlist" | `books list --status TO_BUY --category FICTION --json` |
| "Priority books to buy" | `books list --status TO_BUY --priority --json` |
| "Books I haven't started" | `books list --status NOT_STARTED --json` |
| "All my books" | `books list --json` — paginate if `total > limit` |
| "Recently finished" | `books list --status READ --sort finished_at --order desc --limit 10 --json` |
| "Compact wishlist" | `books list --status TO_BUY --fields id,title,status --json` |

## Count & stats

| User says | Command |
| --- | --- |
| "How many books am I reading?" | `books count --status READING --json` |
| "How many on my wishlist?" | `books count --status TO_BUY --json` |
| "Library stats" | `books stats --json` |
| "How many did I finish this year?" | `books stats --year 2025 --json` |

## Update & other

| User says | Command |
| --- | --- |
| "Mark book 42 as read" | `books update 42 --status READ --finished-at "<RFC3339>" --json` |
| "Mark books 1, 2, 3 as read" | `books update --ids 1,2,3 --status READ --json` |
| "Recategorize book 42 as biography" | `books update 42 --category BIOGRAPHY --json` |
| "Start reading book 7" | `books update 7 --status READING --started-at "<RFC3339>" --json` |
| "Remove priority from book 3" | `books update 3 --no-priority --json` |
| "Show book 42" | `books get 42 --json` |
| "Show Dune" | `books get --title "Dune" --exact --json` |
| "Remove book 42 from my list" | `books update 42 --status ARCHIVED --json` |
| "Delete book 42 permanently" | `books delete 42 -y --json` |
| "Export my library" | `books export --format json --output books.json --json` |
| "Import books from a file" | `books import --input books.json --json` |
| "Backup my database" | `books backup --output /path/to/backup.db --json` |
| "What statuses exist?" | `books schema --json` |
| "Where is my database?" | `books config --json` |
