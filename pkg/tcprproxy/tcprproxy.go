// Package tcprproxy is a TCP reverse proxy adaptation
package tcprproxy

// Import packages
import (
	"bytes"
	"log"
	"math/rand"
	"net"
	"time"

	l "github.com/Gamemastertwig/go-secure/loghelper"
)

// Global Variables
var connSig, done, exit chan string
var incoming, outgoing chan []byte

// TCPForward uses tcp connections to establish a reverse proxy
func TCPForward(front string, back string, logger string, program string) {
	// create tcp connection
	listn, err := net.Listen("tcp", front)
	if err != nil {
		l.ConnLogMess(logger, program+" ERROR: ", front+" Failed to setup listiner for RPROXY")
		log.Fatalf(program+" ERROR: ", front+" Failed to setup listiner for RPROXY")
	}

	l.ConnLogMess(logger, program+" LOG: ", "(TCP) Listening on "+front)
	log.Println(program+" LOG: ", "(TCP) Listening on "+front)

	// create connection signal channel
	connSig = make(chan string)

	// allow for multiple connections
	for {
		go Session(listn, front, back, logger, program)
		<-connSig
	}
}

// TCPForwardLB uses tcp connections to establish a reverse proxy with load balancer
func TCPForwardLB(front string, backends []string, logAddr string, program string) {
	// create tcp connection
	listn, err := net.Listen("tcp", front)
	if err != nil {
		l.ConnLogMess(logAddr, program+" ERROR: ", front+" Failed to setup listiner for LB RPROXY")
		log.Fatalf(program+" ERROR: ", front+" Failed to setup listiner for LB RPROXY")
	}

	l.ConnLogMess(logAddr, program+" LOG: ", "(TCP) Listening on "+front)
	log.Println(program+" LOG: ", "(TCP) Listening on "+front)

	// create connection signal channel
	connSig = make(chan string)

	// allow for multiple connections
	for {
		// call loadBalance here
		bAddr := loadBalanceRand(backends, logAddr, program)
		go Session(listn, front, bAddr, logAddr, program)
		<-connSig
	}
}

// Not working correctly [to do]
func loadBalanceRand(backends []string, logAddr string, program string) string {
	// verify back ends
	for i, b := range backends {
		TempConn, err := net.Dial("tcp", b)
		if err != nil {
			// remove from slice ([]sting) if it fails to make a connection
			backends = append(backends[:i], backends[i+1:]...)
		}
		TempConn.Close()
	}

	// seed rand based on time per OS
	rand.Seed(time.Now().UnixNano())
	// get a random int from 0 to length of string slice -1
	lenBack := len(backends)
	n := rand.Intn(lenBack)

	l.ConnLogMess(logAddr, program+" LOG: ", "Sending load to "+backends[n])
	log.Println(program+" LOG: ", "Sending load to "+backends[n])

	return backends[n]
}

// Session allows for multiple connections from clients at the same time
// listening on front end (net.Listener) and then connects to back end
// address (backAddr string)
func Session(listn net.Listener, frontAddr string, backAddr string, logAddr string, program string) {
	// wait for front end connection
	frontConn, err := listn.Accept()
	if err != nil {
		l.ConnLogMess(logAddr, program+" ERROR:", "Failed to accept connection:: "+err.Error())
		log.Fatalln(program + " ERROR: Failed to accept connection:: " + err.Error())
	}
	l.ConnLogMess(logAddr, program+" LOG:", "Accepted Connection from "+frontConn.LocalAddr().String())
	log.Println(program + " LOG: Accepted Connection from " + frontConn.LocalAddr().String())

	// defer close : member LIFO
	defer frontConn.Close()

	// send message to allow for an new connection
	connSig <- "Done"

	// create connection for server end
	serverConn, err := net.Dial("tcp", backAddr)
	if err != nil {
		l.ConnLogMess(logAddr, program+" ERROR:", "Dial failed for address "+backAddr+":: "+err.Error())
		log.Fatalln(program + " ERROR: Dial failed for address " + backAddr + ":: " + err.Error())
	}
	l.ConnLogMess(logAddr, program+" LOG:", "Dial succesful to "+backAddr)
	log.Println(program + " LOG: Dial succesful to " + backAddr)

	// defer close
	defer serverConn.Close()

	// create message channels
	done = make(chan string)
	incoming = make(chan []byte)
	outgoing = make(chan []byte)

	// listen for message from client and log request
	go TCPListen(frontConn, incoming, logAddr, frontAddr, program)
	// serve message from client to server
	go TCPServe(serverConn, incoming, logAddr, backAddr, program)
	// listen for response from server and log request
	go TCPListen(serverConn, outgoing, logAddr, frontAddr, program)
	// server message from server to client
	go TCPServe(frontConn, outgoing, logAddr, backAddr, program)
	<-done
}

// TCPListen listens on conn (net.Conn) and reads the data ([]byte) into a
// string channel
func TCPListen(conn net.Conn, packet chan []byte, logAddr string, fromAddr string, program string) {
	for {
		// create read buffer
		rBuf := make([]byte, 1024)

		err := conn.SetDeadline(time.Now().Add(1 * time.Second))
		if err != nil {
			l.ConnLogMess(logAddr, program+" ERROR:", "Failed deadline setup for "+
				conn.LocalAddr().String()+":: "+err.Error())
			log.Println(program + " ERROR: Failed deadline setup for " +
				conn.LocalAddr().String() + ":: " + err.Error())
			break
		}

		// Attempt read
		_, err = conn.Read(rBuf)
		if err != nil {
			break
		}

		cBuf := bytes.Trim(rBuf, "\x00")

		// place buffer in packet channel
		packet <- cBuf

		l.ConnLogMess(logAddr, program+" LOG: "+fromAddr+":", "Packet read.")
		log.Println(program + " LOG: " + fromAddr + ": Packet read.")
	}
	done <- "Done"
}

// TCPServe serves the data ([]byte) from the packet (string channel)
// to the connection (net.Conn)
func TCPServe(conn net.Conn, packet chan []byte, logAddr string, fromAddr string, program string) {
	for {
		// get buffer from packet channel
		mBuf := <-packet

		// connection exist write buffer to connection
		if conn != nil {
			conn.Write(mBuf)

			l.ConnLogMess(logAddr, program+" LOG: "+fromAddr+":", "Packet sent.")
			log.Println(program + " LOG: " + fromAddr + ": Packet sent.")
		}
	}
}
