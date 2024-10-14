package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog/log"
)

func New() Config {
	var c Config

	if err := cleanenv.ReadEnv(&c); err != nil {
		log.Fatal().Err(err).Msg("error reading config")
	}
	return c
}

// could add config for PoW algorithm but need to forward to pow package
type Config struct {
	TCP
}

type TCP struct {
	Host string `env:"SVC_HOST" env-default:"localhost"`
	Port string `env:"SVC_PORT" env-default:":9090"`
}

func (c Config) Addr() string {
	return c.Host + c.Port
}
