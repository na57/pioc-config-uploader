package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config 配置文件结构
type Config struct {
	APIKey     string `yaml:"api_key"`
	APISecret  string `yaml:"api_secret"`
	BaseURL    string `yaml:"base_url"`
	ConfigID   string `yaml:"config_id"`
	ConfigFile string `yaml:"config_file"`
}

// APIResponse API响应结构
type APIResponse struct {
	Success bool `json:"success"`
	Data    struct {
		ID            string `json:"id"`
		VersionNumber string `json:"versionNumber"`
	} `json:"data"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

// PIOCClient PIOC API客户端
type PIOCClient struct {
	apiKey    string
	apiSecret string
	baseURL   string
}

// NewPIOCClient 创建新的PIOC客户端
func NewPIOCClient(apiKey, apiSecret, baseURL string) *PIOCClient {
	return &PIOCClient{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		baseURL:   baseURL,
	}
}

// generateSignature 生成HMAC-SHA256签名
func (c *PIOCClient) generateSignature(method, path, queryString, body string) (signature, timestamp string) {
	timestamp = fmt.Sprintf("%d", time.Now().Unix())
	signString := fmt.Sprintf("%s\n%s\n%s\n%s\n%s", method, path, queryString, timestamp, body)

	h := hmac.New(sha256.New, []byte(c.apiSecret))
	h.Write([]byte(signString))
	signature = hex.EncodeToString(h.Sum(nil))

	return signature, timestamp
}

// CreateVersion 创建配置新版本
func (c *PIOCClient) CreateVersion(configID, content, versionNumber string) (*APIResponse, error) {
	path := "/api/external/v1/configsys/versions"
	url := c.baseURL + path

	requestBody := map[string]string{
		"configId":      configID,
		"content":       content,
		"versionNumber": versionNumber,
	}
	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求体失败: %w", err)
	}

	signature, timestamp := c.generateSignature("POST", path, "", string(bodyBytes))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("X-API-Signature", signature)
	req.Header.Set("X-API-Timestamp", timestamp)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &apiResp, nil
}

// LoadConfig 从YAML文件加载配置
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析YAML失败: %w", err)
	}

	// 设置默认值
	if config.BaseURL == "" {
		config.BaseURL = "http://localhost:8080"
	}

	return &config, nil
}

// ReadConfigFile 读取配置文件内容
func ReadConfigFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("读取配置文件失败: %w", err)
	}
	return string(data), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("用法: config-uploader <config.yaml>")
		fmt.Println("\n示例配置文件格式:")
		fmt.Println(`api_key: "pk_your_api_key_here"
api_secret: "your_api_secret_here"
base_url: "http://localhost:8080"
config_id: "your-config-id-here"
config_file: "/path/to/your/config.conf"`)
		os.Exit(1)
	}

	configFile := os.Args[1]

	// 加载YAML配置
	config, err := LoadConfig(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}

	// 验证必要参数
	if config.APIKey == "" || config.APISecret == "" || config.ConfigID == "" || config.ConfigFile == "" {
		fmt.Fprintln(os.Stderr, "错误: api_key, api_secret, config_id, config_file 都是必填项")
		os.Exit(1)
	}

	// 读取配置文件内容
	content, err := ReadConfigFile(config.ConfigFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("正在上传配置文件: %s\n", config.ConfigFile)
	fmt.Printf("配置ID: %s\n", config.ConfigID)
	fmt.Printf("API地址: %s\n", config.BaseURL)

	// 创建PIOC客户端
	client := NewPIOCClient(config.APIKey, config.APISecret, config.BaseURL)

	// 生成版本号：当前日期 YYYYMMDD
	versionNumber := time.Now().Format("20060102")

	// 创建新版本
	resp, err := client.CreateVersion(config.ConfigID, content, versionNumber)
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}

	if resp.Success {
		fmt.Printf("✓ 版本创建成功!\n")
		fmt.Printf("  版本ID: %s\n", resp.Data.ID)
		fmt.Printf("  版本号: %s\n", resp.Data.VersionNumber)
	} else {
		fmt.Fprintf(os.Stderr, "✗ 创建失败: %s\n", resp.Message)
		if resp.Error != "" {
			fmt.Fprintf(os.Stderr, "  错误详情: %s\n", resp.Error)
		}
		os.Exit(1)
	}
}
