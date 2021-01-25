package agscheduler

import (
	"github.com/sirupsen/logrus"
	"github.com/timest/env"
	"time"
)

var Config AGSConfig

type AGSConfig struct {
	Log LogConfig
	PG  PGConfig
}

type LogConfig struct {
	Level string `default:"info"`
	Json  bool   `default:"false"`
}

type PGConfig struct {
	Addr     string `default:"localhost:5432"`
	User     string `default:"postgres"`
	Password string `default:"postgres"`
	Database string `default:"postgres"`
	PoolSize int    `default:"3"`
}

func init() {
	config := AGSConfig{}
	env.IgnorePrefix()
	err := env.Fill(&config)
	if err != nil {
		AGSLog.WithFields(GenASGModule("config")).
			WithError(err).
			Warningln("init env failure.")
	}
	Config = config
	logLevel, err := logrus.ParseLevel(Config.Log.Level)
	if err == nil {
		AGSLog.SetLevel(logLevel)
	}
	if Config.Log.Json {
		AGSLog.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339Nano})
	}
}
