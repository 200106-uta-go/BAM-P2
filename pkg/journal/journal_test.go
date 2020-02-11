package journal

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestViewEntireJournal(t *testing.T) {
	database, err := sql.Open("sqlite3", "./journal_test.db")
	if err != nil {
		log.Fatal(err)
	}

	statement, err := database.Prepare("DROP TABLE IF EXISTS journal_entries")
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()

	statement, err = database.Prepare("CREATE TABLE IF NOT EXISTS journal_entries (id INTEGER PRIMARY KEY, date TEXT, entry TEXT)")
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()

	statement, err = database.Prepare(`INSERT INTO journal_entries (date, entry) VALUES ("01-01-2020", "Today is New Year's Day!")`)
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()

	statement, err = database.Prepare(`INSERT INTO journal_entries (date, entry) VALUES ("01-02-2020", "It's the day after New Year's Day!")`)
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()

	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	ViewEntireJournal(database)

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	if string(out) != "\nView All Entries\n\n01-01-2020:\nToday is New Year's Day!\n01-02-2020:\nIt's the day after New Year's Day!\n" {
		t.Errorf("Got:\n%v\nExpected:\nView All Entries\n\n01-01-2020:\nToday is New Year's Day!\n01-02-2020:\nIt's the day after New Year's Day!\n", string(out))
	}
}

func TestPrintEntry(t *testing.T) {
	database, err := sql.Open("sqlite3", "./journal_test.db")
	if err != nil {
		log.Fatal(err)
	}

	statement, err := database.Prepare("DROP TABLE IF EXISTS journal_entries")
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()

	statement, err = database.Prepare("CREATE TABLE IF NOT EXISTS journal_entries (id INTEGER PRIMARY KEY, date TEXT, entry TEXT)")
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()

	statement, err = database.Prepare(`INSERT INTO journal_entries (date, entry) VALUES ("01-01-2020", "Today is New Year's Day!")`)
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()

	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printEntry(database, "01-01-2020")

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	if string(out) != "\n01-01-2020:\nToday is New Year's Day!\n" {
		t.Errorf("Got:\n%v\nExpected:\n01-01-2020:\nToday is New Year's Day!\n\n", string(out))
	}
}
