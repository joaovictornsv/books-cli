# books-cli

A personal CLI to manage your book collection in a local SQLite database. Track reading status, wishlists, notes, and search your library from the terminal — single binary, no server.

```bash
books add "Dune" --author "Frank Herbert" --status NOT_STARTED
books list --status READING
books search "le guin" --json
```

## Commands

| Command | Description |
| --- | --- |
| `add` | Add a book |
| `get` | Show one book by ID or title |
| `list` | List books with filters, sorting, and pagination |
| `search` | Search by title, description, or author |
| `check` | Find likely duplicates by title before adding |
| `update` | Update fields or reading status (including `--status ARCHIVED`) |
| `delete` | Permanently remove a book |
| `schema` | Show enums and field semantics (`--json`) |
| `count` | Count books matching filters |
| `stats` | Library aggregates |
| `backup` | Copy the SQLite database |
| `export` | Dump library to JSON or CSV |
| `import` | Load books from JSON or CSV |
| `config` | Show resolved configuration |

Use `--json` on any command for scripting. Full flag reference: [docs/COMMANDS.md](docs/COMMANDS.md). For AI agents: [`.cursor/skills/books-cli/SKILL.md`](.cursor/skills/books-cli/SKILL.md).

## Setup

**Requirements:** Go 1.25+

```bash
git clone https://github.com/joaovictornsv/books-cli.git
cd books-cli
go build -o books ./cmd/books
```

Pre-built binaries for linux/amd64 are available on [GitHub Releases](https://github.com/joaovictornsv/books-cli/releases).

### Database path

1. `BOOKS_DB` environment variable
2. `database` in `~/.config/books/config.toml`
3. Default: `~/.local/share/books/books.db`

```toml
database = "/home/user/books.db"
```

Run `books config` to see which path is in use.

## Development

```bash
go test ./...
go build ./cmd/books
```

Changes are tracked in [CHANGELOG.md](CHANGELOG.md).

## License

MIT — see [LICENSE](LICENSE).
