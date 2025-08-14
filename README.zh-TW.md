# Claude Code CLI 使用說明

[English](README.md) | 繁體中文 | [简体中文](README.zh-CN.md)

## 簡介

Claude Code 是 Anthropic 官方的 CLI 工具，提供 AI 編程助手和互動式開發支援。

---

## 取得 API 金鑰

**重要**：使用 Claude Code 前，請先完成以下步驟：

1. 前往 [MediaTek MLOp Gateway](https://mlop-azure-gateway.mediatek.inc/auth/login) 登入
2. 取得您的 GAISF API 金鑰
3. 妥善保存金鑰以供後續使用

**注意**：由於 SSL 憑證設定問題，本說明文件中的網址使用 HTTP 而非 HTTPS，以確保在不同網路環境下的相容性。

---

## 安裝

### 方法一：使用預編譯執行檔（建議）

請前往 [Claude Code Release](https://gitea.mediatek.inc/IT-GAIA/claude-code/releases/latest) 下載最新版本的安裝包。

#### Windows 使用者

⚠️ **重要警告**：Claude Code 對 Windows 的支援有限。如果您無法使用 WSL（Windows 子系統 Linux），我們強烈建議使用 macOS 或 Linux 以獲得最佳體驗。

1. 下載 Windows 版本的 `.exe` 執行檔
2. 為執行檔建立目錄：
   ```cmd
   mkdir %USERPROFILE%\.local\bin
   ```
3. 將下載的執行檔移動到目錄並重新命名：
   ```cmd
   move claude-code-windows-x64.exe %USERPROFILE%\.local\bin\claude.exe
   ```
4. 將目錄加入您的 PATH 環境變數：
   - 開啟系統內容 → 環境變數
   - 將 `%USERPROFILE%\.local\bin` 加入您的 PATH
   - 或使用 PowerShell：
   ```powershell
   $env:PATH += ";$env:USERPROFILE\.local\bin"
   [Environment]::SetEnvironmentVariable("PATH", $env:PATH, "User")
   ```

#### macOS/Linux 使用者
1. 下載適合您平台的執行檔
2. 將執行檔移動至本地 bin 目錄：

```bash
# 設定執行權限
chmod +x claude-code-*

# 移動至 ~/.local/bin（如目錄不存在則建立）
mkdir -p ~/.local/bin
mv claude-code-* ~/.local/bin/claude

# 確保 ~/.local/bin 在您的 PATH 中
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

### 方法二：使用 npm（適合開發者）

使用 npm 安裝 Claude Code（需要 Node.js 環境和程式設計基礎）：

```bash
npm install -g @anthropic-ai/claude-code
```

---

## 平台支援

Claude Code 支援以下平台：
- macOS
- Linux 
- Windows WSL（Windows 子系統 Linux）

### Windows 使用者

對於希望在 WSL 之外使用 Claude Code 的 Windows 使用者，需要安裝 Git Bash：

1. 從以下地址下載並安裝 Git for Windows：https://git-scm.com/downloads/win
2. 設定環境變數指向 Git Bash：
   ```powershell
   $env:CLAUDE_CODE_GIT_BASH_PATH="C:\Program Files\Git\bin\bash.exe"
   ```

更多詳細設定資訊，請參考[官方安裝文件](https://docs.anthropic.com/zh-TW/docs/claude-code/setup)。

---

## 設定檔配置

在 `~/.claude/settings.json` 建立設定檔：

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
    "ANTHROPIC_CUSTOM_HEADERS": "api-key: <<您的 GAISF API 金鑰>>"
  }
}
```

**設定說明：**
- 請將 `<<您的 GAISF API 金鑰>>` 替換為您的實際 API 金鑰
- `ANTHROPIC_BEDROCK_BASE_URL` 設定為使用 HTTP 而非 HTTPS，這是由於 SSL 憑證設定需求
- **重要提醒**：在某些情況下，Claude Code 可能無法讀取 `~/.claude/settings.json` 設定檔。如果您遇到設定問題，請參考 [官方設定檔說明文件](https://docs.anthropic.com/zh-TW/docs/claude-code/settings#%E8%A8%AD%E5%AE%9A%E6%AA%94%E6%A1%88) 了解其他設定檔位置和故障排除步驟

---

## 使用方法

設定完成後，啟動 Claude Code：

```bash
claude
```

---

## 額外資源

- [官方文件](https://docs.anthropic.com/zh-TW/docs/claude-code)
- [設定文件](https://docs.anthropic.com/zh-TW/docs/claude-code/settings)
- [子代理功能](https://docs.anthropic.com/zh-TW/docs/claude-code/sub-agents) - 探索專業任務的代理功能
- [MCP 整合](https://docs.anthropic.com/zh-TW/docs/claude-code/mcp) - 了解模型上下文協定支援
