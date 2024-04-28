package main

import (
	"log"
	"os"
	"path/filepath"
	"sync"
)

const loggerFlag = log.Ldate | log.Ltime | log.Lmsgprefix

func initLogger() {
	log.SetFlags(loggerFlag)
	log.SetPrefix(filepath.Base(os.Args[0]) + ": ")
}

var loggers sync.Map

func getLogger(fn string) *log.Logger {
	if v, ok := loggers.Load(fn); ok {
		return v.(*log.Logger)
	}

	l := log.New(os.Stderr, fn+": ", loggerFlag)
	loggers.Store(fn, l)
	return l
}
