package logger

import(
	log "github.com/sirupsen/logrus"
)

func Init(logLevel *string)  {
	ll, err := log.ParseLevel(*logLevel)
	if err != nil {
		panic(err)
	}

	log.SetLevel(ll)
	log.SetFormatter(&log.TextFormatter{})
}