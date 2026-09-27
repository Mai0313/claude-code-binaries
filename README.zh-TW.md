# Claude Code Binaries

[English](README.md) | 繁體中文 | [简体中文](README.zh-CN.md)

[Claude Code](https://github.com/anthropics/claude-code) 原生執行檔在 GitHub Releases 上的鏡像，給能連上 GitHub、卻連不上 Anthropic 下載伺服器（`downloads.claude.ai`）的環境使用。

- 檔案與 Anthropic 發布的完全相同，未經修改。每個檔案在發布到這裡之前，都會先比對 Anthropic manifest 裡的 SHA-256。
- 上游推出新版後，一小時內會出現在這裡。
- 這不是 Anthropic 的專案。如果你連得上 `claude.ai`，請改用[官方安裝方式](https://code.claude.com/docs/en/setup)。

鏡像只解決下載的問題。Claude Code 執行時仍需要連到 Anthropic API，或你設定的其他服務，例如 Amazon Bedrock、Google Vertex AI 或 LLM gateway。

## 下載

打開[最新 release](https://github.com/Mai0313/claude-code-binaries/releases/latest)，選擇對應平台的檔案：

| 平台                        | 檔案                                |
| --------------------------- | ----------------------------------- |
| macOS，Apple Silicon        | `claude-<version>-darwin-arm64`     |
| macOS，Intel                | `claude-<version>-darwin-x64`       |
| Linux x64                   | `claude-<version>-linux-x64`        |
| Linux ARM64                 | `claude-<version>-linux-arm64`      |
| Linux x64，musl（Alpine）   | `claude-<version>-linux-x64-musl`   |
| Linux ARM64，musl（Alpine） | `claude-<version>-linux-arm64-musl` |
| Windows x64                 | `claude-<version>-win32-x64.exe`    |
| Windows ARM64               | `claude-<version>-win32-arm64.exe`  |

每個執行檔另有一個較小的 `.zst` 壓縮檔，用 `zstd -d` 解壓即可。`manifest.json`、`manifest.json.sig` 與 `claude-code.asc` 用來[驗證下載的檔案](#%E9%A9%97%E8%AD%89%E4%B8%8B%E8%BC%89%E7%9A%84%E6%AA%94%E6%A1%88)。較舊的 release 可能只有 macOS、glibc Linux 與 Windows x64 的執行檔，沒有 `.zst` 檔，也沒有驗證用的檔案。

## 安裝

### macOS 與 Linux

把 `linux-x64` 換成你的平台：

```bash
chmod +x claude-*-linux-x64
mkdir -p ~/.local/bin
mv claude-*-linux-x64 ~/.local/bin/claude
```

如果之後 shell 找不到 `claude`，在 `~/.bashrc` 或 `~/.zshrc` 裡把 `~/.local/bin` 加進 `PATH`：

```bash
export PATH="$HOME/.local/bin:$PATH"
```

在 Alpine 等 musl 發行版上，還要安裝執行時需要的套件，並在下方的設定檔裡把 `USE_BUILTIN_RIPGREP` 設為 `0`：

```bash
apk add libgcc libstdc++ ripgrep
```

### Windows

在 PowerShell 裡執行，把 `win32-x64` 換成你的平台：

```powershell
New-Item -ItemType Directory -Force "$env:USERPROFILE\.local\bin"
Move-Item claude-*-win32-x64.exe "$env:USERPROFILE\.local\bin\claude.exe"
$userPath = [Environment]::GetEnvironmentVariable("PATH", "User")
[Environment]::SetEnvironmentVariable("PATH", "$userPath;$env:USERPROFILE\.local\bin", "User")
```

開一個新的終端機，`PATH` 的變更才會生效。有安裝 [Git for Windows](https://git-scm.com/downloads/win) 時，Claude Code 會使用 Git Bash，否則使用 PowerShell。

### 關閉自動更新

Claude Code 會到 `downloads.claude.ai` 檢查更新，在需要這個鏡像的環境裡必定失敗。在 `~/.claude/settings.json`（Windows 為 `%USERPROFILE%\.claude\settings.json`）裡關閉：

```json
{
  "env": {
    "DISABLE_AUTOUPDATER": "1"
  }
}
```

要更新時，下載較新的 release 並取代原本的檔案。`claude --version` 會顯示目前安裝的版本。

## 驗證下載的檔案

Anthropic 會簽署 `manifest.json`，而 manifest 列出了每個執行檔的 SHA-256。先檢查簽章：

```bash
gpg --show-keys claude-code.asc
gpg --import claude-code.asc
gpg --verify manifest.json.sig manifest.json
```

`gpg --show-keys` 印出的 fingerprint 必須是 `31DDDE24DDFAB679F42D7BD2BAA929FF1A7ECACE`，與 [Anthropic 官方文件](https://code.claude.com/docs/en/setup#binary-integrity-and-code-signing)公布的一致：從這個鏡像下載的公鑰本身無法證明任何事。`gpg --verify` 應該顯示 `Good signature from "Anthropic Claude Code Release Signing <security@anthropic.com>"`。

接著比對執行檔的 checksum 與 manifest 裡的值：

```bash
sha256sum claude-*-linux-x64 # macOS: shasum -a 256
jq -r '.platforms["linux-x64"].checksum' manifest.json
```

在 Windows 上，`(Get-FileHash claude.exe -Algorithm SHA256).Hash.ToLower()` 會印出 checksum。`.zst` 檔要先解壓，解壓後的檔案會符合 `manifest.json` 裡的 checksum。

## 鏡像的運作方式

`.github/workflows/updater.yml` 每小時執行一次。它向 `scripts/fetch.sh` 取得上游的最新版本，如果還沒有 release 使用這個 tag，就下載並驗證該版本的所有檔案，發布成 release，並以上游 changelog 中該版本的段落作為說明。

`scripts/fetch.sh` 是唯一與 Claude Code 相關的檔案。要用這個 repository 鏡像其他專案，改寫它的三個指令與這份 README 即可：

```bash
./scripts/fetch.sh version                # 上游最新版本
./scripts/fetch.sh download VERSION dist  # VERSION 的所有檔案，已驗證
./scripts/fetch.sh notes VERSION          # VERSION 的 release notes
```

Claude Code 版權屬於 Anthropic PBC，使用時須遵守 [Anthropic 的條款](https://github.com/anthropics/claude-code/blob/main/LICENSE.md)。
