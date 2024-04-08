package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string     `yaml:"env" env-default:"local"`
	StoragePath string     `yaml:"storage_path" env-required:"true"`
	HTTPServer  HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Port        int           `yaml:"port"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}

	conf := new(Config)
	if err := cleanenv.ReadConfig(path, conf); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return conf
}

func fetchConfigPath() (path string) {
	flag.StringVar(&path, "config", "", "path to config file")
	flag.Parse()

	if path == "" {
		path = os.Getenv("CONFIG_FILE")

	}
	return
}
