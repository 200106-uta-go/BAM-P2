package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/200106-uta-go/BAM-P2/pkg/httputil"
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

//set up the database connection
func init() {

	//load in environment variables from .env
	//will print error message when running from docker image
	//because env file is passed into docker run command
	envErr := godotenv.Load("/home/ubuntu/go/src/github.com/200106-uta-go/BAM-P2/.env")
	if envErr != nil {
		if !strings.Contains(envErr.Error(), "no such file or directory") {
			log.Println("Error loading .env: ", envErr)
		}
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
	httputil.GenericErrHandler("error", err)

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

	//create endpoints for web client api
	goodRequest := http.HandlerFunc(httputil.GoodRequest)
	http.Handle("/", logRequest(goodRequest))

	login := http.HandlerFunc(userLogin)
	http.Handle("/login", logRequest(login))

	create := http.HandlerFunc(createUser)
	http.Handle("/createUser", logRequest(create))

	addEntry := http.HandlerFunc(addJEntry)
	http.Handle("/addJEntry", logRequest(addEntry))

	getJournal := http.HandlerFunc(getJournalEntries)
	http.Handle("/getJournal", logRequest(getJournal))

	fmt.Printf("HTTP server listening on port %s\n", servPort)
	http.ListenAndServe(servPort, nil)
}

//middleware to send all http requests to the logger
func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
