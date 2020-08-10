package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

type Database struct {
	User              string
	Password          string
	Host              string
	DBName            string
	ConnectionTimeout time.Duration
	CommandTimeout    time.Duration
}

type config struct {
	Database Database
	Server   struct {
		Address string
	}
	SWApi struct {
		BaseUrl string
	}
}

var Data config

func ReadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(filepath.Join("$GOPATH", "src", "b2w", "swapi-challenge", "config"))
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}

	if err := viper.Unmarshal(&Data); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
