package gol

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	_ "github.com/snowzach/rotatefilehook" // https://github.com/sirupsen/logrus/issues/784#issuecomment-403765306
)

type Fields logrus.Fields

// Gol logrus.Logger with additional features
type Gol struct {
	*logrus.Logger
	// Mutex Locking is enabled by Default to protect OSFile and sensitive feilds
	MU             logrus.MutexWrap
	OSFile         *os.File
	FullPath       string
	EnableAutoSync bool // Disable by default
}

// New Creates a new gol logger. Configuration should be set by changing `Formatter`,
// `Out` and `Hooks` directly on the default logger instance. You can also just
// instantiate your own:
//
//  var log = &Gol{
// 		Logger: &logrus.Logger{
// 			Out:          os.Stderr,
// 			Formatter:    new(logrus.TextFormatter),
// 			Hooks:        make(logrus.LevelHooks),
// 			Level:        logrus.InfoLevel,
//	 	},
// 		OSFile: nil,
// 		FullPath: "",
//    }
//
// It's recommended to make this a global instance called `log`.
func New() *Gol {
	g := &Gol{
		Logger: &logrus.Logger{
			Out: os.Stderr,
			Formatter: &logrus.TextFormatter{
				DisableColors:   true,                            // https://github.com/sirupsen/logrus/issues/1069
				TimestampFormat: "2006-01-02T15:04:05.000Z07:00", // add timestamp with format: RFC3339 Micro
			},
			Hooks:        make(logrus.LevelHooks),
			Level:        logrus.DebugLevel,
			ExitFunc:     os.Exit,
			ReportCaller: false,
		},
		// OSFile: nil,
		// FullPath: "",
	}

	return g
}

func NewWithHookPrinter() *Gol {
	g := &Gol{
		Logger: &logrus.Logger{
			Out: os.Stderr,
			Formatter: &logrus.TextFormatter{
				DisableColors:   true,                            // https://github.com/sirupsen/logrus/issues/1069
				TimestampFormat: "2006-01-02T15:04:05.000Z07:00", // add timestamp with format: RFC3339 Micro
			},
			Hooks:        make(logrus.LevelHooks),
			Level:        logrus.DebugLevel,
			ExitFunc:     os.Exit,
			ReportCaller: false,
		},
		// OSFile: nil,
		// FullPath: "",
	}

	g.AddHook(&HookPrinter{})
	return g
}

func NewWithHookAutoFileDate(filepath string, rotatefileTimezone string) (gol *Gol, err error) {
	var tloc *time.Location

	gol = &Gol{
		Logger: &logrus.Logger{
			Out: os.Stderr,
			Formatter: &logrus.TextFormatter{
				DisableColors:   true,                            // https://github.com/sirupsen/logrus/issues/1069
				TimestampFormat: "2006-01-02T15:04:05.000Z07:00", // add timestamp with format: RFC3339 Micro
			},
			Hooks:        make(logrus.LevelHooks),
			Level:        logrus.DebugLevel,
			ExitFunc:     os.Exit,
			ReportCaller: false,
		},
		// OSFile: nil,
		FullPath:       filepath,
		EnableAutoSync: true,
	}

	if rotatefileTimezone == "" {
		tloc, _ = time.LoadLocation("Asia/Bangkok")
	} else {
		tloc, err = time.LoadLocation(rotatefileTimezone)
		if err != nil {
			return nil, err
		}
	}

	h := &HookAutoFileDate{
		Logger:       gol,
		TimeLocation: tloc,
	}

	gol.AddHook(h)
	h.FileDateHandler()

	return gol, nil
}

func NewJSONWithHookAutoFileDate(filepath string, rotatefileTimezone string) (gol *Gol, err error) {
	var tloc *time.Location
	gol = &Gol{
		Logger: &logrus.Logger{
			Out: os.Stderr,
			Formatter: &logrus.JSONFormatter{
				TimestampFormat: "2006-01-02T15:04:05.000Z07:00", // add timestamp with format: RFC3339 Micro
				PrettyPrint:     false,
			},
			Hooks:        make(logrus.LevelHooks),
			Level:        logrus.DebugLevel,
			ExitFunc:     os.Exit,
			ReportCaller: false,
		},
		OSFile:         nil,
		FullPath:       filepath,
		EnableAutoSync: true,
	}

	if rotatefileTimezone == "" {
		tloc, _ = time.LoadLocation("Asia/Bangkok")
	} else {
		tloc, err = time.LoadLocation(rotatefileTimezone)
		if err != nil {
			return nil, err
		}
	}

	h := &HookAutoFileDate{
		Logger:       gol,
		TimeLocation: tloc,
	}

	gol.AddHook(h)
	h.FileDateHandler()

	return gol, nil
}

// SetFileOutput Util method to set output to file
func (gol *Gol) SetFileOutput(filename string) error {
	return gol.setFileOutput(filename, true, true)
}

func (gol *Gol) setFileOutput(filename string, setFullPath bool, muLock bool) error {
	var file *os.File
	var err error

	gol.MU.Lock()
	defer gol.MU.Unlock()

	writerBuf := gol.Out
	osfileBuf := gol.OSFile

	if filename == "" {

		gol.Logger.SetOutput(os.Stderr)

		gol.OSFile = nil
		gol.FullPath = ""

		return nil
	} else {
		if file, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err != nil {
			return err
		}

		file.Sync()

		if muLock == true {
			gol.Logger.SetOutput(file)
		} else {
			gol.Logger.Out = file
		}

		gol.OSFile = file
		if setFullPath == true {
			gol.FullPath = filename
		}
	}

	if writerBuf != nil {
		switch v := writerBuf.(type) {
		case *os.File:
			osfileBuf.Close()
		default:
			fmt.Printf("Not support type %T!\n", v)
			return errors.New(fmt.Sprintf("Not support type %T!\n", v))
		}
	}

	return nil
}

// GWithFields Adds a struct of fields to the log entry.
func (gol *Gol) GWithFields(fields Fields) *Entry {
	return gol.WithFields(fields)
}

// WithFields Shadowing low level WithFields() method to prevent argument type error.
func (gol *Gol) WithFields(fields Fields) *Entry {
	e := gol.Logger.WithFields(logrus.Fields(fields))
	return &Entry{
		Entry:  *e,
		Logger: gol,
	}
}

// FileSync sync os file
func (gol *Gol) FileSync() {
	gol.MU.Lock()
	defer gol.MU.Unlock()

	if gol.OSFile != nil {
		gol.OSFile.Sync()
	}
}

// GDebugf developer log & system information we don't want user to see.
func (gol *Gol) GDebugf(format string, args ...interface{}) {
	gol.Debugf(format, args...)

	if gol.EnableAutoSync == true {
		gol.FileSync()
	}
}

// GInfof user log
func (gol *Gol) GInfof(format string, args ...interface{}) {
	gol.Infof(format, args...)

	if gol.EnableAutoSync == true {
		gol.FileSync()
	}
}

// GErrorf user log case error
func (gol *Gol) GErrorf(format string, args ...interface{}) {
	gol.Errorf(format, args...)

	if gol.EnableAutoSync == true {
		gol.FileSync()
	}
}

// Printf similar to GInfof
func (gol *Gol) Printf(format string, args ...interface{}) {
	gol.GInfof(format, args)
}
