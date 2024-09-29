package config

import (
	"log"

	"github.com/gorilla/securecookie"
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/viper"
)

func ReadConfigFile(configFile string) {
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
	setDefaults()

	viper.ReadInConfig()
	viper.AutomaticEnv()
	if viper.ConfigFileUsed() == "" {

		viper.WriteConfigAs("./config.toml")
	}
	log.Println("done reading config", viper.ConfigFileUsed())
}

func setDefaults() {
	// Web
	viper.SetDefault("web.port", 8080)
	viper.SetDefault("web.host", "127.0.0.1")
	viper.SetDefault("web.tls", false)
  viper.SetDefault("web.cert", "/etc/ssl/certs/ssl-cert-snakeoil.pem")
  viper.SetDefault("web.key", "/etc/ssl/key/ssl-cert-snakeoil.key")
	viper.SetDefault("web.sessionpath", "./sessions")
	viper.SetDefault("web.sessionkey", string(securecookie.GenerateRandomKey(32)))
  viper.SetDefault("web.blockfilebrowser", false)
  viper.SetDefault("web.adminpw", string(securecookie.GenerateRandomKey(32)))
	// Containers
  viper.SetDefault("container.image", "lscr.io/linuxserver/firefox")
  viper.SetDefault("container.port", 3000)
  viper.SetDefault("container.isolated", false)
  viper.SetDefault("container.gpu", "")
}

func SaveConfig() error {
	return viper.WriteConfig()
}

