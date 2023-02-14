package config

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
    Debug bool `toml:"debug"`
    Server Server `toml:"server"`
}

type Server struct  {
    Port string `toml:"port"`
}

var conf Config

func Load(file string) {
	// Read Config File
	data, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	// Parse Config File
	if _, err := toml.Decode(string(data), &conf); err != nil {
		log.Fatal(err)
	}
}

func Get() *Config {
    return &conf
}
