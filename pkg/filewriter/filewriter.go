// Package filewriter is a group of helper functions for creating,
// writing, and reading from files
package filewriter

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// CheckForFile checks to see if file exist
func CheckForFile(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// CreateFile creates a file with filename (string) as its name
// and returns true (bool) if succesful
func CreateFile(filename string) bool {
	_, err := os.Create(filename)
	if err != nil {
		return false
	}
	return true
}

// WriteNew looks for filename (string) [if not present creates it]
// and fills it with message ([]byte) converted to a string and
// returns true (bool) if succesful
func WriteNew(filename string, message []byte) bool {
	if !CheckForFile(filename) {
		if !CreateFile(filename) {
			log.Println("Unable to find or create file: " + filename)
			return false
		}
	} else {
		f, err := os.Open(filename)
		if err != nil {
			log.Println("Unable to open file: " + filename)
			return false
		}
		defer f.Close()

		fmt.Fprint(f, string(message))
	}
	return true
}

// WriteAppend looks for filename (string) [if not present creates it]
// and appends the message ([]byte) converted to string to the EOF and
// returns true (bool) if succesful
func WriteAppend(filename string, message []byte) bool {
	if !CheckForFile(filename) {
		if !CreateFile(filename) {
			log.Println("Unable to find or create file: " + filename)
			return false
		}
	} else {
		f, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
		if err != nil {
			log.Println("Unable to open file: " + filename)
			return false
		}
		defer f.Close()

		fmt.Fprintln(f, string(message))
	}
	return true
}

// WriteRaw looks for filename (string) [if not present creates it]
// and appends the message ([]byte) as raw bytes to the EOF and
// returns true (bool) if succesful
func WriteRaw(filename string, message []byte) bool {
	if !CheckForFile(filename) {
		if !CreateFile(filename) {
			log.Println("Unable to find or create file: " + filename)
			return false
		}
	} else {
		f, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
		if err != nil {
			log.Println("Unable to open file: " + filename)
			return false
		}
		defer f.Close()

		fmt.Fprint(f, message)
	}
	return true
}

// ReadFile reads from the filename (string) and returns the results
// as a []byte slice
func ReadFile(filename string) []byte {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("File reading error", err)
		return nil
	}
	return data
}
