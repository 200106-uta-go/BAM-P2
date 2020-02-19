package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/200106-uta-go/BAM-P2/pkg/httputil"
	"golang.org/x/crypto/bcrypt"
)

// hashPassword generates a hashed and salted version of the plain text user password
// to store in the database
func hashPassword(pass string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	httputil.GenericErrHandler("error", err)
	return string(hash)
}

// userLogin verifies that the username and password given by the client
// are present in the database
func userLogin(w http.ResponseWriter, r *http.Request) {
	var user User

	//read request body
	body, err := ioutil.ReadAll(r.Body)
	httputil.GenericErrHandler("print", err)

	w = httputil.SetHeaders(w)

	//run only if request body exists
	if len(body) != 0 {

		//poopulate a User struct from the request body
		err = json.Unmarshal(body, &user)
		httputil.GenericErrHandler("error", err)

		//if the password is matched, user is authorized
		//otherwise send a bad request back to client
		if checkPassword(database, user.Username, user.Password) {
			response, err := json.Marshal(user)
			httputil.GenericErrHandler("error", err)
			w.Write(response)
		} else {
			httputil.BadRequest(w, "Invalid Authentication")
		}
	} else {
		//send bad request to client if no username/password was given
		httputil.BadRequest(w, "Empty body sent for login")
	}
}

// createUser creates a database entry with the given user credentials
func createUser(w http.ResponseWriter, r *http.Request) {
	var user User

	//read request body
	body, err := ioutil.ReadAll(r.Body)
	httputil.GenericErrHandler("print", err)

	//only run if body exists
	if len(body) != 0 {

		//poopulate a User struct from the request body
		err = json.Unmarshal(body, &user)
		httputil.GenericErrHandler("error", err)

		w = httputil.SetHeaders(w)

		addUser(database, user.Username, user.Password)

		//marshal the user into a byte slice to send back to client
		response, err := json.Marshal(user)
		httputil.GenericErrHandler("error", err)
		w.Write(response)
	} else {
		//send bad request to client if no username/password was given
		httputil.BadRequest(w, "Empty fields sent")
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
	httputil.GenericErrHandler("error", err)
	defer rows.Close()

	//for each row returned from query, check if password matches
	for rows.Next() {
		err := rows.Scan(&uTableID, &uTableUN, &uTablePS)
		httputil.GenericErrHandler("error", err)

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
	httputil.GenericErrHandler("error", err)

	_, err = statement.Exec()
	if err != nil {
		log.Println("User might already exist", err)
	}
}
