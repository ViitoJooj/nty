# Contributing to nty

Thanks for helping out. This guide gets you from clone to pull request.

## Prerequisites

- [Go](https://go.dev/dl/) (see the version in `go.mod`)
- `git` and the [GitHub CLI](https://cli.github.com/) (`gh`) on your `PATH`

## Getting started

```sh
git clone git@github.com:ViitoJooj/nty.git
cd nty
go build ./...
```

Read the **[developer documentation](docs/README.md)** first — it explains the
architecture and the one rule that matters most.

## The one rule

**`cmd/` only wires; all logic lives in `internal/`.**

A Cobra handler reads flags and calls **one** function in `internal/services`.
If you're writing a loop or branching logic inside a `Run:` handler, it's in the
wrong place. See `docs/README.md` for the full layout.

## Before you open a PR

Run all three and make sure they pass:

```sh
go build ./...
go vet ./...
go test ./...
```

If you add non-trivial logic (a branch, a loop, a parser, a security path),
add at least one test that fails if that logic breaks. No need for heavy
frameworks — a plain `Test*` function is enough (see
`internal/services/release_flow_test.go`).

## Commit messages

Use [Conventional Commits](https://www.conventionalcommits.org/): `feat:`,
`fix:`, `chore:`, `docs:`, etc. (`nty release` reads these to decide version
bumps and write release notes.) Fittingly, you can use `nty commit` itself.

## Pull requests

1. Branch off the default branch.
2. Keep the diff focused — one logical change per PR.
3. Describe what changed and why.
4. Make sure build, vet, and tests are green.

## Safety expectations

When adding features, keep `nty` conservative:

- Destructive git operations should refuse to lose unmerged work
  (use `git branch -d`, not `-D`).
- Outward-facing or hard-to-undo actions (publishing, deleting) are
  **dry-run by default** and gated behind an explicit flag.
- **Never execute AI-generated shell commands.**

## License

By contributing, you agree that your contributions are licensed under the
[MIT License](LICENSE).
