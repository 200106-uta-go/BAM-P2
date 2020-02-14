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
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/joho/godotenv"
)

var (
	username    string  // Username of journal
	password    string  // Password of journal
	database    *sql.DB // Pointer to database handle
	err         error   // Temporary reference to error value
	uTableID    int     // Temporary reference to id value of table user_table
	uTableUN    string  // Temporary reference to username value of table user_table
	uTablePS    string  // Temporary reference to password value of table user_table
	jTableID    int     // Temporary reference to id value of table journal_entries
	jTableDate  string  // Temporary reference to date value of table journal_entries
	jTableEntry string  // Temporary reference to entry value of table journal_entries

	// JEntries is a slice of JEntry struct for a json file
	JEntries []JEntry
)

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
	err := godotenv.Load("/home/ubuntu/go/src/github.com/200106-uta-go/BAM-P2/.env")
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

	//delete all db tables
	deleteAllTables(database)

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

	fs := http.FileServer(http.Dir("../../web"))
	http.Handle("/", fs)
	http.HandleFunc("/login", userLogin)
	http.HandleFunc("/createUser", createUser)
	http.HandleFunc("/addJEntry", addJEntry)
	http.HandleFunc("/getJournal", getJournalEntries)

	fmt.Println("Server listening at localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func userLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("userLogin running")
	var user User
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	w = setHeaders(w)
	if len(body) != 0 {
		fmt.Println(string(body))
		err = json.Unmarshal(body, &user)
		if err != nil {
			log.Fatalln(err)
		}
		if checkPassword(database, user.Username, user.Password) {
			fmt.Println("User Authorized")
			response, err := json.Marshal(user)
			genericErrHandler("error", err)
			fmt.Println(string(response))
			w.Write(response)
		} else {
			badRequest(w, "Invalid Authentication")
		}
	} else {
		badRequest(w, "Empty body sent for login")
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("createUser")
	var user User
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	if len(body) != 0 {
		fmt.Println(string(body))
		err = json.Unmarshal(body, &user)
		if err != nil {
			log.Fatalln(err)
		}
		w = setHeaders(w)
		addUser(database, user.Username, user.Password)
		response, err := json.Marshal(user)
		genericErrHandler("error", err)
		fmt.Fprintln(w, response)
	} else {
		fmt.Fprintln(w, "{message: 'Empty fields sent'}")
	}
}

// checkPassword checks for username and password congruency
func checkPassword(db *sql.DB, un string, ps string) bool {
	fmt.Println("checkPassword running")
	query := fmt.Sprintf("SELECT * FROM user_table WHERE username = '%s'", un)
	fmt.Println(query)
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
		fmt.Println(uTablePS, ps)
		if uTablePS == ps {
			fmt.Println("user logged in")
			return true
		}
	}
	fmt.Println("user not logged in")
	return false
}

// addUser adds the username to user_table
func addUser(db *sql.DB, un string, ps string) {
	fmt.Println("addUser running")
	fmt.Println(un, ps)
	query := fmt.Sprintf("INSERT INTO user_table (username, password) VALUES ('%s', '%s')", un, ps)
	statement, err := database.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	_, err = statement.Exec()
	if err != nil {
		log.Println("User might already exist", err)
	}
}

func addJEntry(w http.ResponseWriter, r *http.Request) {
	var journal Journal
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	w = setHeaders(w)
	if len(body) != 0 {
		fmt.Println(string(body))
		err = json.Unmarshal(body, &journal)
		genericErrHandler("error", err)

		//creates the user's journal table if it doesn't exist
		fmt.Println("Creating journal table")
		createJournalTable(database, journal.Username)

		//add journal entry to db
		fmt.Println("adding journal entry", journal)
		inputEntry(database, journal.Username, journal.Journal[0].Entry)
	} else {
		badRequest(w, "Empty body sent")
	}
}

// inputEntry adds the current date as a string adds a journal entry
// to be stored into the database in association with the date.
func inputEntry(db *sql.DB, un string, entry string) {
	fmt.Print("\nAdd An Entry For Today\n\n")

	journalDate := string(time.Now().Format("2006-01-02"))
	query := fmt.Sprintf("INSERT INTO %s (date, entry) VALUES ('%s', '%s')", un, journalDate, entry)
	statement, err := database.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	_, err = statement.Exec()
	if err != nil {
		log.Fatalln(err)
	}
}

func getJournalEntries(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getJournalEntries running")
	var journal Journal
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	w = setHeaders(w)

	if len(body) != 0 {
		fmt.Println(string(body))
		err = json.Unmarshal(body, &journal)
		genericErrHandler("error", err)
		w = setHeaders(w)

		//create a journal table for this user in case it doesn't already exist
		createJournalTable(database, journal.Username)

		//get latest journal entry
		query := fmt.Sprintf("SELECT * FROM %s ORDER BY date DESC", journal.Username)

		fmt.Println(query)
		rows, err := database.Query(query)
		if err != nil {
			log.Fatal(err)
		}

		JEntries = []JEntry{}

		for rows.Next() {
			rows.Scan(&jTableID, &jTableDate, &jTableEntry)
			JEntries = append(JEntries, JEntry{Date: jTableDate, Entry: jTableEntry})
		}

		journal.Journal = JEntries

		//unmarshall journal into stringified json
		j, err := json.Marshal(journal)
		genericErrHandler("error", err)
		w.Write(j)

	} else {
		badRequest(w, "Empty body sent")
	}

}

func setHeaders(w http.ResponseWriter) http.ResponseWriter {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
	return w
}

func genericErrHandler(action string, err error) {
	if action == "error" {
		if err != nil {
			log.Fatalln(err)
		}
	} else if action == "print" {
		if err != nil {
			log.Println(err)
		}
	}
}

func deleteAllTables(db *sql.DB) {
	statement, err := db.Prepare(`while(exists(select 1 from INFORMATION_SCHEMA.TABLES 
        where TABLE_NAME != '__MigrationHistory' 
        AND TABLE_TYPE = 'BASE TABLE'))
        begin
        declare @sql nvarchar(2000)
        SELECT TOP 1 @sql=('DROP TABLE ' + TABLE_SCHEMA + '.[' + TABLE_NAME
        + ']')
        FROM INFORMATION_SCHEMA.TABLES
        WHERE TABLE_NAME != '__MigrationHistory' AND TABLE_TYPE = 'BASE TABLE'
        exec (@sql)
        end`)
	if err != nil {
		log.Fatal(err)
	}
	defer statement.Close()
	statement.Exec()
}

func badRequest(w http.ResponseWriter, message string) {
	response := HTTPResponse{
		Message: message,
	}
	r, err := json.Marshal(response)
	genericErrHandler("error", err)
	w.WriteHeader(http.StatusBadRequest)
	w.Write(r)
}

//creates a database for a user's entries
func createJournalTable(db *sql.DB, un string) {
	query := fmt.Sprintf(`IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='%s' and xtype='U') 
		CREATE TABLE %s (id INT NOT NULL IDENTITY(1,1) PRIMARY KEY, date VARCHAR(255), entry VARCHAR(8000))`, un, un)
	statement, err := database.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(query)
	_, err = statement.Exec()
	if err != nil {
		log.Fatalln(err)
	}
}
