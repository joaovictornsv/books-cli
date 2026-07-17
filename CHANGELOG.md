# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- Renamed `eligible_to_sell` / `sold` to `eligible_to_donate` / `donated` (DB columns, JSON keys, CLI flags, and table output)

### Removed

- `--eligible-to-sell`, `--no-eligible-to-sell`, `--sold`, `--no-sold` flags (replaced by donate equivalents)

## [0.6.0] - 2026-07-09

### Added

- `version` command with `--json` build metadata (`version`, `commit`, `go_version`)
- `SHA256SUMS` checksum file published alongside release binaries; README documents verification steps

## [0.5.0] - 2026-07-07

### Added

- Bulk `update --ids` to change several books in one transaction
- `export` command to dump the library to JSON or CSV (`--format`, `--output`, `--include-archived`)
- `import` command to load books from JSON or CSV with upsert-by-ID and `--dry-run`

## [0.4.0] - 2026-07-06

### Added

- `--sort` and `--order` flags on `list` and `search` (fields: `id`, `title`, `author`, `status`, `added_at`, `started_at`, `finished_at`)
- `search` `--term` now matches author in addition to title and description (`--author` remains an AND filter)
- `get --title` to look up a book by title (`--exact`, `--author` for disambiguation)

## [0.3.0] - 2026-07-03

### Added

- `schema` command for machine-readable status/category enums and book field semantics
- `check` command for pre-add duplicate detection by title (and optional author)
- `--category` filter on `search`
- `--no-priority`, `--no-eligible-to-sell`, and `--no-sold` flags on `update` to clear boolean fields explicitly
- `delete` requires `--yes` / `-y` when using `--json`; interactive mode prompts on a TTY unless `-y` is passed

## [0.2.0] - 2026-07-01

### Added

- `count` command with optional `--status`, `--category`, `--priority`, and `--eligible-to-sell` filters
- `stats` command with aggregates by status/category, `finished_this_year`, and `priority_wishlist`
- `backup` command to create a consistent SQLite database copy via `VACUUM INTO` (`--output`, `--force`)
- JSON output for `config` via `--json`
- `--category` filter on `list`

### Removed

- `archive` command; use `update <id> --status ARCHIVED` instead

## [0.1.0] - 2026-06-27

### Added

- Repeatable `--term` on `search` for OR matching across title/description (positional query still supported)
- `description` column on books (optional text, searchable via `search`)
- `delete` command to permanently remove a book from the database
- SQLite-backed `books` table with embedded migrations and schema validation
- Commands: `add`, `get`, `list`, `search`, `update`, `archive`, `delete`, `config`
- Status and boolean field validation with automatic `started_at` / `finished_at` timestamps
- Human-readable table output and `--json` for scripting
- Configurable database path via `BOOKS_DB`, `~/.config/books/config.toml`, or default path
- Command reference in `docs/COMMANDS.md`
- Pagination for `list` and `search` (`--page`, `--limit`)
- GitHub Actions CI (test and build on push/PR)
