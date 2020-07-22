package db

import (
	"fmt"
	"log"

	// Mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var schema = `
CREATE TABLE IF NOT EXISTS guild (
	id VARCHAR(25) NOT NULL,
	PRIMARY KEY(id)
);`

// Connection stores a global connection to the database
var Connection *sqlx.DB = nil

// Guild stores all guilds this bot is a member of
type Guild struct {
	ID int `db:"id"`
}

// Init opens a connection to the database and initializes standard tables.
func Init(user, pass, host string, port int64, dbname string) {
	var err error

	cn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, pass, host, port, dbname)
	Connection, err = sqlx.Connect("mysql", cn)

	if err != nil {
		log.Fatalln(err)
	}

	Connection.MustExec(schema)
}
