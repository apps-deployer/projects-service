package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env  string     `yaml:"env" env-default:"local"`
	Db   DbConfig   `yaml:"db"`
	Grpc GrpcConfig `yaml:"grpc"`
	Auth AuthConfig `yaml:"auth"`
}

type DbConfig struct {
	Host string `yaml:"host" env-required:"true"`
	Port int    `yaml:"port" env-default:"5432"`
	User string `yaml:"user" env-required:"true"`
	Pass string `yaml:"pass" env:"DB_PASSWORD" env-required:"true"`
	Name string `yaml:"name" env-required:"true"`
}

type GrpcConfig struct {
	Port    int           `yaml:"port" env-default:"50051"`
	Timeout time.Duration `yaml:"timeout" env-default:"5s"`
}

type AuthConfig struct {
	JwtSecret string `yaml:"jwt_secret" env:"JWT_SECRET" env-required:"true"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}
	return &cfg
}

// fetchConfigPath получает путь к конфигурационному файлу из флага командной строки или из переменной окружения.
// Приоритет: флаг > переменная окружения > значение по умолчанию
// Значение по умолчанию пустая строка
func fetchConfigPath() (res string) {
	flag.StringVar(&res, "config", "", "path/to/config")
	flag.Parse()
	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}
	return
}
