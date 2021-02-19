package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Parth576/ytplay/utils"
	"github.com/spf13/viper"
)

// InitConfig creates the config file if it does not exist
func InitConfig(homedir string) {
	configPath := filepath.Join(homedir, ".ytplay.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		file, err := os.Create(configPath)
		utils.PrintErr(err)
		defer file.Close()
	}
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Could not find config file at " + configPath)
		}
	}
	viper.WriteConfig()
}

// CheckAPIKey checks if YT API key exists
func CheckAPIKey(apiKey, keyFlag string) {
	if apiKey == "" && keyFlag == "" {
		fmt.Println("Youtube API key not set, please generate API key from https://console.developers.google.com")
		fmt.Println("And then run the command:")
		fmt.Println("ytplay -key <your-api-key>")
		os.Exit(1)
	}

	if keyFlag != "" {
		viper.Set("YOUTUBE_API_KEY", keyFlag)
		viper.WriteConfig()
		apiKey = keyFlag
		fmt.Println("Youtube API key saved!")
		os.Exit(1)
	}

}
