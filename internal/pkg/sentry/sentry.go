package sentry

import (
	"github.com/getsentry/sentry-go"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
)

var log = logging.MustGetLogger("sentry")

func InitSentry() {
	dsn := viper.GetString("sentry.DSN")
	if err := sentry.Init(sentry.ClientOptions{
		Dsn: viper.GetString("sentry.DSN"),
	}); err != nil {
		log.Warningf("initialization failed: %v\n", err)
	}
	log.Infof("sentry registered with DSN: %v", dsn)
}
