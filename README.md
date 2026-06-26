# nty

A small command-line helper that does your git busywork for you. It writes your
commit messages, manages your branches, and cuts releases — using AI so you
don't have to think about the wording.

Works the same on Windows, macOS, and Linux.

---

## What it does

- **Commits for you.** Run one command and `nty` looks at your changed files,
  writes a clear message for each one, commits them, and pushes.
- **Keeps branches tidy.** It reuses your current work branch instead of piling
  up new ones, and cleans up branches that have already been merged.
- **Cuts releases.** It figures out the next version number, writes the release
  notes, and publishes a GitHub release — all from your commit history.
- **Speaks your language.** Choose English or Portuguese; commit messages and
  release notes follow your choice.

---

## Install

You need [Go](https://go.dev/dl/) installed. Then, from the project folder:

```sh
make install
```

This builds `nty` and puts it on your `PATH`. Restart your terminal afterwards.

> On Windows it installs to `%USERPROFILE%\bin\.nty`; on macOS/Linux to
> `~/.local/bin`.

---

## First-time setup

`nty` uses Claude (Anthropic) to write text. Connect your account once:

```sh
nty config
```

Pick **Register AI credentials** in the menu and follow the browser login. Your
login is saved in `~/.nty/config.yaml` so you only do this once.

While you're there, you can also pick **Language** to choose English or
Portuguese.

---

## Everyday use

**Commit and push everything you changed:**

```sh
nty commit
```

It branches (or reuses your branch), writes a message per file, commits, and
pushes. You don't type a single commit message.

**Cut a release:**

```sh
nty release            # preview — shows the version and notes, changes nothing
nty release --publish  # actually creates the GitHub release
```

By default `release` only *shows* you what it would do. Add `--publish` when
you're happy with it.

---

## Settings

Everything lives in `~/.nty/config.yaml`. Open `nty config` any time to change
your language or re-connect your AI account.

---

## For developers

Building on `nty`, contributing, or want to know how it works under the hood?

➡️ **[Developer documentation](docs/README.md)**

- [Contributing guide](CONTRIBUTING.md)
- [Security policy](SECURITY.md)
- [License (MIT)](LICENSE)
