package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type Config struct {
	Env                        map[string]string       `json:"env"`
	IncludeCoAuthoredBy        bool                    `json:"includeCoAuthoredBy"`
	EnableAllProjectMcpServers bool                    `json:"enableAllProjectMcpServers"`
	Hooks                      map[string][]HookConfig `json:"hooks"`
}

type HookConfig struct {
	Matcher string     `json:"matcher"`
	Hooks   []HookItem `json:"hooks"`
}

type HookItem struct {
	Type    string `json:"type"`
	Command string `json:"command"`
}

func main() {
	fmt.Println("開始安裝 Claude Code CLI 工具...")

	// 1. 檢查作業系統
	osType := runtime.GOOS
	arch := runtime.GOARCH
	fmt.Printf("檢測到作業系統: %s (%s)\n", osType, arch)

	// Windows 特別提示
	if osType == "windows" {
		fmt.Println("⚠️  Windows 用戶注意:")
		fmt.Println("   - 建議以管理員身份執行此程式")
		fmt.Println("   - 程式將自動安裝 Chocolatey 和 Node.js")
		fmt.Println("   - 如果安裝失敗，可能需要重新啟動命令提示字元")
		fmt.Println()
	}

	// 2. 安裝 Node.js
	if err := installNodeJS(osType); err != nil {
		fmt.Printf("安裝 Node.js 失敗: %v\n", err)
		return
	}

	// 3. 驗證 Node.js 和 npm 安裝
	if err := verifyNodeInstallation(); err != nil {
		fmt.Printf("驗證 Node.js 安裝失敗: %v\n", err)
		return
	}

	// 4. 安裝 Claude Code 套件
	if err := installClaudeCode(); err != nil {
		fmt.Printf("安裝 Claude Code 套件失敗: %v\n", err)
		return
	}

	// 5. 下載並安裝 claude_analysis
	if err := installClaudeAnalysis(osType, arch); err != nil {
		fmt.Printf("安裝 claude_analysis 失敗: %v\n", err)
		return
	}

	// 6. 創建配置文件
	if err := createConfig(osType, arch); err != nil {
		fmt.Printf("創建配置文件失敗: %v\n", err)
		return
	}

	// 7. 驗證配置文件
	if err := verifyConfig(); err != nil {
		fmt.Printf("驗證配置文件失敗: %v\n", err)
		return
	}

	fmt.Println("✅ Claude Code CLI 工具安裝完成！")
}

func installNodeJS(osType string) error {
	fmt.Println("正在安裝 Node.js v22.18.0 LTS...")

	var cmd *exec.Cmd

	switch osType {
	case "windows":
		// Windows: 先安裝 Chocolatey，然後安裝 Node.js
		if err := installChocolatey(); err != nil {
			fmt.Printf("安裝 Chocolatey 失敗: %v\n", err)
			fmt.Println("請手動從 https://nodejs.org 下載安裝 Node.js v22.18.0")
			return fmt.Errorf("請手動安裝 Node.js")
		}

		fmt.Println("正在使用 Chocolatey 安裝 Node.js...")
		// 使用 Chocolatey 安裝 Node.js (安裝最新的 LTS 版本)
		cmd = exec.Command("choco", "install", "nodejs", "-y")
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Chocolatey 輸出: %s\n", string(output))
			fmt.Println("❌ Chocolatey 安裝 Node.js 失敗")
			fmt.Println("請手動從 https://nodejs.org 下載安裝 Node.js v22.18.0")
			return fmt.Errorf("請手動安裝 Node.js")
		}
		fmt.Println("✅ Node.js 安裝完成")
	case "darwin":
		// macOS: 使用 homebrew
		cmd = exec.Command("brew", "install", "node@22")
		if err := cmd.Run(); err != nil {
			fmt.Println("Homebrew 安裝失敗，請手動從 https://nodejs.org 下載安裝 Node.js v22.18.0")
			return fmt.Errorf("請手動安裝 Node.js")
		}
	case "linux":
		// Linux: 使用 NodeSource repository
		commands := [][]string{
			{"curl", "-fsSL", "https://deb.nodesource.com/setup_22.x", "-o", "/tmp/nodesource_setup.sh"},
			{"sudo", "bash", "/tmp/nodesource_setup.sh"},
			{"sudo", "apt-get", "install", "-y", "nodejs"},
		}

		for _, cmdArgs := range commands {
			cmd = exec.Command(cmdArgs[0], cmdArgs[1:]...)
			if err := cmd.Run(); err != nil {
				fmt.Printf("執行命令失敗: %v\n", strings.Join(cmdArgs, " "))
				return err
			}
		}
	default:
		return fmt.Errorf("不支援的作業系統: %s", osType)
	}

	return nil
}

