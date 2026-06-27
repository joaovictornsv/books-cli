# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- SQLite-backed `books` table with embedded migrations and schema validation
- Commands: `add`, `get`, `list`, `search`, `update`, `archive`, `config`
- Status and boolean field validation with automatic `started_at` / `finished_at` timestamps
- Human-readable table output and `--json` for scripting
- Configurable database path via `BOOKS_DB`, `~/.config/books/config.toml`, or default path
- Command reference in `docs/COMMANDS.md`
- GitHub Actions CI (test and build on push/PR)
