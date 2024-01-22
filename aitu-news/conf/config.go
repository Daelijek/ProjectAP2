package conf

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	"log"
	"sync"
)

type Config struct {
	DB struct {
		Host     string `config: "host"`
		Port     string `config: "port"`
		User     string `config: "user"`
		Password string `config: "password"`
		DbName   string `config: "dbname"`
	} `config: "db"`
	SERVER struct {
		Port          string `config: port`
		StaticDir     string `config: staticDir'`
		SessionSecret string `config: sessionSecret`
	} `config: server`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{}
		config.WithOptions(config.WithTagName("config"))
		config.AddDriver(yaml.Driver)
		err := config.LoadFiles("./config.yml")
		if err != nil {
			log.Println("Cannot read config file")
			panic(err)
		}
		err1 := config.Decode(instance)
		if err1 != nil {
			panic(err1)
		}
	})
	return instance
}
