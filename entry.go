package gol

import (
	"github.com/sirupsen/logrus"
)

type Entry struct {
	logrus.Entry
	Logger *Gol
}

// GDebugf developer log & system information we don't want user to see.
func (entry *Entry) GDebugf(format string, args ...interface{}) {
	entry.Debugf(format, args...)

	if entry.Logger.EnableAutoSync == true {
		entry.Logger.FileSync()
	}
}

// GInfof user log
func (entry *Entry) GInfof(format string, args ...interface{}) {
	entry.Infof(format, args...)

	if entry.Logger.EnableAutoSync == true {
		entry.Logger.FileSync()
	}
}

// GErrorf user log case error
func (entry *Entry) GErrorf(format string, args ...interface{}) {
	entry.Errorf(format, args...)

	if entry.Logger.EnableAutoSync == true {
		entry.Logger.FileSync()
	}
}

// Printf similar to GInfof
func (entry *Entry) Printf(format string, args ...interface{}) {
	entry.GInfof(format, args)
}
