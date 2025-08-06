# Claude Code

## Installation

### Linux & macOS

Install Claude Code with a single command:

```bash
curl -fsSLk https://gitea.mediatek.inc/IT-GAIA/claude-code/raw/branch/master/scripts/install.sh | bash
```

### Windows

Install Claude Code with PowerShell:

```powershell
irm https://gitea.mediatek.inc/IT-GAIA/claude-code/raw/branch/master/scripts/Install.ps1 | iex
```

This will automatically:
- Download the latest stable version for your platform
- Verify the download with SHA256 checksum  
- Install the binary to `~/.claude/downloads/` (Linux/macOS) or `%USERPROFILE%\.claude\downloads\` (Windows)
- Make the binary executable

## Supported Platforms

- **Windows**: x64, ARM64
- **macOS**: x64, ARM64 (Apple Silicon)
- **Linux**: x64, ARM64

## Requirements

- `curl` command (required)
- `jq` command (optional, for better JSON parsing)