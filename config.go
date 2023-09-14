package main
import (
   "os"
   "fmt"
   "strings"
)

type Config struct {
   Port     string `yaml:"port" env:"MIDDAG_PORT" env-default:"8000"`
   Database struct {
      Host     string `yaml:"host" env:"MIDDAG_DATABASE_HOST" env-default:"0.0.0.0"`
      Port     string `yaml:"port" env:"MIDDAG_DATABASE_PORT" env-default:"5432"`
      Username string `yaml:"username" env:"MIDDAG_DATABASE_USERNAME" env-default:"postgres"`
      Password string `yaml:"password" env:"MIDDAG_DATABASE_PASSWORD"`
      Name     string `yaml:"name" env:"MIDDAG_DATABASE_NAME" env-default:"postgres"`
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

   cfg.Port = env_config("port", "8000")
   cfg.Database.Host = env_config("database_host", "0.0.0.0")
   cfg.Database.Port = env_config("database_port", "5432")
   cfg.Database.Username = env_config("database_username", "postgres")
   cfg.Database.Password = env_config("database_password", "postgres")
   cfg.Database.Name = env_config("database_name", "postgres")

   return cfg
}
