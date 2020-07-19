package db

//"database/sql"

//_ "github.com/go-sql-driver/mysql"
//"github.com/jmoiron/sqlx"

var schema = `
CREATE TABLE IF NOT EXISTS guild (
	id VARCHAR(25) NOT NULL,
	PRIMARY_KEY(id)
);`

type Guild struct {
	ID int `db:"id"`
}
