package util

import (
	"github.com/spf13/viper"
	"log"
	"path/filepath"
)

// GetTestConfig : Return test config based on package name
func GetTestConfig(packageName string) map[string]string {
	path, _ := filepath.Abs("../")
	viper.SetConfigName("test_config")
	viper.AddConfigPath(path)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error : %s", err)
	}

	return viper.GetStringMapString(packageName)
}
