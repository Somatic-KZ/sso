package logger

import (
	stdlog "log"
	"os"

	"github.com/hashicorp/logutils"
)

func Setup(mode logutils.LogLevel) {
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "WARN", "ERROR"},
		MinLevel: mode,
		Writer:   os.Stdout,
	}

	switch mode {
	case "DEBUG":
		stdlog.SetFlags(stdlog.Ldate | stdlog.Ltime | stdlog.Lmicroseconds | stdlog.Lshortfile)
	default:
		stdlog.SetFlags(stdlog.Ldate | stdlog.Ltime)
	}

	stdlog.SetOutput(filter)
}
