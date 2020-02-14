package logger

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/200106-uta-go/BAM-P2/pkg/filewriter"
)

const logPath = "current.log"

//var logchan chan // Channel for exporting log to log server
func init() {
	var multiWriter io.Writer
	var logSvr, logF bool
	var logConn net.Conn
	var logFile *os.File

	err := os.Setenv("LOG_SVR_ADDR", ":9090")
	if err != nil {
		// failed to create env variable
	} else {
		// connect to log server
		logConn, err = net.Dial("tcp", os.Getenv("LOG_SVR_ADDR"))
		if err != nil {
			log.Printf("Failed to connect to %+v :: %+v", os.Getenv("LOG_SVR_ADDR"), err)
			logSvr = false
		} else {
			logSvr = true
		}
	}

	// make sure logSrv is false if Setenv fails
	if logConn == nil {
		logSvr = false
	}

	// create file writer
	if !filewriter.CheckForFile(logPath) {
		if filewriter.CreateFile(logPath) {
			logFile, err = os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				log.Printf("Failed to open file %+v :: %+v", logPath, err)
				logF = false
			} else {
				logF = true
			}
		} else {
			logF = false
		}
	} else {
		logFile, err = os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Printf("Failed to open file %+v :: %+v", logPath, err)
			logF = false
		} else {
			logF = true
		}
	}

	// Create MultiWrite based on if a connections
	switch {
	case logF && logSvr:
		multiWriter = io.MultiWriter(logFile, logConn, os.Stdout)
		fmt.Println("Output to SERVER, FILE, and STDOUT")
	case logF && !logSvr:
		multiWriter = io.MultiWriter(logFile, os.Stdout)
		fmt.Println("Output to FILE and STDOUT")
	case logSvr && !logF:
		multiWriter = io.MultiWriter(logConn, os.Stdout)
		fmt.Println("Output to SERVER and STDOUT")
	default:
		multiWriter = io.MultiWriter(os.Stdout)
		fmt.Println("Output to STDOUT only")
	}

	// sets log to use multiWriter
	log.SetOutput(multiWriter)
}
