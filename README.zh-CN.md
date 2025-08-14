# Claude Code CLI 使用指南

[English](README.md) | [繁體中文](README.zh-TW.md) | 简体中文

## 简介

Claude Code 是 Anthropic 官方的 CLI 工具，提供 AI 编程助手和交互式开发支持。

---

## 获取 API 密钥

**重要**：使用 Claude Code 前，请先完成以下步骤：

1. 访问 [MediaTek MLOp Gateway](https://mlop-azure-gateway.mediatek.inc/auth/login) 登录
2. 获取您的 GAISF API 密钥
3. 妥善保存密钥以供后续使用

**注意**：由于 SSL 证书配置问题，本说明文档中的网址使用 HTTP 而非 HTTPS，以确保在不同网络环境下的兼容性。

---

## 安装

### 方法一：使用预编译可执行文件（推荐）

请访问 [Claude Code Release](https://gitea.mediatek.inc/IT-GAIA/claude-code/releases/latest) 下载最新版本的安装包。

#### Windows 用户

⚠️ **重要警告**：Claude Code 对 Windows 的支持有限。如果您无法使用 WSL（Windows 子系统 Linux），我们强烈建议使用 macOS 或 Linux 以获得最佳体验。

1. 下载 Windows 版本的 `.exe` 可执行文件
2. 为可执行文件创建目录：
   ```cmd
   mkdir %USERPROFILE%\.local\bin
   ```
3. 将下载的可执行文件移动到目录并重命名：
   ```cmd
   move claude-code-windows-x64.exe %USERPROFILE%\.local\bin\claude.exe
   ```
4. 将目录添加到您的 PATH 环境变量：
   - 打开系统属性 → 环境变量
   - 将 `%USERPROFILE%\.local\bin` 添加到您的 PATH
   - 或使用 PowerShell：
   ```powershell
   $env:PATH += ";$env:USERPROFILE\.local\bin"
   [Environment]::SetEnvironmentVariable("PATH", $env:PATH, "User")
   ```

#### macOS/Linux 用户
1. 下载适合您平台的可执行文件
2. 将可执行文件移动到本地 bin 目录：

```bash
# 设置执行权限
chmod +x claude-code-*

# 移动到 ~/.local/bin（如目录不存在则创建）
mkdir -p ~/.local/bin
mv claude-code-* ~/.local/bin/claude

# 确保 ~/.local/bin 在您的 PATH 中
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

### 方法二：使用 npm（适合开发者）

使用 npm 安装 Claude Code（需要 Node.js 环境和编程基础）：

```bash
npm install -g @anthropic-ai/claude-code
```

---

## 平台支持

Claude Code 支持以下平台：
- macOS
- Linux 
- Windows WSL（Windows 子系统 Linux）

### Windows 用户

对于希望在 WSL 之外使用 Claude Code 的 Windows 用户，需要安装 Git Bash：

1. 从以下地址下载并安装 Git for Windows：https://git-scm.com/downloads/win
2. 设置环境变量指向 Git Bash：
   ```powershell
   $env:CLAUDE_CODE_GIT_BASH_PATH="C:\Program Files\Git\bin\bash.exe"
   ```

更多详细设置信息，请参考[官方安装文档](https://docs.anthropic.com/zh-CN/docs/claude-code/setup)。

---

## 配置文件设置

在 `~/.claude/settings.json` 创建配置文件：

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
    "ANTHROPIC_CUSTOM_HEADERS": "api-key: <<您的 GAISF API 密钥>>"
  }
}
```

**配置说明：**
- 请将 `<<您的 GAISF API 密钥>>` 替换为您的实际 API 密钥
- `ANTHROPIC_BEDROCK_BASE_URL` 配置为使用 HTTP 而非 HTTPS，这是由于 SSL 证书配置需求
- **重要提醒**：在某些情况下，Claude Code 可能无法读取 `~/.claude/settings.json` 配置文件。如果您遇到配置问题，请参考 [官方设置文档](https://docs.anthropic.com/zh-CN/docs/claude-code/settings#%E8%AE%BE%E7%BD%AE%E6%96%87%E4%BB%B6) 了解其他配置文件位置和故障排除步骤

---

## 使用方法

配置完成后，启动 Claude Code：

```bash
claude
```

---

## 其他资源

- [官方文档](https://docs.anthropic.com/zh-CN/docs/claude-code)
- [设置文档](https://docs.anthropic.com/zh-CN/docs/claude-code/settings)
- [子代理功能](https://docs.anthropic.com/zh-CN/docs/claude-code/sub-agents) - 探索专业任务的代理功能
- [MCP 集成](https://docs.anthropic.com/zh-CN/docs/claude-code/mcp) - 了解模型上下文协议支持