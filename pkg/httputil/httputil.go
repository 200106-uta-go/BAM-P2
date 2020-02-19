package httputil

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// HTTPResponse defines a generic struct for sending a http response message as JSON
type HTTPResponse struct {
	Message string `json:"message"`
}

// SetHeaders sets the response headers for an outgoing HTTP response
func SetHeaders(w http.ResponseWriter) http.ResponseWriter {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
	return w
}

// GoodRequest sends an "OK" message along with the 200 status code
func GoodRequest(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintln(w, "OK")
	GenericErrHandler("error", err)
}

// BadRequest sends an HTTP response with a Bad Request status code along
// with the message passed into the function back to the client in JSON format
func BadRequest(w http.ResponseWriter, message string) {
	response := HTTPResponse{
		Message: message,
	}
	r, err := json.Marshal(response)
	GenericErrHandler("error", err)

	w.WriteHeader(http.StatusBadRequest)
	w.Write(r)
}

// GenericErrHandler is a function to replace the common
// generic error handler written throughout the code
func GenericErrHandler(action string, err error) {
	switch action {
	case "print":
		if err != nil {
			log.Println(err)
		}
	case "error":
		if err != nil {
			log.Fatalln(err)
		}
	default:
		if err != nil {
			log.Fatalln(err)
		}
	}
}
