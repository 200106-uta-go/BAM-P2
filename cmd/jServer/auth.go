package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// hashPassword generates a hashed and salted version of the plain text user password
// to store in the database
func hashPassword(pass string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	genericErrHandler("error", err)
	return string(hash)
}

// userLogin verifies that the username and password given by the client
// are present in the database
func userLogin(w http.ResponseWriter, r *http.Request) {
	var user User

	//read request body
	body, err := ioutil.ReadAll(r.Body)
	genericErrHandler("print", err)

	w = setHeaders(w)

	//run only if request body exists
	if len(body) != 0 {

		//poopulate a User struct from the request body
		err = json.Unmarshal(body, &user)
		genericErrHandler("error", err)

		//if the password is matched, user is authorized
		//otherwise send a bad request back to client
		if checkPassword(database, user.Username, user.Password) {
			response, err := json.Marshal(user)
			genericErrHandler("error", err)
			w.Write(response)
		} else {
			badRequest(w, "Invalid Authentication")
		}
	} else {
		//send bad request to client if no username/password was given
		badRequest(w, "Empty body sent for login")
	}
}

// createUser creates a database entry with the given user credentials
func createUser(w http.ResponseWriter, r *http.Request) {
	var user User

	//read request body
	body, err := ioutil.ReadAll(r.Body)
	genericErrHandler("print", err)

	//only run if body exists
	if len(body) != 0 {

		//poopulate a User struct from the request body
		err = json.Unmarshal(body, &user)
		genericErrHandler("error", err)

		w = setHeaders(w)

		addUser(database, user.Username, user.Password)

		//marshal the user into a byte slice to send back to client
		response, err := json.Marshal(user)
		genericErrHandler("error", err)
		w.Write(response)
	} else {
		//send bad request to client if no username/password was given
		badRequest(w, "Empty fields sent")
	}
}

// checkPassword checks for username and password congruency
func checkPassword(db *sql.DB, un string, ps string) bool {
	//set up variables to hold each row's contents
	var uTableID int
	var uTableUN string
	var uTablePS string

	//create and execute query on database to find user
	query := fmt.Sprintf("SELECT * FROM user_table WHERE username = '%s'", un)
	rows, err := db.Query(query)
	genericErrHandler("error", err)
	defer rows.Close()

	//for each row returned from query, check if password matches
	for rows.Next() {
		err := rows.Scan(&uTableID, &uTableUN, &uTablePS)
		genericErrHandler("error", err)

		//compare database's hashed password with client's plain password
		hashErr := bcrypt.CompareHashAndPassword([]byte(uTablePS), []byte(ps))
		return hashErr == nil
	}

	//return false if password could not be matched
	return false
}

// addUser adds the username to user_table
func addUser(db *sql.DB, un string, ps string) {
	//creates a hash from password to save into database
	hash := hashPassword(ps)

	//create query to add user
	query := fmt.Sprintf("INSERT INTO user_table (username, password) VALUES ('%s', '%s')", un, hash)
	statement, err := database.Prepare(query)
	genericErrHandler("error", err)

	_, err = statement.Exec()
	if err != nil {
		log.Println("User might already exist", err)
	}
}
