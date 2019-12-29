package main

import (
	"os"
	"time"

	"github.com/Napat/gol"
	"github.com/sirupsen/logrus"
)

func customGol() *gol.Gol {
	g := &gol.Gol{
		Logger: &logrus.Logger{
			Out: os.Stderr,
			Formatter: &logrus.TextFormatter{
				DisableColors:   true,
				TimestampFormat: time.RFC3339, // time format RFC3339
			},
			Hooks:        make(logrus.LevelHooks),
			Level:        logrus.InfoLevel, // disable debug level
			ExitFunc:     os.Exit,
			ReportCaller: false,
		},
		EnableAutoSync: false, // disable gol sync
	}

	return g
}

func main() {

	x := customGol()

	x.GDebugf("GDebugf debug") // will not print
	x.GInfof("GInfof info")

}
