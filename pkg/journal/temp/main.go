package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
)

type JEntry struct {
	Date  string `json:"date"`
	Entry string `json:"entry"`
}

var (
	username    string   // Username of journal
	password    string   // Password of journal
	database    *sql.DB  // Pointer to database handle
	err         error    // Temporary reference to error value
	uTableID    int      // Temporary reference to id value of table user_table
	uTableUN    string   // Temporary reference to username value of table user_table
	uTablePS    string   // Temporary reference to password value of table user_table
	jTableID    int      // Temporary reference to id value of table journal_entries
	jTableDate  string   // Temporary reference to date value of table journal_entries
	jTableEntry string   // Temporary reference to entry value of table journal_entries
	JEntries    []JEntry // Slice of JEntry struct
)

func init() {
	var server = "<your_server.database.windows.net>"
	var port = 1433
	var user = "<your_username>"
	var password = "<your_password>"
	var database = "<your_database>"

	// Build connection string
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;", server, user, password, port, database)

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

// // Tentative plan for user authorization
// func main() {
// 	/*
// 		Implement acquisition methods for username and password
// 	*/
// 	usercheck := CheckUsername(database, username)

// 	if usercheck {
// 		for {
// 			passwordcheck := CheckPassword(database, username, password)

// 			if passwordcheck {
// 				/*
// 					Implementation of main functions
// 				*/
// 			} else {
// 				/*
// 					Implement condition if password is invalid
// 				*/
// 			}
// 		}
// 	} else {
// 		for {
// 			/*
// 				Implement prompt and input method of new username
// 			*/
// 			checkUsernameChars := CheckInvalidChars(username)
// 			if checkUsernameChars {
// 				for {
// 					checkPasswordChars := CheckInvalidChars(password)
// 					if checkPasswordChars {
// 						AddUsernamePassword(database, username, password)
// 						break
// 					} else {
// 						/*
// 							Implement condition if password is invalid
// 						*/
// 					}
// 				}
// 				break
// 			} else {
// 				/*
// 					Implement condition if username is invalid
// 				*/
// 			}
// 		}
// 	}

// 	// Creates the username table if it does not exist
// 	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS ? (id INTEGER PRIMARY KEY, username TEXT, password TEXT)")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	statement.Exec(username)
// }

// // CheckUsername checks for existence of username in user_table
// func CheckUsername(db *sql.DB, un string) bool {
// 	rows, err := db.Query(`SELECT * FROM ?`, un)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		err := rows.Scan(&jTableID, &jTableDate, &jTableEntry)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		if uTableUN == un {
// 			return true
// 		}
// 	}
// }

// CheckInvalidUNChars checks username for the existence of special characters
func CheckInvalidUNChars(un string) bool {
	if strings.ContainsAny(un, " !@#$%^&()[]{}`~:;<>,./\\+*\"?'") == false {
		return false
	}
	return true
}

// CheckInvalidPSChars checks password for the existence of special characters
func CheckInvalidPSChars(un string) bool {
	if strings.ContainsAny(un, " ()[]{}~:;<>,./\\+\"'") == false {
		return false
	}
	return true
}

// AddUsernamePassword adds the username to user_table
func AddUsernamePassword(db *sql.DB, un string, ps string) {
	statement, err := database.Prepare(`INSERT INTO ? (username, password) VALUES (?, ?)`)
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec(un, un, ps)
}

// CheckPassword checks for username and password congruency
func CheckPassword(db *sql.DB, un string, ps string) bool {
	rows, err := db.Query(`SELECT * FROM ? WHERE username = ?`, un, un)
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

// EditEntry replaces the journal entry of a particular date
func EditEntry(db *sql.DB, un string, jDate string, jEntry string) {
	statement, err := db.Prepare("UPDATE ? SET entry = ? WHERE date = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer statement.Close()
	statement.Exec(un, jEntry, jDate)
}

// IfEntryExists checks to see if an entry for a certain date already exists
// in the username table in a specified SQL database.
func IfEntryExists(db *sql.DB, un string, jEntry string, jDate string) {
	rows, err := db.Query(`SELECT * FROM ? WHERE date = ?`, un, jDate)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&jTableID, &jTableDate, &jTableEntry)
		if err != nil {
			log.Fatal(err)
		}
		// If the date of the entry already exists, the entry will be appended	to
		// the preexisting entry after a new line.
		if jTableDate == jDate {
			rows, err = db.Query("SELECT * FROM ? WHERE date = ?", un, jDate)
			if err != nil {
				log.Fatal(err)
			}
			for rows.Next() {
				err := rows.Scan(&jTableID, &jTableDate, &jTableEntry)
				if err != nil {
					log.Fatal(err)
				}
				jEntry = fmt.Sprint(jTableEntry + "\n\n" + jEntry)
			}

			EditEntry(db, un, jDate, jEntry)

			// If the date of the entry does not exist, the entry will be added to
			// the added to the table.
		} else {
			statement, err := db.Prepare("INSERT INTO ? (date, entry) VALUES (?, ?)")
			if err != nil {
				log.Fatal(err)
			}
			defer statement.Close()
			statement.Exec(un, jDate, jEntry)
		}
	}
	return
}

// OutputJournal returns entire journal of username in []jEntry format
func OutputJournal(un string) []JEntry {
	rows, err := database.Query("SELECT * FROM ? ORDER BY date DESC", un)
	if err != nil {
		log.Fatal(err)
	}

	JEntries = []JEntry{}

	for rows.Next() {
		rows.Scan(&jTableID, &jTableDate, &jTableEntry)
		JEntries = append(JEntries, JEntry{Date: jTableDate, Entry: jTableEntry})
	}

	return JEntries
}
