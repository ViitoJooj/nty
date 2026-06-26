# nty — Developer Documentation

`nty` is a Go CLI built on [Cobra](https://github.com/spf13/cobra) for commands,
[Viper](https://github.com/spf13/viper) for config, and
[Bubble Tea](https://github.com/charmbracelet/bubbletea) for the interactive
config UI. It talks to the Anthropic (Claude) API to generate commit messages
and release notes.

- **Module:** `github.com/ViitoJooj/nty`
- **Entry point:** `main.go` → `cmd.Execute()`

---

## Architecture at a glance

The golden rule: **`cmd/` only wires; all logic lives in `internal/`.** A Cobra
handler should read flags and call one function. If you find yourself writing a
loop or branching logic inside a `Run:` handler, it belongs in `internal/services`.

```
main.go                     entry point
cmd/                        Cobra commands — thin wiring only
  root.go                   root command + Execute()
  commit.go                 `nty commit`  → services.RunCommit()
  config.go                 `nty config`  → tui.Run() / tui.RunLang()
  release.go                `nty release` → services.RunRelease()
internal/
  services/                 all business logic
    git.go                  git plumbing (run git, branches, tags)
    claude.go               Anthropic API client + prompts
    claude_oauth.go         Claude OAuth (PKCE) login flow
    commit_flow.go          RunCommit orchestration
    release_flow.go         RunRelease orchestration
  config/
    config.go               load/save ~/.nty/config.yaml (Viper)
  tui/                      Bubble Tea app
    app.go                  App model + screen switching
    run.go                  Run / RunLang launchers
    screens/                one file per screen
```

---

## The three commands

### `nty commit` — `services.RunCommit()`

1. List changed files (`git status --porcelain`).
2. Resolve the branch:
   - On a work branch → reuse it.
   - On a base branch (`main`/`master`/`develop`) → delete merged branches
     (`git branch -d`, never `-D`), then create a fresh AI-named branch.
3. For each file: stage, diff, ask Claude for a message, commit.
4. `git push -u origin <branch>`.

### `nty config` — `tui.Run()` / `tui.RunLang()`

Launches the Bubble Tea menu. `--lang` jumps straight to language selection.
Each menu option opens a screen under `internal/tui/screens/`. Selections are
persisted via `internal/config`.

### `nty release` — `services.RunRelease(publish bool)`

1. Last tag (`git describe --tags --abbrev=0`, fallback `v0.0.0`).
2. Commits since that tag.
3. **AI classifies** the project as `binary` or `library` from its manifests +
   file tree (language-independent — the model generalizes, no per-ecosystem code).
4. **AI decides** the next semver and writes the notes. The version is validated
   against `^v\d+\.\d+\.\d+$`; junk is rejected.
5. If `binary`, build an artifact from a **fixed command table** (Go only today)
   — never AI-generated shell.
6. Print the plan. **Dry-run unless `--publish`.** With `--publish`, call
   `gh release create`.

---

## Configuration

Stored at `~/.nty/config.yaml`, managed by Viper (`internal/config`).

| Key                       | Meaning                                  |
|---------------------------|------------------------------------------|
| `lang`                    | `en` or `pt` — drives AI output language |
| `claude_access_token`     | Claude OAuth access token                |
| `claude_refresh_token`    | Claude OAuth refresh token               |
| `claude_expires_at`       | token expiry (unix seconds)              |
| `github_user` / `_email`  | GitHub identity                          |

Add a setting: extend the `Config` struct and add a `SaveX` helper in
`config.go` following the existing `viper.Set` + `write()` pattern.

---

## AI integration

`internal/services/claude.go` wraps the Anthropic Messages API
(`claude-haiku-4-5`) via the OAuth token from config. The shared `complete()`
helper takes a system + user prompt; the public helpers build the prompts:

- `CommitMessage(file, diff)` — one-line commit message, language from `lang`.
- `BranchName(files)` — kebab-case branch name.
- `ClassifyArtifact(manifests, tree)` — `binary` vs `library`.
- `ReleaseNotes(lastTag, commits)` — next version + notes.

Login uses OAuth PKCE (`claude_oauth.go`): `Start()` returns a URL + PKCE,
`Exchange(code, pkce)` swaps the code for tokens, saved via
`config.SaveClaudeAuth`.

---

## Adding a command

1. Create `cmd/<name>.go` with a `cobra.Command`. The `Run:` handler reads flags
   and calls **one** `services` function — no logic here.
2. Put the logic in `internal/services/<name>_flow.go`.
3. Register it in `init()` with `rootCmd.AddCommand(...)`.

Outward-facing or hard-to-undo actions (publishing, deleting) should be
**dry-run by default** and gated behind an explicit flag, like `release --publish`.

---

## Build & test

```sh
go build ./...        # compile
go vet ./...          # static checks
go test ./...         # tests
make install          # build + install to PATH
```

---

## Conventions

- **Thin `cmd/`, fat `internal/services`.** Logic never lives in a Cobra handler.
- **Errors bubble up with `%w`;** the command layer prints them.
- **Destructive git is safe by default:** merged-branch cleanup uses `-d`
  (refuses unmerged), releases are dry-run until `--publish`.
- **Never run AI-generated shell.** Build commands come from a fixed table.

---

## See also

- [Contributing guide](../CONTRIBUTING.md) — setup, the one rule, PR checklist
- [Security policy](../SECURITY.md) — reporting, what `nty` touches on your machine
- [License (MIT)](../LICENSE)
