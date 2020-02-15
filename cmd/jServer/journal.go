package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

//creates a database for a user's entries
func createJournalTable(db *sql.DB, un string) {

	//create and execute query to create table if it doesn't already exist
	query := fmt.Sprintf(`IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='%s' and xtype='U') 
		CREATE TABLE %s (id INT NOT NULL IDENTITY(1,1) PRIMARY KEY, date VARCHAR(255), entry VARCHAR(8000))`, un, un)

	statement, err := database.Prepare(query)
	genericErrHandler("error", err)

	_, err = statement.Exec()
	genericErrHandler("error", err)
}

//addJEntry adds the given entry into the user's database table
func addJEntry(w http.ResponseWriter, r *http.Request) {
	var journal Journal

	//read from body
	body, err := ioutil.ReadAll(r.Body)
	genericErrHandler("print", err)

	w = setHeaders(w)

	//only run if body exists
	if len(body) != 0 {

		//populate journal struct from request body
		err = json.Unmarshal(body, &journal)
		genericErrHandler("error", err)

		//creates the user's journal table if it doesn't exist
		createJournalTable(database, journal.Username)

		//add journal entry to db
		inputEntry(database, journal.Username, journal.Journal[0].Entry)
	} else {
		//send bad request if user request doesn't have a body
		badRequest(w, "Empty body sent")
	}
}

// inputEntry adds the current date as a string adds a journal entry
// to be stored into the database in association with the date.
func inputEntry(db *sql.DB, un string, entry string) {

	//get today's date to add into journal entry
	journalDate := string(time.Now().Format("2006-01-02"))

	//create and execute query to add entry into user's database
	query := fmt.Sprintf("INSERT INTO %s (date, entry) VALUES ('%s', '%s')", un, journalDate, entry)
	statement, err := database.Prepare(query)
	genericErrHandler("error", err)

	_, err = statement.Exec()
	genericErrHandler("error", err)
}

//getJournalEntries send all journal entries from the logged in user's journal
func getJournalEntries(w http.ResponseWriter, r *http.Request) {
	var journal Journal

	//read from request body
	body, err := ioutil.ReadAll(r.Body)
	genericErrHandler("print", err)

	w = setHeaders(w)

	//only run if body exists
	if len(body) != 0 {

		//populate journal struct from request body
		err = json.Unmarshal(body, &journal)
		genericErrHandler("error", err)

		//create a journal table for this user in case it doesn't already exist
		createJournalTable(database, journal.Username)

		//get all journal entries sorted by date
		query := fmt.Sprintf("SELECT * FROM %s ORDER BY date DESC", journal.Username)
		rows, err := database.Query(query)
		genericErrHandler("error", err)

		//create variables for storing journal entries from database
		jEntries := []JEntry{}
		var jTableID int
		var jTableDate string
		var jTableEntry string

		//for each row returned by query, append to jEntries slice
		for rows.Next() {
			rows.Scan(&jTableID, &jTableDate, &jTableEntry)
			jEntries = append(jEntries, JEntry{Date: jTableDate, Entry: jTableEntry})
		}

		//add the jEntries slice into the journal struct
		journal.Journal = jEntries

		//unmarshall journal into stringified json and send as response
		j, err := json.Marshal(journal)
		genericErrHandler("error", err)
		w.Write(j)

	} else {
		//send bad request if user request doesn't have a body
		badRequest(w, "Empty body sent")
	}
}
