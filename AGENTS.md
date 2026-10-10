# AGENTS.md

This file provides guidance to coding agents (Claude Code, Codex and others) when working with code in this repository.

## What this repository is

A mirror of Claude Code's native binaries on GitHub Releases, for machines that can reach GitHub but not `downloads.claude.ai`. There is no product code: the releases are the product, and the repository holds one script, the workflows that publish with it, three READMEs and a contributing guide.

The repository is also meant to become a template for mirroring other projects' binaries. That is why everything specific to Claude Code lives in `scripts/fetch.sh` and the READMEs, and nowhere else.

## How a release gets published

`.github/workflows/updater.yml` runs hourly and knows nothing about Claude Code. It calls three commands of `scripts/fetch.sh` and publishes whatever `download` leaves in `dist/`:

- `version` prints `latest` from upstream, rejecting anything that does not look like a version, because an error page served with 200 would otherwise become a tag.
- `download VERSION DIR` takes the file list from upstream's `manifest.json` (plain binaries) and `manifest.zst.json` (their zstd copies) rather than a hard-coded platform list, plus both manifests' signatures and the signing key. Each binary is named `claude-<version>-<platform>` followed by whatever its upstream name has after `claude` (`.exe`, `.zst`, `.exe.zst`). Users and old links depend on these names, so keep them stable. Any failed download or checksum mismatch must exit non-zero: a partial release is worse than none, because the next hourly run retries only while no release exists.
- `notes VERSION` prints that version's `## VERSION` section of anthropics/claude-code's CHANGELOG.md, falling back to a link to the changelog.

The workflow publishes only when no *published* release carries the tag. The check uses `gh api repos/$GITHUB_REPOSITORY/releases/tags/$VERSION` on purpose: `gh release view` also matches a draft left by a cancelled `gh release create` (which drafts, uploads, then publishes), and treats an empty tag as "the latest release", either of which would silently stop the mirror. For the same reason the version step assigns `VERSION=$(...)` before writing `$GITHUB_OUTPUT`; `echo "version=$(cmd)"` hides a failing `cmd` under `bash -e`.

A pull request touching `scripts/**` or `updater.yml` runs the full download (about 2.5 GB) without publishing; that run is the test. The publish step itself cannot run before merge, so after changing it, check the next real release: 23 assets and the right notes.

## Commands

```bash
./scripts/fetch.sh version
./scripts/fetch.sh download "$(./scripts/fetch.sh version)" dist  # ~2.5 GB; dist/ is gitignored
./scripts/fetch.sh notes 2.1.283
uvx pre-commit run -a  # what the Code Quality Check workflow runs
```

`fetch.sh` needs bash, curl, jq and sha256sum; it only has to run on Ubuntu runners and Linux. pre-commit only sees tracked files, so `git add` a new file before running it.

## Things that look simplifiable and are not

- The READMEs install by copying the binary to `~/.local/bin/claude` and set `DISABLE_AUTOUPDATER=1`. `claude install` re-downloads from upstream and fails offline, which is exactly where this mirror is used.
- The READMEs quote the signing key fingerprint unspaced, as `gpg --show-keys` prints it; codespell reads the last group of the spaced form as a misspelling of "cache".
- mdformat re-aligns the tables in the Chinese READMEs and URL-encodes CJK anchors. Let it.
- The repository setting "Allow GitHub Actions to create and approve pull requests" must stay on: Dependabot auto-approve and merge (`auto_review_merge.yml`) and `pre-commit-updater.yml` fail without it.

## Conventions

- The three READMEs (`README.md`, `README.zh-TW.md`, `README.zh-CN.md`) say the same thing; change all three together.
- Every workflow except `updater.yml`, plus `.pre-commit-config.yaml`, `.gitattributes`, `.gitignore` and `.github/{CODEOWNERS,dependabot.yml,labeler.yml}`, is shared with the Mai0313 template repositories. A change to one of them belongs in the templates too.
- Commit messages and PR titles are English Conventional Commits; `semantic-pull-request.yml` enforces the title. `labeler.yml` labels a PR by its branch prefix (`feat/`, `fix/`, `docs/`, `refactor/`, `chore/`, ...), so a branch without one gets no type label.
