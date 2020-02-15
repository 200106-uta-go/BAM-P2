package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/joho/godotenv"
)

var database *sql.DB

//set up the database connection
func init() {

	//load in environment variables from .env
	err := godotenv.Load("../.env")
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
		log.Fatalln(err)
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("WARNING...this program will delete all tables from the database!")
	fmt.Print("Confirm you would like to delete all tables [y/n]: ")
	byte, err := reader.ReadByte()
	if err != nil {
		log.Fatalln(err)
	}

	if string(byte) == "y" {
		//delete all db tables
		deleteAllTables(database)
		fmt.Println("All tables in database have been deleted")

	} else {
		fmt.Println("Exiting. . . Database has not been deleted.")
	}
}

// deleteAllTables deletes all tables from the database
// this is only used for testing/development
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
		log.Fatalln(err)
	}
	defer statement.Close()
	statement.Exec()
}
