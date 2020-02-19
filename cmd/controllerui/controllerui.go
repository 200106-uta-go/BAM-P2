package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/200106-uta-go/BAM-P2/pkg/httputil"
)

// Request is a struct to hold the values included in
// the requests from the web client
type Request struct {
	Filepath string `json:"filepath"`
	Object   string `json:"object"`
	Name     string `json:"name"`
	Replicas string `json:"replicas"`
}

var port = ":4040"

func main() {
	//setup fileserver for serving index.html
	fs := http.FileServer(http.Dir("./"))

	//handle routes used by client
	http.HandleFunc("/apply", apply)
	http.HandleFunc("/delete", delete)
	http.HandleFunc("/get", get)
	http.HandleFunc("/describe", describe)
	http.HandleFunc("/scale", scale)
	http.HandleFunc("/logs", logs)
	http.HandleFunc("/cluster", cluster)

	fmt.Println("Server listening on localhost", port)
	log.Fatalln(http.ListenAndServe(port, fs))
}

func apply(w http.ResponseWriter, r *http.Request) {
	//get filepath from request
	body := readBody(r)
	newRequest := jsonToRequest(body)

	controller.apply(newRequest.Filepath)
	w = httputil.SetHeaders(w)
	w.Write([]byte("OK"))
}

func delete(w http.ResponseWriter, r *http.Request) {
	//get object and name from request
	body := readBody(r)
	newRequest := jsonToRequest(body)

	controller.delete(newRequest.Object, newRequest.Name)
	w = httputil.SetHeaders(w)
	w.Write([]byte("OK"))
}

func get(w http.ResponseWriter, r *http.Request) {
	//get object and name from request
	body := readBody(r)
	newRequest := jsonToRequest(body)

	json := controller.get(newRequest.Object, newRequest.Name)
	w = httputil.SetHeaders(w)
	w.Write(json)
}

func describe(w http.ResponseWriter, r *http.Request) {
	//get object and name from request
	body := readBody(r)
	newRequest := jsonToRequest(body)

	json := controller.describe(newRequest.Object, newRequest.Name)
	w = httputil.SetHeaders(w)
	w.Write(json)
}

func scale(w http.ResponseWriter, r *http.Request) {
	//get replicas and deployment filepath from request
	body := readBody(r)
	newRequest := jsonToRequest(body)

	controller.scale(newRequest.Replicas, newRequest.Filepath)
	w = httputil.SetHeaders(w)
	w.Write([]byte("OK"))
}

func logs(w http.ResponseWriter, r *http.Request) {
	//get pod name from request
	body := readBody(r)
	newRequest := jsonToRequest(body)

	json := controller.logs(newRequest.Name)
	w = httputil.SetHeaders(w)
	w.Write(json)
}

func cluster(w http.ResponseWriter, r *http.Request) {
	//

	json := controller.cluster()
	w = httputil.SetHeaders(w)
	w.Write(json)
}

func readBody(r *http.Request) string {
	body, err := ioutil.ReadAll(r.Body)
	httputil.GenericErrHandler("error", err)
	return string(body)
}

func (r Request) tojson() string {
	bytes, err := json.Marshal(r)
	httputil.GenericErrHandler("error", err)
	return string(bytes)
}

func jsonToRequest(js string) Request {
	newRequest := Request{}
	err := json.Unmarshal([]byte(js), &newRequest)
	httputil.GenericErrHandler("error", err)
	return newRequest
}
