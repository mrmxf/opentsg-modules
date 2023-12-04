// package errhandle is used for generating logs and writing any errors to logs
package errhandle

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

/*
insert an embed struct here or something for our own error processing
now there is a new struct for errors these changes can be added more easily and incorporated.

type idea struct {
	error string
	level warn debug info error and anyothers
	wrapped
	count e.g. if this error number is repeated for several wrapped errors then just slap it in for certain levels of user base
}

errors.Is(err, fs.ErrExist) for enumerating errors

 simple way to create wrapped errors is to call fmt.Errorf and apply the %w verb to the error argument:

errors.Unwrap(fmt.Errorf("... %w ...", ..., err, ...))

the error design is the error number - main text as a go template
design is translation of message
and translation of generic items

Errors that I have to wrap:
- file systems errors these are provided
- yaml unmarshall errors




*/

// logInit initialises a log with a writer specified by the input. If no input is given it does not write.
// mnt gives a mount point for any files/folders to be written
// It initialises a middleware that prefixes the error with
// the time the error was written.
// The keys are:
//
// - "stdout" this pipes to standard out
//
// -"stderr" this pipes to standard error
//
// - ^file:[a-zA-Z0-9\.\/]{1,30}\.[lL][oO][gG]$ this generates a log file with the same name
//
// - "file" this generates a folder with a  log name in the time format of 2006-01-02_150405
func LogInit(logType, mnt string) *Logger {
	logType = strings.ToLower(logType)
	logs := log.Default()
	fileName := regexp.MustCompile(`^file:[a-zA-Z0-9\.\/]{1,30}\.[lL][oO][gG]$`)
	//middle := handler{}
	switch {
	case logType == "stdout":
		log.SetOutput(os.Stdout)
	case logType == "stderr":
		log.SetOutput(os.Stderr)
	case logType == "file":
		folTformat := "2006-01-02"
		fol := time.Now().Format(folTformat)
		log.SetOutput(logToFile(mnt + "open-tpg_" + fol + ".log"))
	case fileName.MatchString(logType):
		log.SetOutput(logToFile(filepath.Join(mnt, logType[5:])))
	default:
		log.SetOutput(io.Discard)
	}
	// set the writer as the middleware writer
	//log.SetOutput(middle)
	// set so there's no flags and we use our own magic
	logs.SetFlags(0)

	//logsWrapper :=
	//	fmt.Println(logsWrapper)
	return &Logger{log: logs, errors: make(chan loggedMessage, 1000)}
}

func logToFile(name string) io.Writer {
	logFile, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666) // create a log file and panic if it can't run
	if err != nil {
		log.Fatal(fmt.Errorf("error initialising the logging: %v", err))
	}

	return logFile
}

type loggedMessage struct {
	errorMessage, errorTime string
}

type Logger struct {
	//channel is safe for concurrency when the logger is being used by several widgets
	errors chan loggedMessage
	log    *log.Logger
	prefix string
}

// ErrorCount returns the number of error
func (l *Logger) ErrorCount() int {
	return len(l.errors)

}

// PrintErrorMessage writes a string to the log with the prefix. If debug mode is off than just the
// error number is returned, instead of the full message.
func (l *Logger) PrintErrorMessage(prefix string, message error, debug bool) {

	var completeError string
	errTime := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")

	if debug {
		completeError = fmt.Sprintf("%v%v", prefix, message)
	} else {
		mes := message.Error()
		if len(mes) > 4 {
			completeError = fmt.Sprintf("%v%v", prefix, mes[0:4])
		} else {
			completeError = fmt.Sprintf("%v%v", prefix, mes)
		}
	}

	// error handling for if there's too many errors filling up the channel,
	// so something else happens instead of blocking
	select {
	case l.errors <- loggedMessage{errorMessage: completeError, errorTime: errTime}:
	default:
		erroCount := len(l.errors)
		l.LogFlush()
		panic(fmt.Sprintf("channel is full, %v\n", erroCount))
		//panic errors out of range, more errors assigned then we have capacity for etc?
	}
}

// LogFlush flushes the error messages, writing the stored time of error, the prefix and the error message in that order
func (l *Logger) LogFlush() {
	//loop through the channel until its empty
	for len(l.errors) > 0 {
		mes := <-l.errors
		l.log.Printf("%v %v%v", mes.errorTime, l.prefix, mes.errorMessage)
	}
}

// SerPrefix is a pseudo wrapper for the log.SetPrefix function
func (l *Logger) SetPrefix(prefix string) {
	l.prefix = prefix
}

