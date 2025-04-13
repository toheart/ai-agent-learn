package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config 存储应用程序的配置信息
type Config struct {
	DatabaseURL  string `json:"database_url"`
	ServerPort   int    `json:"server_port"`
	LogLevel     string `json:"log_level"`
	MaxConnections int  `json:"max_connections"`
	Debug        bool   `json:"debug"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		DatabaseURL:    "postgres://user:password@localhost:5432/myapp",
		ServerPort:     8080,
		LogLevel:       "info",
		MaxConnections: 100,
		Debug:          false,
	}
}

// LoadConfig 从指定路径加载配置文件
func LoadConfig(path string) (*Config, error) {
	// 首先检查文件是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", path)
	}

	// 读取文件内容
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	// 解析JSON
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return &config, nil
}

// SaveConfig 将配置保存到指定路径
func SaveConfig(config *Config, path string) error {
	// 确保目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// 序列化为JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize config: %v", err)
	}

	// 写入文件
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

// Initialize 初始化应用程序的配置
func Initialize() (*Config, error) {
	configPath := "config.json"
	
	// 尝试加载配置
	config, err := LoadConfig(configPath)
	if err != nil {
		// 如果配置不存在，创建默认配置
		config = DefaultConfig()
		if err := SaveConfig(config, configPath); err != nil {
			return nil, fmt.Errorf("failed to save default config: %v", err)
		}
		fmt.Println("Created default configuration file")
	}
	
	return config, nil
} 