# Claude Code Binaries

English | [繁體中文](README.zh-TW.md) | [简体中文](README.zh-CN.md)

A mirror of the [Claude Code](https://github.com/anthropics/claude-code) native binaries on GitHub Releases, for machines that can reach GitHub but not Anthropic's download server (`downloads.claude.ai`).

- The files are the ones Anthropic publishes, unmodified. Each one is checked against the SHA-256 in Anthropic's manifest before it is published here.
- A new upstream version appears here within an hour.
- This is not an Anthropic project. If you can reach `claude.ai`, use the [official installer](https://code.claude.com/docs/en/setup) instead.

The mirror only solves the download. Claude Code still needs to reach the Anthropic API, or whichever provider you configure it for, such as Amazon Bedrock, Google Vertex AI or an LLM gateway.

## Download

Open the [latest release](https://github.com/Mai0313/claude-code-binaries/releases/latest) and pick the file for your platform:

| Platform                   | File                                |
| -------------------------- | ----------------------------------- |
| macOS, Apple Silicon       | `claude-<version>-darwin-arm64`     |
| macOS, Intel               | `claude-<version>-darwin-x64`       |
| Linux x64                  | `claude-<version>-linux-x64`        |
| Linux ARM64                | `claude-<version>-linux-arm64`      |
| Linux x64, musl (Alpine)   | `claude-<version>-linux-x64-musl`   |
| Linux ARM64, musl (Alpine) | `claude-<version>-linux-arm64-musl` |
| Windows x64                | `claude-<version>-win32-x64.exe`    |
| Windows ARM64              | `claude-<version>-win32-arm64.exe`  |

Every binary also comes as a smaller `.zst` file; unpack it with `zstd -d`. `manifest.json`, `manifest.json.sig` and `claude-code.asc` are there to [verify a download](#verify-a-download). Older releases may carry only the macOS, glibc Linux and Windows x64 binaries, without the `.zst` files or the files to verify them.

## Install

### macOS and Linux

Replace `linux-x64` with your platform:

```bash
chmod +x claude-*-linux-x64
mkdir -p ~/.local/bin
mv claude-*-linux-x64 ~/.local/bin/claude
```

If the shell cannot find `claude` afterwards, add `~/.local/bin` to your `PATH` in `~/.bashrc` or `~/.zshrc`:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

On Alpine and other musl distributions, also install the runtime dependencies and set `USE_BUILTIN_RIPGREP` to `0` in the settings file below:

```bash
apk add libgcc libstdc++ ripgrep
```

### Windows

In PowerShell, replace `win32-x64` with your platform:

```powershell
New-Item -ItemType Directory -Force "$env:USERPROFILE\.local\bin"
Move-Item claude-*-win32-x64.exe "$env:USERPROFILE\.local\bin\claude.exe"
$userPath = [Environment]::GetEnvironmentVariable("PATH", "User")
[Environment]::SetEnvironmentVariable("PATH", "$userPath;$env:USERPROFILE\.local\bin", "User")
```

Open a new terminal for the `PATH` change to apply. Claude Code uses Git Bash when [Git for Windows](https://git-scm.com/downloads/win) is installed and PowerShell otherwise.

### Turn off auto-update

Claude Code checks `downloads.claude.ai` for updates, which fails wherever this mirror is needed. Turn the check off in `~/.claude/settings.json` (`%USERPROFILE%\.claude\settings.json` on Windows):

```json
{
  "env": {
    "DISABLE_AUTOUPDATER": "1"
  }
}
```

To update, download a newer release and replace the file. `claude --version` shows the installed version.

## Verify a download

Anthropic signs `manifest.json`, and the manifest lists the SHA-256 of every binary. Check the signature first:

```bash
gpg --show-keys claude-code.asc
gpg --import claude-code.asc
gpg --verify manifest.json.sig manifest.json
```

The fingerprint `gpg --show-keys` prints must be `31DDDE24DDFAB679F42D7BD2BAA929FF1A7ECACE`, as published in [Anthropic's documentation](https://code.claude.com/docs/en/setup#binary-integrity-and-code-signing): a key downloaded from this mirror proves nothing on its own. `gpg --verify` should report `Good signature from "Anthropic Claude Code Release Signing <security@anthropic.com>"`.

Then compare the checksum of the binary with the one in the manifest:

```bash
sha256sum claude-*-linux-x64 # macOS: shasum -a 256
jq -r '.platforms["linux-x64"].checksum' manifest.json
```

On Windows, `(Get-FileHash claude.exe -Algorithm SHA256).Hash.ToLower()` prints the checksum. For a `.zst` file, unpack it first; the result matches the checksum in `manifest.json`.

## Development

How the mirror publishes a release, and how to reuse it for another project, is in [CONTRIBUTING.md](./.github/CONTRIBUTING.md).

Claude Code is © Anthropic PBC, and its use is subject to [Anthropic's terms](https://github.com/anthropics/claude-code/blob/main/LICENSE.md).
