# Stoic

## Description

Keep calm and think in the stormy terminal.

## Installation

### Direct Download from GitHub Releases
The easiest way to install Stoic is to download the pre-built binary for your platform from the [GitHub Releases page](https://github.com/narcilee7/stoic/releases).

1. Go to the latest release.
2. Download the file matching your OS and architecture (e.g., `stoic_Darwin_arm64.tar.gz` for macOS on Apple Silicon).
3. Extract the archive:
   - For tar.gz (macOS/Linux): `tar -xzf stoic_*.tar.gz`
   - For zip (Windows): Unzip using File Explorer or `unzip`.
4. Run the executable: `./stoic` (or `stoic.exe` on Windows).

### Homebrew (macOS)
If you're on macOS, you can install via Homebrew:

```bash
brew install narcilee7/stoic/stoic
```

(Note: This assumes a Homebrew tap is set up; see Contributing for details on creating one.)

### Building from Source
If you prefer, build from source:

1. Clone the repo: `git clone https://github.com/narcilee7/stoic.git`
2. `cd stoic`
3. `go build ./cmd/main.go`
4. Run `./main` (rename to stoic if desired).
