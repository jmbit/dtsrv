package config

import (
	"log"

	"github.com/jmbit/dtsrv/internal/utils"
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/viper"
)

func ReadConfigFile(configFile string) {
	setDefaults()
	if configFile != "" {
		log.Println("reading Config file", configFile)
		viper.SetConfigFile(configFile)
	} else {
		viper.AddConfigPath(".")
		viper.AddConfigPath("/etc/dtsrv/")
		viper.SetConfigName("config")
		viper.SetConfigType("toml")
	}
	viper.SetEnvPrefix("DTSRV")

	viper.ReadInConfig()
	viper.AutomaticEnv()
	if viper.ConfigFileUsed() == "" {

		viper.WriteConfigAs("./config.toml")
	}
	log.Println("done reading config", viper.ConfigFileUsed())
}

func setDefaults() {
  adminPW, _ := utils.RandomString(32)
  sessionKey, _ := utils.RandomString(32)
	// Web
	viper.SetDefault("web.port", 8080)
	viper.SetDefault("web.host", "127.0.0.1")
	viper.SetDefault("web.tls", false)
  viper.SetDefault("web.cert", "/etc/ssl/certs/ssl-cert-snakeoil.pem")
  viper.SetDefault("web.key", "/etc/ssl/key/ssl-cert-snakeoil.key")
	viper.SetDefault("web.sessionpath", "./sessions")
	viper.SetDefault("web.sessionkey", sessionKey)
  viper.SetDefault("web.blockfilebrowser", false)
  viper.SetDefault("web.adminpw", adminPW)
  viper.SetDefault("web.loghttp", true)
	// Containers
  viper.SetDefault("container.image", "lscr.io/linuxserver/firefox")
  viper.SetDefault("container.port", 3000)
  viper.SetDefault("container.isolated", false)
  viper.SetDefault("container.gpu", "")
  // Default: 12h (in s)
  viper.SetDefault("container.maxage", 43200)
}

func SaveConfig() error {
	return viper.WriteConfig()
}

