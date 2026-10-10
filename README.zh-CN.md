# Claude Code Binaries

[English](README.md) | [繁體中文](README.zh-TW.md) | 简体中文

[Claude Code](https://github.com/anthropics/claude-code) 原生可执行文件在 GitHub Releases 上的镜像，供能访问 GitHub、却无法访问 Anthropic 下载服务器（`downloads.claude.ai`）的环境使用。

- 文件与 Anthropic 发布的完全相同，未经修改。每个文件在发布到这里之前，都会先比对 Anthropic manifest 中的 SHA-256。
- 上游发布新版本后，一小时内会出现在这里。
- 这不是 Anthropic 的项目。如果你能访问 `claude.ai`，请改用[官方安装方式](https://code.claude.com/docs/en/setup)。

镜像只解决下载的问题。Claude Code 运行时仍需要连接 Anthropic API，或你配置的其他服务，例如 Amazon Bedrock、Google Vertex AI 或 LLM gateway。

## 下载

打开[最新 release](https://github.com/Mai0313/claude-code-binaries/releases/latest)，选择对应平台的文件：

| 平台                        | 文件                                |
| --------------------------- | ----------------------------------- |
| macOS，Apple Silicon        | `claude-<version>-darwin-arm64`     |
| macOS，Intel                | `claude-<version>-darwin-x64`       |
| Linux x64                   | `claude-<version>-linux-x64`        |
| Linux ARM64                 | `claude-<version>-linux-arm64`      |
| Linux x64，musl（Alpine）   | `claude-<version>-linux-x64-musl`   |
| Linux ARM64，musl（Alpine） | `claude-<version>-linux-arm64-musl` |
| Windows x64                 | `claude-<version>-win32-x64.exe`    |
| Windows ARM64               | `claude-<version>-win32-arm64.exe`  |

每个可执行文件另有一个较小的 `.zst` 压缩文件，用 `zstd -d` 解压即可。`manifest.json`、`manifest.json.sig` 和 `claude-code.asc` 用于[验证下载的文件](#%E9%AA%8C%E8%AF%81%E4%B8%8B%E8%BD%BD%E7%9A%84%E6%96%87%E4%BB%B6)。较旧的 release 可能只有 macOS、glibc Linux 和 Windows x64 的可执行文件，没有 `.zst` 文件，也没有验证用的文件。

## 安装

### macOS 和 Linux

把 `linux-x64` 换成你的平台：

```bash
chmod +x claude-*-linux-x64
mkdir -p ~/.local/bin
mv claude-*-linux-x64 ~/.local/bin/claude
```

如果之后 shell 找不到 `claude`，在 `~/.bashrc` 或 `~/.zshrc` 中把 `~/.local/bin` 加入 `PATH`：

```bash
export PATH="$HOME/.local/bin:$PATH"
```

在 Alpine 等 musl 发行版上，还需要安装运行时依赖，并在下方的配置文件中把 `USE_BUILTIN_RIPGREP` 设为 `0`：

```bash
apk add libgcc libstdc++ ripgrep
```

### Windows

在 PowerShell 中运行，把 `win32-x64` 换成你的平台：

```powershell
New-Item -ItemType Directory -Force "$env:USERPROFILE\.local\bin"
Move-Item claude-*-win32-x64.exe "$env:USERPROFILE\.local\bin\claude.exe"
$userPath = [Environment]::GetEnvironmentVariable("PATH", "User")
[Environment]::SetEnvironmentVariable("PATH", "$userPath;$env:USERPROFILE\.local\bin", "User")
```

打开一个新的终端，`PATH` 的修改才会生效。安装了 [Git for Windows](https://git-scm.com/downloads/win) 时，Claude Code 会使用 Git Bash，否则使用 PowerShell。

### 关闭自动更新

Claude Code 会到 `downloads.claude.ai` 检查更新，在需要这个镜像的环境中必然失败。在 `~/.claude/settings.json`（Windows 为 `%USERPROFILE%\.claude\settings.json`）中关闭：

```json
{
  "env": {
    "DISABLE_AUTOUPDATER": "1"
  }
}
```

需要更新时，下载较新的 release 并替换原来的文件。`claude --version` 会显示当前安装的版本。

## 验证下载的文件

Anthropic 会签署 `manifest.json`，而 manifest 列出了每个可执行文件的 SHA-256。先检查签名：

```bash
gpg --show-keys claude-code.asc
gpg --import claude-code.asc
gpg --verify manifest.json.sig manifest.json
```

`gpg --show-keys` 输出的 fingerprint 必须是 `31DDDE24DDFAB679F42D7BD2BAA929FF1A7ECACE`，与 [Anthropic 官方文档](https://code.claude.com/docs/en/setup#binary-integrity-and-code-signing)公布的一致：从这个镜像下载的公钥本身无法证明任何事。`gpg --verify` 应该显示 `Good signature from "Anthropic Claude Code Release Signing <security@anthropic.com>"`。

接着比对可执行文件的 checksum 与 manifest 中的值：

```bash
sha256sum claude-*-linux-x64 # macOS: shasum -a 256
jq -r '.platforms["linux-x64"].checksum' manifest.json
```

在 Windows 上，`(Get-FileHash claude.exe -Algorithm SHA256).Hash.ToLower()` 会输出 checksum。`.zst` 文件要先解压，解压后的文件会符合 `manifest.json` 中的 checksum。

## 开发

镜像如何发布 release，以及如何改用于镜像其他项目，请见 [CONTRIBUTING.md](./.github/CONTRIBUTING.md)。

Claude Code 版权归 Anthropic PBC 所有，使用时须遵守 [Anthropic 的条款](https://github.com/anthropics/claude-code/blob/main/LICENSE.md)。
