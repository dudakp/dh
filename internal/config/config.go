package config

// TODO: refactor this to use zerolog

import (
	"log"
	"os"
)

var (
	InfoLog *log.Logger
	WarnLog *log.Logger
	ErrLog  *log.Logger
)

func init() {
	InfoLog = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarnLog = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrLog = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
