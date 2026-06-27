---
name: github-releases
description: >-
  Create and publish simple GitHub releases for books-cli using semantic
  versioning. Use when cutting a release, tagging a version, writing release
  notes, or updating CHANGELOG.md for a new version.
---

# GitHub Releases (books-cli)

Simple release workflow for this project: semver tag, changelog notes, and a Linux amd64 binary attached to every GitHub release.

## Versioning

- Follow [Semantic Versioning](https://semver.org/): `MAJOR.MINOR.PATCH` (e.g. `v0.1.0`).
- Tag format: `v` prefix required (`v0.1.0`, not `0.1.0`).
- Pre-1.0 (`0.x.y`): breaking changes may bump MINOR; PATCH for fixes and small features.

## Before releasing

1. Ensure CI is green on `main`.
2. Update `CHANGELOG.md`:
   - Move items from `[Unreleased]` into a new section: `## [x.y.z] - YYYY-MM-DD`.
   - Leave an empty `[Unreleased]` section at the top.
3. Commit the changelog update (message e.g. `chore: release v0.1.0`).

## Tag the release

Tag the changelog commit **before** pushing or publishing the GitHub release:

```bash
VERSION=v0.1.0

git tag -a "$VERSION" -m "$VERSION"
```

Verify the tag points at the right commit: `git show "$VERSION" --no-patch`.

## Create the release

```bash
VERSION=v0.1.0
ASSET=books-linux-amd64

# Push commit and tag
git push origin main
git push origin "$VERSION"

# Build release binary (linux/amd64) from the tagged commit
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o "$ASSET" ./cmd/books

# Create GitHub release with binary attached
gh release create "$VERSION" \
  --title "$VERSION" \
  --notes "$(awk -v v="${VERSION#v}" '/^## \['v'\]/{flag=1; next} /^## \[/{flag=0} flag' CHANGELOG.md)

See CHANGELOG.md for details." \
  "$ASSET"

rm -f "$ASSET"
```

Every release must include the built binary. Do not publish a release without uploading `books-linux-amd64`.

Add extra GOOS/GOARCH builds only when explicitly requested; default is linux/amd64 only.

## Release notes content

Keep notes short:

- 3–5 bullet points of user-visible changes.
- Link to full changelog: `See CHANGELOG.md for details.`

## Checklist

- [ ] CHANGELOG updated for this version
- [ ] Annotated tag created locally (`git tag -a vX.Y.Z -m vX.Y.Z`)
- [ ] Commit and tag pushed to `origin`
- [ ] `books-linux-amd64` built from the tagged commit
- [ ] GitHub release created with matching notes and binary uploaded
