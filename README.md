## Golang CLI App Bootstrap

This project provides common bootstrapping and features used across most Go CLI apps.

### Features

- **CLI:** `urfave/cli/v2`
- **Structured logging:** `log/slog`
- **TOML configs:** `burntsushi/toml`
- **Environment files:** `joho/godotenv`
- Error and exit handling
- Logrotate config (under `data/logrotate.d/`)
- Sample `version` subcommand
- Config consolidation: config file → env file → env vars → CLI args
- Makefile: build, dev build, run, install, test, clean

### Prerequisites

- Go 1.21+
- Make

### Usage

#### Creating a new project

From the bootstrap repo root, run:

```bash
make bootstrap OWNER=<owner> APP=<appname>
```

Use your **owner/org name** (e.g. GitHub username or module path) and your **app name** (single word). Placeholders in the copied files are replaced with these values. Example:

```bash
make bootstrap OWNER=myorg APP=mycli
```
This will resolve the project path to `github.com/myorg/mycli/...`

The command will:

- Copy bootstrap files into the current directory
- Replace `appname` / `ownername` with your APP and OWNER
- Rename `cmd/appname` and the logrotate file to your app name
- Run `git init` and `go get` for dependencies

When you're done, remove the `.bootstrap` directory. You can delete it right away or keep it if you might reset and re-bootstrap later.

#### Resetting to re-bootstrap

To wipe the generated project and restore the repo so you can run `make bootstrap` again:

1. `cd .bootstrap`
2. `make clean`

This only works when run from inside the `.bootstrap` directory.

### AI Disclosure

This bootstrap cli project was hacked together with the assistance of AI. I make no guarantees about anything, this was only made to quickly bootstrap golang cli projects.
Is there a better way to do this? Yeah, but this works for me.