func installChocolatey() error {
	fmt.Println("檢查 Chocolatey 是否已安裝...")

	// 先檢查 Chocolatey 是否已經安裝
	checkCmd := exec.Command("choco", "--version")
	if err := checkCmd.Run(); err == nil {
		fmt.Println("✅ Chocolatey 已安裝")
		return nil
	}

	fmt.Println("正在安裝 Chocolatey...")
	fmt.Println("⚠️  注意: 此操作需要管理員權限")

	// 使用 PowerShell 安裝 Chocolatey (更完整的安裝命令)
	installScript := `Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))`

	cmd := exec.Command("powershell", "-Command", installScript)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("PowerShell 輸出: %s\n", string(output))
		fmt.Println("\n❌ 自動安裝 Chocolatey 失敗")
		fmt.Println("請以管理員身份執行以下命令手動安裝:")
		fmt.Println(`powershell -c "irm https://community.chocolatey.org/install.ps1|iex"`)
		return fmt.Errorf("安裝 Chocolatey 失敗: %v", err)
	}

	fmt.Println("✅ Chocolatey 安裝完成")

	// 刷新環境變數
	refreshCmd := exec.Command("powershell", "-Command", "refreshenv")
	refreshCmd.Run()

	// 再次檢查安裝是否成功
	verifyCmd := exec.Command("choco", "--version")
	if err := verifyCmd.Run(); err != nil {
		fmt.Println("⚠️  Chocolatey 可能需要重新啟動命令提示字元才能生效")
		fmt.Println("如果繼續安裝失敗，請:")
		fmt.Println("1. 重新啟動命令提示字元")
		fmt.Println("2. 以管理員身份重新執行此程式")
		return fmt.Errorf("Chocolatey 安裝驗證失敗，請重新啟動命令提示字元後再試")
	}

	return nil
}

func verifyNodeInstallation() error {
	fmt.Println("驗證 Node.js 和 npm 安裝...")

	// 檢查 node 版本
	nodeCmd := exec.Command("node", "-v")
	nodeOutput, err := nodeCmd.Output()
	if err != nil {
		return fmt.Errorf("無法執行 node -v: %v", err)
	}
	fmt.Printf("Node.js 版本: %s", string(nodeOutput))

	// 檢查 npm 版本
	npmCmd := exec.Command("npm", "-v")
	npmOutput, err := npmCmd.Output()
	if err != nil {
		return fmt.Errorf("無法執行 npm -v: %v", err)
	}
	fmt.Printf("npm 版本: %s", string(npmOutput))

	return nil
}

func installClaudeCode() error {
	fmt.Println("正在安裝 @anthropic-ai/claude-code 套件...")

	cmd := exec.Command("npm", "install", "-g", "@anthropic-ai/claude-code", "--registry=http://oa-mirror.mediatek.inc/repository/npm")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("npm install 失敗: %v", err)
	}

	// 驗證安裝
	verifyCmd := exec.Command("claude", "--version")
	output, err := verifyCmd.Output()
	if err != nil {
		return fmt.Errorf("無法執行 claude --version: %v", err)
	}
	fmt.Printf("Claude Code 版本: %s", string(output))

	return nil
}

