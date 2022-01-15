package logging

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
)

var Log Logger

type LogLevel int64

const (
	Fatal LogLevel = iota
	Error
	Warn
	Info
	Debug
	Trace
)
const maxLevel = 5

var Levels = []LogLevel{
	Fatal,
	Error,
	Warn,
	Info,
	Debug,
	Trace,
}

func ListLogLevels() string {
	ll := [maxLevel + 1]string{}

	for i, l := range Levels {
		ll[i] = fmt.Sprintf("%d - %s", l, l.String())
	}

	return strings.Join(ll[:], ", ")
}

func (ll LogLevel) String() string {
	switch ll {
	case Trace:
		return "TRACE"
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	case Warn:
		return "WARN"
	case Error:
		return "ERROR"
	case Fatal:
		return "FATAL"
	}
	return ""
}

func (ll LogLevel) Color() *color.Color {
	switch ll {
	case Trace:
		return color.New(color.FgWhite, color.Faint)
	case Debug:
		return color.New(color.FgWhite, color.Italic)
	case Info:
		return color.New(color.FgWhite)
	case Warn:
		return color.New(color.FgYellow)
	case Error:
		return color.New(color.FgRed)
	case Fatal:
		return color.New(color.BgRed, color.FgBlack)
	}
	return color.New(color.FgWhite)
}

func (ll LogLevel) Valid() LogLevel {
	if ll < 0 {
		return Fatal
	}

	if ll > maxLevel {
		return Trace
	}

	return ll
}

type Logger struct {
	logLevel      LogLevel
	traceLogger   *log.Logger
	debugLogger   *log.Logger
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
	fatalLogger   *log.Logger
}

func (logger Logger) logMessage(level LogLevel, log *log.Logger, format string, v ...interface{}) {
	if logger.logLevel >= level {
		level.Color().Set()
		defer color.Unset()

		if len(v) == 0 {
			log.Output(3, fmt.Sprintln(format))
		} else {
			log.Output(3, fmt.Sprintln(format, v))
		}
	}
}

func (logger Logger) Trace(format string, v ...interface{}) {
	logger.logMessage(Trace, logger.traceLogger, format, v...)
}

func (logger Logger) Debug(format string, v ...interface{}) {
	logger.logMessage(Debug, logger.debugLogger, format, v...)
}

func (logger Logger) Info(format string, v ...interface{}) {
	logger.logMessage(Info, logger.infoLogger, format, v...)
}

func (logger Logger) Warn(format string, v ...interface{}) {
	logger.logMessage(Warn, logger.warningLogger, format, v...)
}

func (logger Logger) Error(format string, v ...interface{}) {
	logger.logMessage(Error, logger.errorLogger, format, v...)
}

func (logger Logger) Fatal(format string, v ...interface{}) {
	logger.logMessage(Fatal, logger.fatalLogger, format, v...)
	os.Exit(1)
}

func (logger Logger) SetLogLevel(logLevel LogLevel) {
	logger.logLevel = logLevel
}

func init() {
	Log = Logger{
		traceLogger:   log.New(os.Stdout, Trace.String()+": ", log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix),
		debugLogger:   log.New(os.Stdout, Debug.String()+": ", log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix),
		infoLogger:    log.New(os.Stdout, Info.String()+": ", log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix),
		warningLogger: log.New(os.Stdout, Warn.String()+": ", log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix),
		errorLogger:   log.New(os.Stderr, Error.String()+": ", log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix),
		fatalLogger:   log.New(os.Stderr, Fatal.String()+": ", log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix),
	}
}
