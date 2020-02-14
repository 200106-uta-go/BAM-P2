package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/joho/godotenv"
)

var (
	database *sql.DB // Pointer to database handle
	err      error   // Temporary reference to error value
	uTableID int     // Temporary reference to id value of table user_table
	uTableUN string  // Temporary reference to username value of table user_table
	uTablePS string  // Temporary reference to password value of table user_table
)

type user struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

//set up the database connection
func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalln("Error loading .env: ", err)
	}

	var server = os.Getenv("DB_SERVER")
	var port = os.Getenv("DB_PORT")
	var dbUser = os.Getenv("DB_USER")
	var dbPass = os.Getenv("DB_PASS")
	var db = os.Getenv("DB_NAME")

	// Build connection string
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;", server, dbUser, dbPass, port, db)

	// Create connection pool
	database, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	ctx := context.Background()
	err = database.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {

	fs := http.FileServer(http.Dir("../../web"))
	http.Handle("/", fs)
	http.HandleFunc("/login", userLogin)
	http.HandleFunc("/createUser", createUser)

	http.ListenAndServe(":8080", nil)
}

func userLogin(w http.ResponseWriter, r *http.Request) {
	user := parseBody(r)
	w = setHeaders(w)
	checkPassword(database, user.Name, user.Password)
	fmt.Fprintln(w, "OK")
}

func createUser(w http.ResponseWriter, r *http.Request) {
	user := parseBody(r)
	w = setHeaders(w)
	addUser(database, user.Name, user.Password)
	fmt.Fprintln(w, "OK")
}

func parseBody(r *http.Request) user {
	thisUser := user{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	json.Unmarshal(body, &thisUser)

	return thisUser
}

// checkPassword checks for username and password congruency
func checkPassword(db *sql.DB, un string, ps string) bool {
	query := fmt.Sprintf("SELECT * FROM %s WHERE username = %s", un, un)
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&uTableID, &uTableUN, &uTablePS)
		if err != nil {
			log.Fatal(err)
		}
		if uTablePS == ps {
			return true
		}
	}
	return false
}

// addUser adds the username to user_table
func addUser(db *sql.DB, un string, ps string) {
	query := fmt.Sprintf("INSERT INTO %s (username, password) VALUES (%s, %s)", un, un, ps)
	statement, err := database.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec(un, un, ps)
}

func setHeaders(w http.ResponseWriter) http.ResponseWriter {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
	return w
}