func installClaudeAnalysis(osType, arch string) error {
	fmt.Println("正在下載並安裝 claude_analysis...")

	// 確定下載 URL 和文件名
	var url, filename string

	switch osType {
	case "darwin":
		if arch == "arm64" {
			url = "https://gitea.mediatek.inc/IT-GAIA/claude-code-monitor/releases/download/v0.0.2/claude_analysis-darwin-arm64"
			filename = "claude_analysis-darwin-arm64"
		} else {
			url = "https://gitea.mediatek.inc/IT-GAIA/claude-code-monitor/releases/download/v0.0.2/claude_analysis-darwin-amd64"
			filename = "claude_analysis-darwin-amd64"
		}
	case "linux":
		if arch == "arm64" {
			url = "https://gitea.mediatek.inc/IT-GAIA/claude-code-monitor/releases/download/v0.0.2/claude_analysis-linux-arm64"
			filename = "claude_analysis-linux-arm64"
		} else {
			url = "https://gitea.mediatek.inc/IT-GAIA/claude-code-monitor/releases/download/v0.0.2/claude_analysis-linux-amd64"
			filename = "claude_analysis-linux-amd64"
		}
	case "windows":
		url = "https://gitea.mediatek.inc/IT-GAIA/claude-code-monitor/releases/download/v0.0.2/claude_analysis-windows-amd64.exe"
		filename = "claude_analysis-windows-amd64.exe"
	default:
		return fmt.Errorf("不支援的作業系統: %s", osType)
	}

	// 創建 ~/.claude 目錄
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("無法取得用戶主目錄: %v", err)
	}

	claudeDir := filepath.Join(homeDir, ".claude")
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		return fmt.Errorf("無法創建 .claude 目錄: %v", err)
	}

	// 下載文件
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("下載失敗: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下載失敗，HTTP 狀態碼: %d", resp.StatusCode)
	}

	// 保存文件
	filePath := filepath.Join(claudeDir, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("無法創建文件: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("無法寫入文件: %v", err)
	}

	// 設置執行權限 (Unix 系統)
	if osType != "windows" {
		if err := os.Chmod(filePath, 0755); err != nil {
			return fmt.Errorf("無法設置執行權限: %v", err)
		}
	}

	fmt.Printf("claude_analysis 已安裝到: %s\n", filePath)
	return nil
}

func createConfig(osType, arch string) error {
	fmt.Println("正在創建配置文件...")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("無法取得用戶主目錄: %v", err)
	}

	// 確定 claude_analysis 的文件名
	var analysisFilename string
	switch osType {
	case "darwin":
		if arch == "arm64" {
			analysisFilename = "claude_analysis-darwin-arm64"
		} else {
			analysisFilename = "claude_analysis-darwin-amd64"
		}
	case "linux":
		if arch == "arm64" {
			analysisFilename = "claude_analysis-linux-arm64"
		} else {
			analysisFilename = "claude_analysis-linux-amd64"
		}
	case "windows":
		analysisFilename = "claude_analysis-windows-amd64.exe"
	}

	commandPath := filepath.Join("~", ".claude", analysisFilename)

	config := Config{
		Env: map[string]string{
			"DISABLE_TELEMETRY":                        "1",
			"CLAUDE_CODE_USE_BEDROCK":                  "1",
			"ANTHROPIC_BEDROCK_BASE_URL":               "http://mlop-azure-rddmz.mediatek.inc",
			"CLAUDE_CODE_ENABLE_TELEMETRY":             "1",
			"CLAUDE_CODE_SKIP_BEDROCK_AUTH":            "1",
			"CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC": "1",
		},
		IncludeCoAuthoredBy:        true,
		EnableAllProjectMcpServers: true,
		Hooks: map[string][]HookConfig{
			"Stop": {
				{
					Matcher: "*",
					Hooks: []HookItem{
						{
							Type:    "command",
							Command: commandPath,
						},
					},
				},
			},
		},
	}

	configJSON, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("無法序列化配置: %v", err)
	}

	configPath := filepath.Join(homeDir, ".claude", "settings.json")
	if err := os.WriteFile(configPath, configJSON, 0644); err != nil {
		return fmt.Errorf("無法寫入配置文件: %v", err)
	}

	fmt.Printf("配置文件已創建: %s\n", configPath)
	return nil
}

func verifyConfig() error {
	fmt.Println("驗證配置文件...")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("無法取得用戶主目錄: %v", err)
	}

	configPath := filepath.Join(homeDir, ".claude", "settings.json")

	// 檢查文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("配置文件不存在: %s", configPath)
	}

	// 讀取並驗證 JSON 格式
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("無法讀取配置文件: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("配置文件格式錯誤: %v", err)
	}

	fmt.Println("✅ 配置文件驗證成功")
	return nil
}
