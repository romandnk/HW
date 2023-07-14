package main

import (
	"log"

	"github.com/joho/godotenv"
	internalhttp "github.com/romandnk/HW/hw12_13_14_15_calendar/internal/server/http"
	dbconf "github.com/romandnk/HW/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/spf13/viper"
)

type Config struct {
	Logger LoggerConf
	DB     dbconf.DBConf
	Server internalhttp.ServerConf
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
		Logger: NewLoggerConf(),
		DB:     NewDBConf(),
		Server: NewServerConf(),
	}

	return &config
}

type LoggerConf struct {
	Level string
}

func NewLoggerConf() LoggerConf {
	level := viper.GetString("logger.level")
	return LoggerConf{
		Level: level,
	}
}

func NewServerConf() internalhttp.ServerConf {
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

func NewDBConf() dbconf.DBConf {
	host := viper.GetString("database.host")
	port := viper.GetString("database.port")
	username := viper.GetString("DB_USER")
	password := viper.GetString("DB_PASSWORD")
	dbName := viper.GetString("database.db_name")
	sslmode := viper.GetString("database.sslmode")
	maxConns := viper.GetInt32("database.max_conns")
	minConns := viper.GetInt32("database.min_conns")
	maxConnLifetime := viper.GetDuration("database.max_conn_lifetime")
	maxConnIdleTime := viper.GetDuration("database.max_conn_idle_time")
	return dbconf.DBConf{
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
}
