# Merriam-Webster CLI

A command-line dictionary tool that fetches definitions from Merriam-Webster's Collegiate Dictionary API.

## Features

- Fast dictionary lookups from the command line
- Colored output for better readability
- Shows word pronunciations
- Displays multiple definitions
- Suggests alternative spellings for misspelled words

## Installation

### Pre-built Binaries (Recommended)

Download pre-built binaries from the [releases page](https://github.com/ahacop/mw-cli/releases):

```bash
# Example for Linux (replace with your platform)
wget https://github.com/ahacop/mw-cli/releases/download/v0.1.0/mw-cli_v0.1.0_linux_amd64.tar.gz
tar -xzf mw-cli_v0.1.0_linux_amd64.tar.gz
sudo mv mw-cli /usr/local/bin/
```

Available platforms: Linux, macOS, Windows (amd64, arm64, arm)

### Go Install

If you have Go installed:

```bash
go install github.com/ahacop/mw-cli@latest
```

The binary will be installed to `$GOPATH/bin/mw-cli` (typically `~/go/bin/mw-cli`).

### Nix

If you use Nix:

```bash
# Run directly without installing
nix run github:ahacop/mw-cli -- serendipity

# Or install to your profile
nix profile install github:ahacop/mw-cli
```

### NixOS (Flake-based Configuration)

To install mw-cli system-wide in your NixOS configuration:

**Option 1: Direct package reference (recommended)**

```nix
{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    mw-cli.url = "github:ahacop/mw-cli";
    # ... other inputs
  };

  outputs = { self, nixpkgs, mw-cli, ... }: {
    nixosConfigurations.your-hostname = nixpkgs.lib.nixosSystem {
      system = "x86_64-linux";  # or your architecture
      modules = [
        ({ pkgs, ... }: {
          environment.systemPackages = [
            mw-cli.packages.${pkgs.system}.default
            # ... other packages
          ];

          # Set API key system-wide
          environment.variables.DICTIONARY_KEY = "your-api-key-here";
        })
      ];
    };
  };
}
```

**Option 2: With Home Manager**

```nix
{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    home-manager.url = "github:nix-community/home-manager";
    mw-cli.url = "github:ahacop/mw-cli";
  };

  outputs = { self, nixpkgs, home-manager, mw-cli, ... }: {
    homeConfigurations.your-user = home-manager.lib.homeManagerConfiguration {
      pkgs = nixpkgs.legacyPackages.x86_64-linux;
      modules = [
        ({ pkgs, ... }: {
          home.packages = [
            mw-cli.packages.${pkgs.system}.default
          ];

          # Set API key for user
          home.sessionVariables.DICTIONARY_KEY = "your-api-key-here";
        })
      ];
    };
  };
}
```

After adding to your configuration, rebuild with `sudo nixos-rebuild switch --flake .#your-hostname`.

**Note:** For better security, consider using [agenix](https://github.com/ryantm/agenix) or [sops-nix](https://github.com/Mic92/sops-nix) to manage the API key secret instead of hardcoding it.

### Build from Source

```bash
git clone https://github.com/ahacop/mw-cli.git
cd mw-cli
go build -o mw-cli
sudo cp mw-cli /usr/local/bin/
```

## Setup

1. **Get a Merriam-Webster API Key** (free)
   - Visit [https://dictionaryapi.com/](https://dictionaryapi.com/)
   - Register for a free API key for the Collegiate Dictionary

2. **Set your API key**
   ```bash
   export DICTIONARY_KEY=your-api-key-here
   ```

   Add this to your shell's rc file (~/.bashrc, ~/.zshrc, etc.) to make it permanent.

## Usage

```bash
mw-cli serendipity
```

### Example Output

```
serendipity [noun]
  /ser-ən-ˈdi-pə-tē/
  1. luck that takes the form of finding valuable or pleasant things that are not looked for
  2. the faculty or phenomenon of finding valuable or agreeable things not sought for
```

## Development

### Prerequisites

- Go 1.25.2 or later
- [just](https://github.com/casey/just) (optional, for convenient commands)

### Development Commands

If you have `just` installed:

```bash
just --list              # Show all available commands
just build               # Build the application
just run <word>          # Build and run with a word
just test                # Run tests
just test-coverage       # Run tests with coverage report
just fmt                 # Format code
just lint                # Run linter (requires golangci-lint)
just vet                 # Run go vet
just check               # Run all checks (fmt, vet, lint, test)
just clean               # Remove build artifacts
just deps                # Download and tidy dependencies
just install             # Install binary to GOPATH/bin
```

Without `just`:

```bash
go build -o mw-cli       # Build
go run main.go <word>    # Run
go test -v ./...         # Test
go fmt ./...             # Format
go vet ./...             # Vet
```

### Nix Development Environment (Optional)

The project includes a Nix flake for reproducible development environments:

```bash
# Enter the development shell (includes Go, gopls, golangci-lint, just)
nix develop

# Or use direnv for automatic loading
direnv allow
```

### Creating a Release

To create a new release:

1. Update version in `flake.nix` if needed
2. Create and push a git tag:
   ```bash
   git tag v0.1.0
   git push origin v0.1.0
   ```
3. GitHub Actions will automatically:
   - Build binaries for Linux, macOS, and Windows
   - Create a GitHub release
   - Upload pre-built binaries

## Credits

Inspired by [define-rs](https://github.com/fosslife/define) - a similar tool written in Rust.

## License

GPLv3
