package util

import (
	"log"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application
// The values are read by viper from config file or environment variables
type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)  // set config path
	viper.SetConfigName("app") // set config name
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env: ", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}
	return

}
