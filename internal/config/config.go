package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	DB     DB     `yaml:"db"`
	Server Server `yaml:"server"`
}

type DB struct {
	Name            string        `yaml:"name"`
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	Host            string        `yaml:"host"`
	Port            string        `yaml:"port"`
	ConnectAttempts int           `yaml:"connect_attempts"`
	ConnectTimeout  time.Duration `yaml:"connect_timeout"`
	MaxConn         int           `yaml:"max_conn"`
	ConnLifeTime    time.Duration `yaml:"conn_life"`
	ConnIdleTime    time.Duration `yaml:"conn_idle"`
}

type Server struct {
	Host            string        `yaml:"host"`
	Port            string        `yaml:"port"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	IdleTimeout     time.Duration `yaml:"idle_timeout"`
}

func New() (*Config, error) {
	conf := &Config{
		DB: DB{
			Name:     os.Getenv("DB_NAME"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
		},
		Server: Server{
			Host: os.Getenv("SERVER_HOST"),
			Port: os.Getenv("SERVER_PORT"),
		},
	}
	dbConnectAttempts, err := strconv.Atoi(os.Getenv("DB_CONNECT_ATTEMPTS"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse DB_CONNECT_ATTEMPTS: %w", err)
	}
	dbConnectTimeout, err := time.ParseDuration(os.Getenv("DB_CONNECT_TIMEOUT"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse DB_CONNECT_TIMEOUT: %w", err)
	}
	dbMaxConn, err := strconv.Atoi(os.Getenv("DB_MAX_CONN"))
	if err != nil {
		return nil, err
	}
	dbConnLifeTime, err := time.ParseDuration(os.Getenv("DB_CONN_LIFE_TIME"))
	if err != nil {
		return nil, err
	}
	dbConnIdleTime, err := time.ParseDuration(os.Getenv("DB_CONN_IDLE_TIME"))
	if err != nil {
		return nil, err
	}
	conf.DB.ConnectAttempts = dbConnectAttempts
	conf.DB.ConnectTimeout = dbConnectTimeout
	conf.DB.MaxConn = dbMaxConn
	conf.DB.ConnLifeTime = dbConnLifeTime
	conf.DB.ConnIdleTime = dbConnIdleTime

	serverShutdownTimeout, err := time.ParseDuration(os.Getenv("SERVER_SHUTDOWN_TIMEOUT"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse SERVER_SHUTDOWN_TIMEOUT: %w", err)
	}
	serverReadTimeout, err := time.ParseDuration(os.Getenv("SERVER_READ_TIMEOUT"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse SERVER_READ_TIMEOUT: %w", err)
	}
	serverWriteTimeout, err := time.ParseDuration(os.Getenv("SERVER_WRITE_TIMEOUT"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse SERVER_WRITE_TIMEOUT: %w", err)
	}
	serverIdleTimeout, err := time.ParseDuration(os.Getenv("SERVER_IDLE_TIMEOUT"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse SERVER_IDLE_TIMEOUT: %w", err)
	}
	conf.Server.ShutdownTimeout = serverShutdownTimeout
	conf.Server.ReadTimeout = serverReadTimeout
	conf.Server.WriteTimeout = serverWriteTimeout
	conf.Server.IdleTimeout = serverIdleTimeout

	return conf, nil
}
