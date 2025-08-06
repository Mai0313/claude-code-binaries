# Claude Code

## 安装

### Linux & macOS

使用单条命令安装 Claude Code：

```bash
curl -fsSLk https://gitea.mediatek.inc/IT-GAIA/claude-code/raw/branch/master/scripts/install.sh | bash
```

### Windows

使用 PowerShell 安装 Claude Code：

```powershell
irm https://gitea.mediatek.inc/IT-GAIA/claude-code/raw/branch/master/scripts/Install.ps1 | iex
```

此命令将自动：
- 下载适合您平台的最新稳定版本
- 使用 SHA256 校验和验证下载文件
- 将可执行文件安装到 `~/.claude/downloads/` (Linux/macOS) 或 `%USERPROFILE%\.claude\downloads\` (Windows)
- 设置可执行权限

## 支持平台

- **Windows**: x64, ARM64
- **macOS**: x64, ARM64 (Apple Silicon)
- **Linux**: x64, ARM64

## 系统要求

- `curl` 命令（必需）
- `jq` 命令（可选，用于更好的 JSON 解析）