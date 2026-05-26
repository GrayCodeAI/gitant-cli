# gitant-cli

Client for [Gitant](https://github.com/GrayCodeAI/gitant-daemon) — use from the terminal while a Gitant node runs locally or on a hosted URL.

Install **`gitant-cli`** on your laptop. Install **`gitant-daemon`** (or use Gitant Cloud) as the server.

**Full setup guide:** [gitant-daemon docs/QUICKSTART.md](https://github.com/GrayCodeAI/gitant-daemon/blob/main/docs/QUICKSTART.md)

## Install

### Pre-built binary (releases)

```bash
curl -fsSL https://raw.githubusercontent.com/GrayCodeAI/gitant-cli/main/scripts/install.sh | bash
```

### Go install

```bash
go install github.com/GrayCodeAI/gitant-cli/cmd/gitant@latest
go install github.com/GrayCodeAI/gitant-cli/cmd/git-remote-gitant@latest
```

### From source

```bash
git clone https://github.com/GrayCodeAI/gitant-cli.git
cd gitant-cli
make build
./bin/gitant version
```

## Quick start

Point at your node (self-hosted or hosted):

```bash
export GITANT_DAEMON_URL=http://localhost:7777   # or https://gitant.example.com

gitant doctor
gitant repo list
git init && git commit -m "init"
gitant push --remote "$GITANT_DAEMON_URL" --repo my-app
gitant issue create --repo my-app --title "First issue"
```

Or use the [web dashboard](https://github.com/GrayCodeAI/gitant-web) against the same daemon URL.

## What runs where

| Component | Repo | Role |
|-----------|------|------|
| **gitant** (this repo) | `gitant-cli` | Developer CLI — git push/pull, issues, PRs, agents |
| **gitant serve** | `gitant-daemon` | HTTP API + git smart HTTP + optional P2P |
| **Dashboard** | `gitant-web` | Browser UI |
| **Agents** | `gitant-mcp` | MCP tools + `@gitant/sdk` |

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `GITANT_DAEMON_URL` | `http://localhost:7777` | Base URL of your Gitant node |
| `GITANT_UCAN_TOKEN` | — | UCAN token for authenticated API calls |

Most commands accept `--daemon-url` to override the env var.

## Commands

```bash
gitant version
gitant doctor
gitant quickstart

# Git
gitant init
gitant push --remote URL --repo ID
gitant pull --remote URL --repo ID
gitant clone REPO [DIR] --remote URL

# Repos, issues, PRs, tasks, labels, protection, webhooks, agents, UCAN
gitant repo list
gitant issue list --repo ID
gitant pr create --repo ID --title T -s branch
gitant agent delegate --did DID --resource repo:ID --actions read,write

# Daemon data (when self-hosting)
gitant backup -o ./backups
gitant restore -i ./backups/gitant-backup-...
```

See `gitant --help` and [gitant-daemon README](https://github.com/GrayCodeAI/gitant-daemon) for the full API.

## git remote helper

`git-remote-gitant` enables `gitant://` remotes:

```bash
git remote add origin gitant://my-repo
export GITANT_DAEMON_URL=http://localhost:7777
git push origin main
```

## Development

```bash
make test
make build
```

This repo is self-contained. For local full-stack dev, clone it beside `gitant-daemon`, `gitant-web`, and `gitant-mcp` in a folder such as `gitant-core/`.

## License

MIT — see [LICENSE](LICENSE).
