package main
import (
   "os"
   "fmt"
   "strings"
)

type Config struct {
   port     string `yaml:"port" env:"MIDDAG_PORT" env-default:"8000"`
   database struct {
      host     string `yaml:"host" env:"MIDDAG_DATABASE_HOST" env-default:"0.0.0.0"`
      port     string `yaml:"port" env:"MIDDAG_DATABASE_PORT" env-default:"5432"`
      username string `yaml:"username" env:"MIDDAG_DATABASE_USERNAME" env-default:"postgres"`
      password string `yaml:"password" env:"MIDDAG_DATABASE_PASSWORD"`
      name     string `yaml:"name" env:"MIDDAG_DATABASE_NAME" env-default:"postgres"`
   }
}

func env_config(key string, def string) string {
   env_key := fmt.Sprintf("MIDDAG_%s", strings.ToUpper(key))
   value := os.Getenv(env_key)
   if value == ""{
      return def
   }
   return value
}

func ReadConfig() Config {
   cfg := Config{}

   cfg.port = env_config("port", "8000")
   cfg.database.host = env_config("database_host", "0.0.0.0")
   cfg.database.port = env_config("database_port", "5432")
   cfg.database.username = env_config("database_username", "postgres")
   cfg.database.password = env_config("database_password", "postgres")
   cfg.database.name = env_config("database_name", "postgres")

   return cfg
}
