// Package tcplogger is a package to assist other applications
// with connecting to the TCP logging server
package tcplogger

// imports
import (
	"log"
	"net"
	"time"
)

// TCPMessage is a struct used to hold information to be sent to the
// TCP logging server
type TCPMessage struct {
	LogAddr     string
	MessageType string
	Message     string
}

// ConnLogMess creates a connection the log server and sends it a message
// then closes it
func ConnLogMess(logAddr string, mType string, message string) {
	var logConn net.Conn
	currentTime := time.Now()

	// create connection to logging server
	logConn, _ = net.Dial("tcp", logAddr)
	if logConn == nil {
		///log.Printf("Dial to logger at %+v failed:: %+v ", logAddr, err)
	} else {
		log.Println("Connected to logging server at " + logConn.LocalAddr().String())
		mes := currentTime.Format("2006/01/02 15:04:05") + " " + mType + " " + message
		logConn.Write([]byte(mes))
	}
}

func (m *TCPMessage) Write() {
	// connect to log server
	logConn, _ := net.Dial("tcp", m.LogAddr)
	if logConn == nil {
		log.Printf("Dial to log server failed at %+v", m.LogAddr)
	} else {
		// 
		if m.MessageType == "" {
			m.MessageType = "TCP LOG"
		}
		mes := m.MessageType + " " + m.Message
		logConn.Write([]byte(mes))
	}
}
