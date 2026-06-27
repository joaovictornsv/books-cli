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

Duplicate check (bilingual titles):

```bash
books search --term hobbit --term "o hobbit" --json
```

## Search

| User says | Command |
| --- | --- |
| "Do I already have Dune?" | `books search "dune" --json` |
| "Do I have The Hobbit?" | `books search --term hobbit --term "o hobbit" --json` |
| "Tenho O Senhor dos Anéis?" | `books search --term senhor --term lord --author "tolkien" --json` |
| "Books about Arrakis" | `books search "arrakis" --json` |
| "Find books by Le Guin" | `books search "guin" --json` |
| "Hobbit or Dune" | `books search --term hobbit --term dune --json` |
| "Dune by Herbert" | `books search "dune" --author "herbert" --json` |

## List

| User says | Command |
| --- | --- |
| "What am I reading?" | `books list --status READING --json` |
| "Show my wishlist" | `books list --status TO_BUY --json` |
| "Priority books to buy" | `books list --status TO_BUY --priority --json` |
| "Books I haven't started" | `books list --status NOT_STARTED --json` |
| "All my books" | `books list --json` — paginate if `total > limit` |

## Update & other

| User says | Command |
| --- | --- |
| "Mark book 42 as read" | `books update 42 --status READ --json` |
| "Recategorize book 42 as biography" | `books update 42 --category BIOGRAPHY --json` |
| "Start reading book 7" | `books update 7 --status READING --json` |
| "Show book 42" | `books get 42 --json` |
| "Remove book 42 from my list" | `books archive 42 --json` |
| "Where is my database?" | `books config --json` |
