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
	PRIMARY_KEY(id)
);`

// Guild stores all guilds this bot is a member of
type Guild struct {
	ID int `db:"id"`
}

func New(dbname, user, pass string, port uint) *sqlx.DB {
	cn := fmt.Sprintf("dbname=%s user=%s pass=%s port=%d sslmode=enable", dbname, user, pass, port)
	db, err := sqlx.Connect("mysql", cn)
	if err != nil {
		log.Fatalln(err)
	}

	db.MustExec(schema)
	return db
}
