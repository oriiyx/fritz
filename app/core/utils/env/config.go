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
	ENV          string        `env:"APP_ENV,required"`
	Port         int           `env:"SERVER_PORT,required"`
	Debug        bool          `env:"SERVER_DEBUG,required"`
	Secret       []byte        `env:"SECRET_KEY,required"`
	SupportMail  string        `env:"SUPPORT_MAIL,required"`
	TimeoutIdle  time.Duration `env:"SERVER_TIMEOUT_IDLE,required"`
	TimeoutRead  time.Duration `env:"SERVER_TIMEOUT_READ,required"`
	TimeoutWrite time.Duration `env:"SERVER_TIMEOUT_WRITE,required"`
}

type ConfDB struct {
	Host     string `env:"DB_HOST"`
	Port     int    `env:"DB_PORT"`
	DBURL    string `env:"DATABASE_URL,required"`
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
