package minecraft

import (
	"github.com/EmilyBjartskular/BlueBerryBunsBot/db"
)

var schema = `
CREATE TABLE IF NOT EXISTS minecraft (
	guild INT NOT NULL,
	host VARCHAR(50) NOT NULL,
	port INT UNSIGNED NOT NULL,
	pass VARCHAR(20) NOT NULL,
	prim BOOL DEFAULT FALSE,
	PRIMARY_KEY(guild, host, port, pass),
	CONSTRAINT id
		FOREIGN KEY (guild)
		REFERENCES guild (id)
		ON DELETE NO ACTION
		ON UPDATE NO ACTION,
	CONSTRAINT primary UNIQUE(guild, prim)
);`

// Minecraft stores minecraft server settings for guilds
type Minecraft struct {
	Guild string `db:"guild"`
	Host  string `db:"host"`
	Port  uint   `db:"port"`
	Pass  string `db:"pass"`
	Prim  bool   `db:"prim"`
}

// Init creates the Minecraft table in the database
func Init() {
	db.Connection.MustExec(schema)
}

func addServer(guild, host, pass string, port int) {
	db.Connection.MustExec("INSERT INTO minecraft (guild, host, port, pass) VALUES($1, $2, $3, $4)", guild, host, port, pass)
}

func setPrimary(guild string) {
}
