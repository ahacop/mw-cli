# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

mw-cli is a command-line dictionary tool that fetches definitions from the Merriam-Webster Collegiate Dictionary API. It's a single-file Go application (main.go) that provides fast dictionary lookups with colored terminal output.

## Environment Setup

**Required environment variable:**
- `DICTIONARY_KEY` - Merriam-Webster API key (obtain free from dictionaryapi.com)

The project uses Nix flakes with direnv for reproducible development environments:
```bash
direnv allow              # Auto-loads environment
nix develop               # Or enter dev shell manually
```

## Common Commands

The project uses `just` as a command runner. All commands are defined in `justfile`:

```bash
# Development workflow
just build                # Build binary: go build -o mw-cli main.go
just run <word>           # Build and run with a word
just dev <word>           # Alias for run

# Testing and quality
just test                 # Run tests with: go test -v ./...
just test-coverage        # Generate coverage.html report
just fmt                  # Format code with go fmt
just vet                  # Run go vet
just lint                 # Run golangci-lint (requires golangci-lint in PATH)
just check                # Run fmt, vet, lint, and test together

# Installation and cleanup
just install              # Install to $GOPATH/bin
just clean                # Remove build artifacts and coverage files
just deps                 # Download and tidy dependencies
```

## Architecture

**Single-file application** (main.go ~125 lines):

1. **Main Flow**: CLI argument parsing → API request → JSON parsing → colored output
2. **API Integration**: Makes GET requests to Merriam-Webster's `/api/v3/references/collegiate/json` endpoint
3. **Error Handling**:
   - API returns JSON array of `DictionaryEntry` structs when word is found
   - API returns JSON array of strings (suggestions) when word is misspelled
   - The code attempts to unmarshal as DictionaryEntry first, then falls back to suggestions
4. **Output Formatting**: Uses `github.com/fatih/color` for terminal colors (cyan for words, yellow for pronunciation, green for definitions)

**Key Structs:**
- `DictionaryEntry`: Represents API response with nested Meta (id, functional label), Hwi (headword, pronunciations), and Shortdef (short definitions) fields

## Release Process

Releases are automated via GitHub Actions (.github/workflows/release.yml):

1. Create and push a git tag: `git tag v0.1.0 && git push origin v0.1.0`
2. GitHub Actions automatically runs goreleaser to:
   - Build for Linux/macOS/Windows on amd64/arm64/arm architectures
   - Create GitHub release with binaries and checksums
   - Package binaries as tar.gz (Unix) or zip (Windows)

**GoReleaser config** (.goreleaser.yaml):
- Builds with `CGO_ENABLED=0` for static binaries
- Uses ldflags `-s -w` to strip debug info and reduce binary size
- Excludes Windows arm64/arm builds
- Groups changelog by Features/Bug fixes using conventional commit prefixes

## Nix Integration

The flake.nix provides:
- Package definition using `buildGoModule` with vendorHash for reproducible builds
- Dev shell with Go, gopls, golangci-lint, and just pre-installed
- Version currently set to "0.1.0" in flake.nix (update when releasing)
- Users can run directly: `nix run github:ahacop/mw-cli -- <word>`
