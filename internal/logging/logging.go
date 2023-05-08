package logging

import (
	"io"
	"log"
	"os"
)

/**
TODO: add logging to file
	* create logger factory
*/

var (
	loggerCache map[string]*log.Logger
)

func init() {
	loggerCache = map[string]*log.Logger{}
}

func GetLoggerFor(pkg string) *log.Logger {
	return GetLoggerForWithFileOutput(pkg, nil)
}

func GetLoggerForWithFileOutput(pkg string, fileOutput *os.File) *log.Logger {
	if logger, ok := loggerCache[pkg]; ok {
		return logger
	} else {
		var logWriters io.Writer
		if fileOutput == nil {
			logWriters = os.Stdout
		}
		logWriters = io.MultiWriter(os.Stdout, fileOutput)
		l := log.New(logWriters, pkg+": ", log.Ldate|log.Ltime|log.Lshortfile)
		loggerCache["pkg"] = l
		return l
	}
}
