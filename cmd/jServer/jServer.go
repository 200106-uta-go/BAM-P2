package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/joho/godotenv"
)

// Pointer to database handle
var database *sql.DB

// User struct holds user information sent in requests
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Journal holds a user's entire journal
type Journal struct {
	Username string   `json:"username"`
	Journal  []JEntry `json:"journal"`
}

// JEntry is the struct for setting up a json file
type JEntry struct {
	Date  string `json:"date"`
	Entry string `json:"entry"`
}

// HTTPResponse defines a generic struct for sending a http response message as JSON
type HTTPResponse struct {
	Message string `json:"message"`
}

//set up the database connection
func init() {

	//load in environment variables from .env
	//will print error message when running from docker image
	//because env file is passed into docker run command
	envErr := godotenv.Load("/home/ubuntu/go/src/github.com/200106-uta-go/BAM-P2/.env")
	if envErr != nil {
		log.Println("Error loading .env: ", envErr)
	}

	var server = os.Getenv("DB_SERVER")
	var dbPort = os.Getenv("DB_PORT")
	var dbUser = os.Getenv("DB_USER")
	var dbPass = os.Getenv("DB_PASS")
	var db = os.Getenv("DB_NAME")

	// Build connection string
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;", server, dbUser, dbPass, dbPort, db)

	// Create connection pool
	var err error
	database, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	ctx := context.Background()
	err = database.PingContext(ctx)
	genericErrHandler("error", err)

	//create user table if it doesn't exist
	statement, err := database.Prepare(`IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='user_table' and xtype='U') 
		CREATE TABLE user_table (id INT NOT NULL IDENTITY(1,1) PRIMARY KEY, username VARCHAR(255), password VARCHAR(255))`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = statement.Exec()
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {

	servPort := ":" + os.Getenv("SERV_PORT")

	//set up file server to serve html
	fs := http.FileServer(http.Dir("../../web"))

	//create endpoints for web client
	http.Handle("/", fs)
	http.HandleFunc("/login", userLogin)
	http.HandleFunc("/createUser", createUser)
	http.HandleFunc("/addJEntry", addJEntry)
	http.HandleFunc("/getJournal", getJournalEntries)

	fmt.Printf("HTTP server listening on port %s\n", servPort)
	http.ListenAndServe(servPort, nil)
}

// setHeaders sets the response headers for an outgoing HTTP response
func setHeaders(w http.ResponseWriter) http.ResponseWriter {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
	return w
}

// genericErrHandler is a function to replace the common
// generic error handler written throughout the code
func genericErrHandler(action string, err error) {
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

// badRequest sends an HTTP response with a Bad Request status code along
// with the message passed into the function back to the client in JSON format
func badRequest(w http.ResponseWriter, message string) {
	response := HTTPResponse{
		Message: message,
	}
	r, err := json.Marshal(response)
	genericErrHandler("error", err)

	w.WriteHeader(http.StatusBadRequest)
	w.Write(r)
}
