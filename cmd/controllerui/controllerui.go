package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/200106-uta-go/BAM-P2/pkg/controller"
	"github.com/200106-uta-go/BAM-P2/pkg/httputil"
)

// Request is a struct to hold the values included in
// the requests from the web client
type Request struct {
	Deployment string `json:"deployment"`
	Object     string `json:"object"`
	Name       string `json:"name"`
	Replicas   string `json:"replicas"`
}

var port = ":4040"

func main() {
	//setup fileserver for serving index.html
	fs := http.FileServer(http.Dir("web/"))
	http.Handle("/web/", http.StripPrefix("/web/", fs))

	//handle routes used by client
	http.HandleFunc("/apply", apply)
	http.HandleFunc("/delete", delete)
	http.HandleFunc("/get", get)
	http.HandleFunc("/scale", scale)
	http.HandleFunc("/logs", logs)
	http.HandleFunc("/run", run)
	http.HandleFunc("/cluster", cluster)

	exec.Command("xdg-open", "http://localhost:4040/web").Run()
	fmt.Printf("Server listening on localhost%s\n", port)
	fmt.Printf("Access the controller dashbaord at http://localhost%s/web\a \n", port)
	log.Fatalln(http.ListenAndServe(port, nil))
}

func apply(w http.ResponseWriter, r *http.Request) {
	//get filepath from request
	body := readBody(r)
	newRequest := jsonToRequest(body)

	temp, err := os.Create("temp.yaml")
	if err != nil {
		log.Fatalln(err)
	}
	temp.WriteString(newRequest.Deployment)
	temp.Close()

	fmt.Println(newRequest.Deployment)
	controller.KubeApply("./temp.yaml")
	w = httputil.SetHeaders(w)
	w.Write([]byte("OK"))
	os.Remove("temp.yaml")
}

func delete(w http.ResponseWriter, r *http.Request) {
	//get object and name from request
	body := readBody(r)
	newRequest := jsonToRequest(body)

	controller.KubeDelete(newRequest.Object, newRequest.Name)
	w = httputil.SetHeaders(w)
	w.Write([]byte("OK"))
}

func get(w http.ResponseWriter, r *http.Request) {
	//get object and name from request
	body := readBody(r)
	newRequest := jsonToRequest(body)

	json := controller.KubeGet(newRequest.Object, newRequest.Name)
	w = httputil.SetHeaders(w)
	w.Write([]byte(json))
}

func run(w http.ResponseWriter, r *http.Request) {
	//get image name from request
	body := readBody(r)
	newRequest := jsonToRequest(body)

	controller.KubeRun(newRequest.Name)
	w = httputil.SetHeaders(w)
	w.Write([]byte("OK"))
}

func scale(w http.ResponseWriter, r *http.Request) {
	//get replicas and deployment filepath from request
	body := readBody(r)
	newRequest := jsonToRequest(body)

	controller.KubeScale(newRequest.Replicas, newRequest.Deployment)
	w = httputil.SetHeaders(w)
	w.Write([]byte("OK"))
}

func logs(w http.ResponseWriter, r *http.Request) {
	//get pod name from request
	body := readBody(r)
	newRequest := jsonToRequest(body)

	json := controller.KubeLogs(newRequest.Name)
	w = httputil.SetHeaders(w)
	w.Write([]byte(json))
}

//retrieves cluster-info string from cluster
func cluster(w http.ResponseWriter, r *http.Request) {
	json := controller.KubeClusterInfo()
	fmt.Println(json)
	w = httputil.SetHeaders(w)
	w.Write([]byte(json))
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
