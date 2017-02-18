package applog

import (
	"io"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

//------------------------------------------------------------------------------
//------------------------------------------------------------------------------

var (
	// Trace is the trace logger
	Trace *log.Logger

	// Info is the info logger.
	Info *log.Logger

	// Warn is the warning logger.
	Warn *log.Logger

	// Error is the error logger.
	Error *log.Logger

	DefaultOutput io.Writer
)

func init() {
	StartupNoOp()
}

// StartupNoOp starts up the logging system in a non-operational mode.
func StartupNoOp() {
	Trace = log.New(ioutil.Discard, "", log.LstdFlags)
	Info = log.New(ioutil.Discard, "", log.LstdFlags)
	Warn = log.New(ioutil.Discard, "", log.LstdFlags)
	Error = log.New(ioutil.Discard, "", log.LstdFlags)
}

// StartupInteractive starts up the logging system for an interactive application.
func StartupInteractive() {
	var traceWriter io.Writer = os.Stdout

	fileLogger := &lumberjack.Logger{
		Filename:   "discordeq.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     7, //days
	}

	multiOut := io.MultiWriter(fileLogger, os.Stdout)
	multiErr := io.MultiWriter(fileLogger, os.Stderr)
	DefaultOutput = multiOut

	Trace = log.New(traceWriter, "[TRACE] ", log.LstdFlags)
	Info = log.New(multiOut, "[INFO] ", log.Ldate|log.Ltime|log.LstdFlags)
	Warn = log.New(multiOut, "[WARN] ", log.Ldate|log.Ltime|log.LstdFlags)
	Error = log.New(multiErr, "[ERROR] ", log.Ldate|log.Ltime|log.LstdFlags)
}
