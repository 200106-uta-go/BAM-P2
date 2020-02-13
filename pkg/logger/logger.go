package logger

import (
	"io"
	"log"
	"os"
)

const logPath = "./logs/project2log.log"

var (
	Logger *log.Logger // Writes a log to logfile
	logchan chan // Channel for exporting log to log server
)
func init() {
	logFile, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(err)
	}

	// Creates io.Writer that writes to both the console and the log file
	multiWriter := io.MultiWriter(logFile, os.Stdout)

	// Creates a pointer to a multiwriter log.Logger
	Logger = log.New(multiWriter, "Project-2 Log:\t", log.LstdFlags|log.Lshortfile)

	/*
	// Creates io.Writer that writes to both the console and a log file to 
	// a log server
	multiWriter := io.MultiWriter(logchan<-, os.Stdout)

	// Creates a pointer to a multiwriter log.Logger
	Logger = log.New(multiWriter, "Project-2 Log:\t", log.LstdFlags|log.Lshortfile)
	*/
}
