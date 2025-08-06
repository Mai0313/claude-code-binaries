# Claude Code

## 安裝

### Linux & macOS

使用單一指令安裝 Claude Code：

```bash
curl -fsSLk https://gitea.mediatek.inc/IT-GAIA/claude-code/raw/branch/master/scripts/install.sh | bash
```

### Windows

使用 PowerShell 安裝 Claude Code：

```powershell
irm https://gitea.mediatek.inc/IT-GAIA/claude-code/raw/branch/master/scripts/Install.ps1 | iex
```

此指令會自動：
- 下載適合您平台的最新穩定版本
- 使用 SHA256 校驗和驗證下載檔案
- 將執行檔安裝到 `~/.claude/downloads/` (Linux/macOS) 或 `%USERPROFILE%\.claude\downloads\` (Windows)
- 設定執行檔權限

## 支援平台

- **Windows**: x64, ARM64
- **macOS**: x64, ARM64 (Apple Silicon)
- **Linux**: x64, ARM64

## 系統需求

- `curl` 指令（必需）
- `jq` 指令（選用，用於更好的 JSON 解析）