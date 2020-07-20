package minecraft

import "github.com/jmoiron/sqlx"

var schema = `
CREATE TABLE IF NOT EXISTS minecraft (
	guild INT NOT NULL,
	host VARCHAR(50) NOT NULL,
	port INT UNSIGNED NOT NULL,
	pass VARCHAR(20) NOT NULL,
	PRIMARY_KEY(guild, host, port, pass),
	CONSTRAINT id
		FOREIGN KEY (guild)
		REFERENCES guild (id)
		ON DELETE NO ACTION
		ON UPDATE NO ACTION
);`

// Minecraft stores minecraft server settings for guilds
type Minecraft struct {
	Guild string `db:"guild"`
	Host  string `db:"host"`
	Port  uint   `db:"port"`
	Pass  string `db:"pass"`
}

// Init creates the Minecraft table in the database
func Init(db *sqlx.DB) {
	db.MustExec(schema)
}
