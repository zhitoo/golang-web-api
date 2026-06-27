package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Env  string
	Name string
	Lang string
}
type ServerConfig struct {
	Port    string
	RunMode string
}

type LoggerConfig struct {
	FileName string
	Encoding string
	Level    string
	Logger   string
}

type CorsConfig struct {
	AllowOrigins string
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DB       string
	SslMode  string
	TimeZone string
}

type RedisConfig struct {
	Host               string
	Port               string
	Password           string
	DB                 int
	DialTimeout        time.Duration
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	PoolSize           int
	PoolTimeout        time.Duration
	IdleCheckFrequency time.Duration
	IdleTimeout        time.Duration
}

type JWTConfig struct {
	Secret                     string
	RefreshSecret              string
	AccessTokenExpireDuration  time.Duration
	RefreshTokenExpireDuration time.Duration
}

type NatsConfig struct {
	Host     string
	Port     string
	Password string
}

type Config struct {
	App      AppConfig
	Server   ServerConfig
	Logger   LoggerConfig
	Cors     CorsConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Nats     NatsConfig
}

func getConfigFileName(env string) string {
	if env == "" {
		env = "local"
	}
	return fmt.Sprintf("config-%s", env)
}

func getConfigDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Dir(filename)
}

func LoadConfig(fileName string, fileType string) (*viper.Viper, error) {
	log.Println("start loading config file")
	v := viper.New()
	v.SetConfigType(fileType)
	v.SetConfigName(fileName)

	//all possible path
	v.AddConfigPath(getConfigDir())

	v.AutomaticEnv()

	err := v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("config file not found: %v", err)
			return nil, fmt.Errorf("config file not found: %v", err)
		}
		return nil, err
	}

	return v, nil
}

func ParseConfig(v *viper.Viper) (*Config, error) {
	var config Config
	err := v.Unmarshal(&config)

	if err != nil {
		log.Printf("Unable to parse config: %v", err)
		return nil, err
	}

	return &config, nil

}

func GetConfig() *Config {
	_ = godotenv.Load() // loads .env from current working directory

	cfgPath := getConfigFileName(os.Getenv("APP_ENV"))

	log.Printf("Loading config from: %s", cfgPath)

	v, err := LoadConfig(cfgPath, "yml")
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := ParseConfig(v)
	if err != nil {
		log.Fatal(err)
	}

	return cfg
}
