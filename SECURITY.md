# Security Policy

## Reporting a vulnerability

If you find a security issue in `nty`, please **do not open a public issue**.

Email the maintainer privately: **joaovitor819oqueres@gmail.com**

Include:

- a description of the issue and its impact,
- steps to reproduce (or a proof of concept),
- the `nty` version / commit you tested.

You'll get an acknowledgement as soon as possible. Please give a reasonable
window to fix the issue before disclosing it publicly.

## Supported versions

`nty` is pre-1.0. Only the latest commit on the default branch receives security
fixes.

## What `nty` touches on your machine

Understanding this helps you assess risk:

- **Credentials at rest.** Your Claude OAuth tokens and GitHub identity are
  stored in `~/.nty/config.yaml` in **plain text**. Protect that file with your
  OS file permissions; anyone who can read it can use your Claude session.
- **Network calls.** `nty` sends file diffs and commit messages to the Anthropic
  (Claude) API to generate text. Don't run `nty commit` on a repository whose
  diffs you are not comfortable sending to a third-party API.
- **It runs git and `gh`.** `nty` shells out to your local `git` and the GitHub
  CLI. It uses whatever credentials those tools already have.
- **It never runs AI-generated shell.** Release build commands come from a fixed,
  reviewed table in the source — the AI classifies the project but does not get
  to execute arbitrary commands.

## Destructive operations

`nty` is conservative by design:

- Merged-branch cleanup uses `git branch -d` (which refuses unmerged branches),
  never `-D`.
- `nty release` is a **dry-run** until you pass `--publish`.
