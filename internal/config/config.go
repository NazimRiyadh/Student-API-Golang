package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HttpServer struct {
	Address string `yaml:"address" env:"ADDRESS" env-required:"true"`
}
type Config struct {
	Env         string `yaml:"env" env:"ENV" env-required:"true" env-default:"production"`
	StoragePath string `yaml:"storage_path" env:"STORAGE_PATH" env-required:"true"`
	HttpServer  `yaml:"http_server"`
}

func MustLoad() *Config {
	var ConfigPath string
	ConfigPath = os.Getenv("CONFIG_PATH")

	if ConfigPath == "" {
		//flag.String() returns a value of string pointer
		flags := flag.String("config", "", "path to the configuration file")
		flag.Parse()

		ConfigPath = *flags //that why we use * to access the value in flag pointer
		if ConfigPath == "" {
			log.Fatal("Config file is not set")
		}
	}

	//Use log.Fatal only when: The application cannot start or continue safely,
	// and immediate termination is required, usually only in main()
	if _, err := os.Stat(ConfigPath); os.IsNotExist(err) {
		log.Fatalf("Config file not exists: %s", err)
	}

	var cnf Config

	err := cleanenv.ReadConfig(ConfigPath, &cnf)

	if err != nil {
		log.Fatalf("cannot read config file: %s", err.Error())
	}

	return &cnf
}
