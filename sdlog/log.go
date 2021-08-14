package sdlog

import (
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gaorx/stardust3/sdparse"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	Init(Options{
		Level:    "debug",
		Stdout:   true,
		Filename: "",
	})
	//logrus.AddHook(&filelineHook{})
}

type Options struct {
	Level      string `json:"level" toml:"level"`
	Stdout     bool   `json:"stdout" toml:"stdout"`
	Filename   string `json:"filename" toml:"filename"`
	MaxSizeMB  int    `json:"max_size" toml:"max_size"`
	MaxAgeDays int    `json:"max_age" toml:"max_age"`
	MaxBackups int    `json:"max_backups" toml:"max_backups"`
}

func Init(opts Options) {
	pretty := sdparse.BoolDef(os.Getenv("SDLOG_PRETTY"), false)
	if opts.MaxSizeMB <= 0 {
		opts.MaxSizeMB = 200
	}

	setPrettyFormat := func(pretty bool) {
		if pretty {
			logrus.SetFormatter(&logrus.TextFormatter{
				DisableColors: false,
				FullTimestamp: true,
			})
		} else {
			logrus.SetFormatter(&logrus.TextFormatter{
				DisableColors: true,
				FullTimestamp: true,
			})
		}
	}

	// level
	var logLevel = logrus.DebugLevel
	if opts.Level != "" {
		logLevel1, err := logrus.ParseLevel(strings.ToLower(opts.Level))
		if err != nil {
			logLevel = logrus.DebugLevel // default
		} else {
			logLevel = logLevel1
		}
	}
	logrus.SetLevel(logLevel)

	// format
	if opts.Filename != "" {
		setPrettyFormat(false)
	} else {
		setPrettyFormat(pretty)
	}

	// output
	if opts.Filename != "" && opts.Stdout {
		// Both file and stdout
		logrus.SetOutput(io.MultiWriter(os.Stdout, &lumberjack.Logger{
			Filename:   opts.Filename,
			MaxSize:    opts.MaxSizeMB,
			MaxAge:     opts.MaxAgeDays,
			MaxBackups: opts.MaxBackups,
		}))
	} else if opts.Filename != "" && !opts.Stdout {
		// Only file
		logrus.SetOutput(&lumberjack.Logger{
			Filename:   opts.Filename,
			MaxSize:    opts.MaxSizeMB,
			MaxAge:     opts.MaxAgeDays,
			MaxBackups: opts.MaxBackups,
		})
	} else if opts.Filename == "" && opts.Stdout {
		// Only stdout
		logrus.SetOutput(os.Stdout)
	} else {
		// No log
		logrus.SetOutput(ioutil.Discard)
	}
}

type Fields = logrus.Fields

var (
	WithFields = logrus.WithFields

	// Level
	SetLevel = logrus.SetLevel
	GetLevel = logrus.GetLevel

	// With info
	WithError = logrus.WithError
	WithField = logrus.WithField

	// Log
	Debug   = logrus.Debug
	Print   = logrus.Print
	Info    = logrus.Info
	Warn    = logrus.Warn
	Warning = logrus.Warning
	Error   = logrus.Error
	Panic   = logrus.Panic
	Fatal   = logrus.Fatal

	// Logf
	Debugf   = logrus.Debugf
	Printf   = logrus.Printf
	Infof    = logrus.Infof
	Warnf    = logrus.Warnf
	Warningf = logrus.Warningf
	Errorf   = logrus.Errorf
	Panicf   = logrus.Panicf
	Fatalf   = logrus.Fatalf
)
