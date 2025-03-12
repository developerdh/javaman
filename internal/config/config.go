package config

import (
	"fmt"
	"os"
	"path/filepath"

	"javaman/internal/detect"

	"github.com/spf13/viper"
)

type Config struct {
	Versions map[string]string `mapstructure:"versions"`
	Settings ConfigSettings    `mapstructure:"settings"`
	Aliases  map[string]string `mapstructure:"aliases"`
}

type ConfigSettings struct {
	Default  string `mapstructure:"default"`
	LastUsed string `mapstructure:"last_used"`
}

const (
	configFileName = "config"
	configFileType = "toml"
	configDirName  = ".javaman"
)

var (
	config *Config
)

// Initialize 初始化配置并自动检测JDK
func Initialize() error {
	// 获取用户主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// 创建配置目录
	configDir := filepath.Join(homeDir, configDirName)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// 设置配置文件路径
	configFile := filepath.Join(configDir, configFileName+"."+configFileType)

	// 设置viper配置
	viper.SetConfigName(configFileName)
	viper.SetConfigType(configFileType)
	viper.AddConfigPath(configDir)
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// 如果配置文件不存在，创建新的配置实例
		config = &Config{
			Versions: make(map[string]string),
			Settings: ConfigSettings{},
			Aliases:  make(map[string]string),
		}

		// 先创建空的配置文件
		if err := viper.SafeWriteConfig(); err != nil {
			return fmt.Errorf("failed to create config file: %w", err)
		}

		// 自动检测并添加JDK
		detected, detectErr := detect.DetectJDKs()
		if detectErr != nil {
			return fmt.Errorf("failed to detect JDKs: %w", detectErr)
		}

		// 添加检测到的JDK
		for version, path := range detected {
			config.Versions[version] = path
		}

		// 如果有版本被检测到，设置最新版本为默认版本
		if len(config.Versions) > 0 {
			if latestVer, _, err := detect.GetLatestJDK(); err == nil {
				config.Settings.Default = latestVer
			}
		}

		if err := SaveConfig(); err != nil {
			return fmt.Errorf("failed to save initial config: %w", err)
		}
	} else if err == nil {
		// 如果配置文件存在，读取配置
		if err := viper.ReadInConfig(); err != nil {
			return fmt.Errorf("failed to read config file: %w", err)
		}

		config = &Config{}
		if err := viper.Unmarshal(config); err != nil {
			return fmt.Errorf("failed to unmarshal config: %w", err)
		}

		// 如果配置中没有版本信息，尝试自动检测
		if len(config.Versions) == 0 {
			detected, detectErr := detect.DetectJDKs()
			if detectErr == nil && len(detected) > 0 {
				// 添加检测到的JDK
				for version, path := range detected {
					config.Versions[version] = path
				}

				// 设置最新版本为默认版本
				if latestVer, _, err := detect.GetLatestJDK(); err == nil {
					config.Settings.Default = latestVer
				}

				if err := SaveConfig(); err != nil {
					return fmt.Errorf("failed to save detected JDKs: %w", err)
				}
			}
		}
	} else {
		return fmt.Errorf("failed to check config file: %w", err)
	}

	return nil
}

// GetConfig 获取配置实例
func GetConfig() *Config {
	return config
}

// SaveConfig 保存配置到文件
func SaveConfig() error {
	// 清除现有配置以避免残留
	viper.Set("versions", nil)
	viper.Set("settings", nil)
	viper.Set("aliases", nil)

	// 逐个设置版本路径
	for version, path := range config.Versions {
		viper.Set(fmt.Sprintf("versions.%s", version), path)
	}

	// 设置其他配置
	viper.Set("settings.default", config.Settings.Default)
	viper.Set("settings.last_used", config.Settings.LastUsed)

	// 逐个设置别名
	for alias, version := range config.Aliases {
		viper.Set(fmt.Sprintf("aliases.%s", alias), version)
	}
	// 保存到文件
	return viper.WriteConfig()
}

// AddVersion 添加新的JDK版本
func AddVersion(version, path string) error {
	if config.Versions == nil {
		config.Versions = make(map[string]string)
	}
	config.Versions[version] = path
	return SaveConfig()
}

// RemoveVersion 删除JDK版本
func RemoveVersion(version string) error {
	delete(config.Versions, version)
	// 删除viper中的版本信息
	delete(viper.Get("versions").(map[string]interface{}), version)
	return SaveConfig()
}

// GetVersions 获取所有已配置的JDK版本
func GetVersions() map[string]string {
	return config.Versions
}
