package main

import (
	"log"

	"github.com/joho/godotenv"
	internalhttp "github.com/romandnk/HW/hw12_13_14_15_calendar/internal/server/http"
	dbconf "github.com/romandnk/HW/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/spf13/viper"
)

type Config struct {
	Logger  LoggerConf
	Server  internalhttp.ServerConf
	Storage StorageConf
}

func NewConfig(path string) *Config {
	viper.SetConfigFile(path) // find config file with specific path

	err := viper.ReadInConfig() // read config file
	if err != nil {
		log.Fatalf("fatal error config file: %s\n", err)
	}

	if err := godotenv.Load("./configs/.env"); err != nil { // load .env into system
		log.Fatalf("error loading .env: %s", err.Error())
	}

	viper.SetEnvPrefix("calendar") // out env variables will look like CALENDAR_PASSWORD=password
	viper.AutomaticEnv()           // read env variables

	config := Config{
		Logger:  newLoggerConf(),
		Storage: newStorageConf(),
		Server:  newServerConf(),
	}

	return &config
}

type LoggerConf struct {
	Level          string
	Representation string
}

func newLoggerConf() LoggerConf {
	level := viper.GetString("logger.level")
	representation := viper.GetString("logger.representation")
	return LoggerConf{
		Level:          level,
		Representation: representation,
	}
}

func newServerConf() internalhttp.ServerConf {
	host := viper.GetString("server.host")
	port := viper.GetString("server.port")
	readTimeout := viper.GetDuration("server.read_timeout")
	writeTimeout := viper.GetDuration("server.write_timeout")
	return internalhttp.ServerConf{
		Host:         host,
		Port:         port,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
}

type StorageConf struct {
	Memory bool
	DB     dbconf.DBConf
}

func newStorageConf() StorageConf {
	memory := viper.GetBool("storage.memory")
	if memory {
		return StorageConf{
			Memory: true,
		}
	}

	host := viper.GetString("storage.database.host")
	port := viper.GetString("storage.database.port")
	username := viper.GetString("DB_USER")
	password := viper.GetString("DB_PASSWORD")
	dbName := viper.GetString("storage.database.db_name")
	sslmode := viper.GetString("storage.database.sslmode")
	maxConns := viper.GetInt("storage.database.max_conns")
	minConns := viper.GetInt("storage.database.min_conns")
	maxConnLifetime := viper.GetDuration("storage.database.max_conn_lifetime")
	maxConnIdleTime := viper.GetDuration("storage.database.max_conn_idle_time")
	DBconf := dbconf.DBConf{
		Host:            host,
		Port:            port,
		Username:        username,
		Password:        password,
		DBName:          dbName,
		SSLMode:         sslmode,
		MaxConns:        maxConns,
		MinConns:        minConns,
		MaxConnLifetime: maxConnLifetime,
		MaxConnIdleTime: maxConnIdleTime,
	}

	return StorageConf{
		DB: DBconf,
	}
}
