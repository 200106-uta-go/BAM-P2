package logger

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/200106-uta-go/BAM-P2/pkg/filewriter"
	"github.com/joho/godotenv"
)

const logPath = "current.log"

//var logchan chan // Channel for exporting log to log server
func init() {
	var multiWriter io.Writer
	var logSvr, logF bool
	var logConn net.Conn
	var logFile *os.File

	//load in environment variables from .env
	//will print error message when running from docker image
	//because env file is passed into docker run command
	err := godotenv.Load("/home/ubuntu/go/src/github.com/200106-uta-go/BAM-P2/.env")
	if err != nil {
		if !strings.Contains(err.Error(), "no such file or directory") {
			log.Println("Error loading .env: ", err)
		}
	}
	logAddr := os.Getenv("LOG_ADDR")

	if len(logAddr) < 1 {
		// failed to get log address env variable
	} else {
		// connect to log server
		logConn, err = net.Dial("tcp", logAddr)
		if err != nil {
			log.Printf("Failed to connect to %+v :: %+v", logAddr, err)
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
