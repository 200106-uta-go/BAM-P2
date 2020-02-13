package main

import (
	"database/sql"
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
)

func init() {
	// **Need to alter database access for SQL Server or Postgres**
	// Makes a handle for the database journal
	database, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	/*
		Implement acquisition methods for username and password
	*/
	usercheck := CheckUsername(database, username)

	if usercheck {
		passwordcheck := CheckPassword(database, username, password)

		if passwordcheck {

		}
	} else {
		for {
			/*
				Implement prompt and input method of new username
			*/
			checkUsernameChars := CheckInvalidChars(username)
			if checkUsernameChars {
				for {
					/*
						Implement prompt and input method of new username
					*/
					checkPasswordChars := CheckInvalidChars(password)
					if checkPasswordChars {
						AddUsernamePassword(database, username, password)
						break
					} else {
						/*
							Implement condition if password is invalid
						*/
					}
				}
				break
			} else {
				/*
					Implement condition if username is invalid
				*/
			}
		}
	}

	// Creates the table user_table if it does not exist
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS ? (id INTEGER PRIMARY KEY, username TEXT, password TEXT)")
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec(username)
}

// CheckUsername checks for existence of username in user_table
func CheckUsername(db *sql.DB, un string) bool {
	rows, err := db.Query(`SELECT * FROM ?`, un)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&jTableID, &jTableDate, &jTableEntry)
		if err != nil {
			log.Fatal(err)
		}
		if uTableUN == un {
			return true
		}
	}
}

// CheckInvalidUNChars checks for the existence of special characters
func CheckInvalidUNChars(un string) bool {
	if strings.ContainsAny(un, " !@#$%^&()[]{}`~:;<>,./\\+*\"?'") == false {
		return false
	}
}

// CheckInvalidPSChars checks for the existence of special characters
func CheckInvalidPSChars(un string) bool {
	if strings.ContainsAny(un, " ()[]{}~:;<>,./\\+\"'") == false {
		return false
	}
}

// AddUsername adds the username to user_table
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
}

// InputEntry adds the current date as a string and prompts the user for
// a journal entry input to be stored into the database in association
// with the date.
func InputEntry(db *sql.DB, un string, jEntry string) {
	jDate := string(time.Now().Format("2006-01-02"))

	ifEntryExists(db, un, jEntry, jDate)

	printEntry(db, jDate)
}

// ifEntryExists checks to see if an entry for a certain date already exists
// in the username table in a specified SQL database.
func ifEntryExists(db *sql.DB, un string, jEntry string, jDate string) {
	rows, err := db.Query(`SELECT * FROM ? WHERE date = ?`, un, jDate)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	dateExists := false

	for rows.Next() {
		err := rows.Scan(&jTableID, &jTableDate, &jTableEntry)
		if err != nil {
			log.Fatal(err)
		}
		if jTableDate == jDate {
			dateExists = true
		}
	}

	// If the date of the entry already exists, the entry will be added	to
	// the preexisting entry after a new line.
	if dateExists {
		rows, err = db.Query("SELECT * FROM un WHERE date = ?", un, jDate)
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

		statement, err := db.Prepare("UPDATE un SET entry = ? WHERE date = ?")
		if err != nil {
			log.Fatal(err)
		}
		defer statement.Close()
		statement.Exec(un, jEntry, jDate)

	} else {
		statement, err := db.Prepare("INSERT INTO un (date, entry) VALUES (?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		defer statement.Close()
		statement.Exec(un, jDate, jEntry)
	}
}
