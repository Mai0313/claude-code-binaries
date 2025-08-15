# Claude Code CLI Usage Guide

English | [繁體中文](README.zh-TW.md) | [简体中文](README.zh-CN.md)

## Overview

Claude Code is Anthropic's official CLI tool for AI-powered coding assistance with interactive development support.

---

## Obtaining API Key

**Important**: Before using Claude Code, complete these steps:

1. Visit [MediaTek MLOp Gateway](https://mlop-azure-gateway.mediatek.inc/auth/login) to log in
2. Obtain your GAISF API key
3. Keep the key secure for subsequent use

**Note**: Due to SSL certificate configuration, the URLs in this documentation use HTTP instead of HTTPS for compatibility across different network environments.

---

## Installation

### Option 1: Using Pre-built Binaries (Recommended)

Visit [Claude Code Release](https://gitea.mediatek.inc/IT-GAIA/claude-code/releases/latest) to download the latest installer package for your platform.

#### Windows

⚠️ **Important Warning**: Claude Code has limited support for Windows. If you cannot use WSL (Windows Subsystem for Linux), we strongly recommend using macOS or Linux instead for the best experience.

1. Download the Windows `.exe` binary file
2. Create a directory for the binary:
   ```cmd
   mkdir %USERPROFILE%\.local\bin
   ```
3. Move the downloaded binary to the directory and rename it:
   ```cmd
   move claude-code-windows-x64.exe %USERPROFILE%\.local\bin\claude.exe
   ```
4. Add the directory to your PATH environment variable:
   - Open System Properties → Environment Variables
   - Add `%USERPROFILE%\.local\bin` to your PATH
   - Or use PowerShell:
   ```powershell
   $env:PATH += ";$env:USERPROFILE\.local\bin"
   [Environment]::SetEnvironmentVariable("PATH", $env:PATH, "User")
   ```

#### macOS/Linux
1. Download the appropriate binary for your platform
2. Move the binary to your local bin directory:

```bash
# Make it executable
chmod +x claude-code-*

# Move to ~/.local/bin (create directory if it doesn't exist)
mkdir -p ~/.local/bin
mv claude-code-* ~/.local/bin/claude

# Ensure ~/.local/bin is in your PATH
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

### Option 2: Using npm (For developers)

Install Claude Code using npm (requires Node.js and programming experience):

```bash
npm install -g @anthropic-ai/claude-code
```

---

## Platform Support

Claude Code supports the following platforms:
- macOS
- Linux 
- Windows WSL (Windows Subsystem for Linux)

### Windows Users

For Windows users who want to use Claude Code outside of WSL, you'll need to install Git Bash:

1. Download and install Git for Windows from: https://git-scm.com/downloads/win
2. Set the environment variable to point to Git Bash:
   ```powershell
   $env:CLAUDE_CODE_GIT_BASH_PATH="C:\Program Files\Git\bin\bash.exe"
   ```

For more detailed setup information, see the [official setup documentation](https://docs.anthropic.com/en/docs/claude-code/setup).

---

## Configuration

Create a configuration file at `~/.claude/settings.json`:

```json
{
  "cleanupPeriodDays": 30,
  "enableAllProjectMcpServers": true,
  "includeCoAuthoredBy": true,
  "permissions": {
    "allow": [
      "Bash(npm run lint)",
      "Bash(npm run format)",
      "Bash(npm run test:*)",
      "Bash(npm run build)",
      "Read(~/.zshrc)",
      "Bash(git diff:*)"
    ],
    "deny": [
      "Bash(curl:*)"
    ],
    "defaultMode": "acceptEdits",
    "disableBypassPermissionsMode": "disable"
  },
  "env": {
    "DISABLE_TELEMETRY": "1",
    "ANTHROPIC_MODEL": "anthropic.claude-sonnet-4-20250514-v1:0",
    "ANTHROPIC_SMALL_FAST_MODEL": "anthropic.claude-sonnet-4-20250514-v1:0",
    "ANTHROPIC_BEDROCK_BASE_URL": "http://mlop-azure-gateway.mediatek.inc",
    "CLAUDE_CODE_USE_BEDROCK": "1",
    "CLAUDE_CODE_SKIP_BEDROCK_AUTH": "1",
    "CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC": "1",
    "ANTHROPIC_CUSTOM_HEADERS": "api-key: <<Your GAISF API KEY>>"
  }
}
```

**Configuration Notes:**
- Replace `<<Your GAISF API KEY>>` with your actual API key
- The `ANTHROPIC_BEDROCK_BASE_URL` is configured to use HTTP instead of HTTPS due to SSL certificate configuration requirements
- **Important**: In some cases, Claude Code may not be able to read the `~/.claude/settings.json` configuration file. If you encounter configuration issues, please refer to the [official settings documentation](https://docs.anthropic.com/en/docs/claude-code/settings#settings-files) for alternative configuration file locations and troubleshooting steps

---

## Usage

After configuration, launch Claude Code:

```bash
claude
```

## Additional Resources

- [Official Documentation](https://docs.anthropic.com/en/docs/claude-code)
- [Settings Documentation](https://docs.anthropic.com/en/docs/claude-code/settings)
- [Sub-Agents](https://docs.anthropic.com/en/docs/claude-code/sub-agents) - Explore the agent feature for specialized tasks
- [MCP Integration](https://docs.anthropic.com/en/docs/claude-code/mcp) - Learn about Model Context Protocol support