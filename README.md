
# DevTool

DevTool is a Go-based command-line utility for repository inspection, GitHub analytics, branch maintenance, and project scaffolding. It helps developers work efficiently from the terminal without switching between Git commands, GitHub pages, and local project setup tools.

---

## What DevTool Can Do

- Inspect repository activity and health
- Show contributor and pull request metrics
- List, review, and clean local branches
- Analyze repository size and language breakdown
- Create starter projects for common frameworks
- Generate shell completion scripts

---

## Requirements

- Git installed and available on your PATH
- Go 1.24+ if you want to build from source
- A Git repository for repository analysis commands
- A GitHub Personal Access Token for GitHub-backed metrics

---

# Installing DevTool

DevTool can be installed using **prebuilt binaries from GitHub Releases**, so users **do not need Go installed**.

Download the latest release:

https://github.com/johannesalemu01/devtool/releases/latest

Each release contains binaries for the main platforms.

| Platform | Binary |
|--------|--------|
| macOS | devtool-macos |
| Linux | devtool-linux |
| Windows | devtool.exe |

---

## Install on macOS

```bash
curl -L https://github.com/johannesalemu01/devtool/releases/latest/download/devtool-macos -o devtool
chmod +x devtool
sudo mv devtool /usr/local/bin/devtool
```

Verify installation:

```bash
devtool --help
```

---

## Install on Linux

```bash
curl -L https://github.com/johannesalemu01/devtool/releases/latest/download/devtool-linux -o devtool
chmod +x devtool
sudo mv devtool /usr/local/bin/devtool
```

Verify:

```bash
devtool --help
```

---

## Install on Windows

Download:

https://github.com/johannesalemu01/devtool/releases/latest/download/devtool.exe

Then either run directly:

```
devtool.exe
```

Or place the file in a folder that exists in your **PATH** so it can be executed globally.

Example location:

```
C:\Program Files\devtool\
```

Then run:

```
devtool --help
```

---

## Install with Go (Optional)

If Go is installed:

```bash
go install github.com/johannesalemu01/devtool@latest
```

Make sure `$GOPATH/bin` is included in your PATH.

---

## Quick Start

After installation, navigate to any Git repository and run:

```bash
devtool contributors
```

or

```bash
devtool list-branches
```

or

```bash
devtool repo-size
```

---

## Verify Installation

```bash
devtool --help
```

Expected output:

```
Usage:
  devtool [command]

Available Commands:
  contributors
  repo-size
  list-branches
  clean-branches
  init
  dashboard
```

---

## Built With

- Go
- Cobra CLI framework
- Lip Gloss for terminal styling
- Git CLI
- GitHub REST API

---

## License

MIT License

---

## Author

Yohannes Alemu  
https://github.com/johannesalemu01
