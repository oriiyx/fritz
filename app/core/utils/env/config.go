package env

import (
	"log"
	"time"

	"github.com/joeshaw/envdecode"
)

type Conf struct {
	DB           ConfDB
	Auth         ConfAuth
	Server       ConfServer
	IsProduction bool
}

type ConfServer struct {
	ENV    string `env:"APP_ENV,required"`
	Port   int    `env:"SERVER_PORT,required"`
	Debug  bool   `env:"SERVER_DEBUG,required"`
	Secret []byte `env:"SECRET_KEY,required"`

	CORSMaxAge           *int     `env:"CORS_MAX_AGE"`
	CORSOrigins          []string `env:"CORS_ORIGINS"`
	CORSMethods          []string `env:"CORS_METHODS"`
	CORSHeaders          []string `env:"CORS_HEADERS"`
	CORSAllowCredentials *bool    `env:"CORS_ALLOW_CREDENTIALS"`
	CORSExposedHeaders   []string `env:"CORS_EXPOSED_HEADERS"`

	TimeoutIdle  time.Duration `env:"SERVER_TIMEOUT_IDLE,required"`
	TimeoutRead  time.Duration `env:"SERVER_TIMEOUT_READ,required"`
	TimeoutWrite time.Duration `env:"SERVER_TIMEOUT_WRITE,required"`
}

type ConfDB struct {
	Host     string `env:"DB_HOST"`
	Port     int    `env:"DB_PORT"`
	Debug    bool   `env:"DB_DEBUG"`
	DBName   string `env:"DB_NAME"`
	Username string `env:"DB_USER"`
	Password string `env:"DB_PASS"`
}

type ConfAuth struct {
	MaxAge int `env:"AUTH_MAX_AGE,required"`
}

func New() *Conf {
	var c Conf
	if err := envdecode.StrictDecode(&c); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}

	AppEnv = c.Server.ENV
	IsProduction = c.Server.ENV == PROD
	c.IsProduction = IsProduction

	return &c
}
