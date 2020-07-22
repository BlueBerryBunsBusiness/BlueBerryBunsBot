package minecraft

import (
	"log"

	"github.com/EmilyBjartskular/BlueBerryBunsBot/db"
)

var schema = `
CREATE TABLE IF NOT EXISTS minecraft (
	guild VARCHAR(25) NOT NULL,
	host VARCHAR(50) NOT NULL,
	port INT UNSIGNED NOT NULL,
	pass VARCHAR(20) NOT NULL,
	prim BOOL DEFAULT FALSE,
	PRIMARY KEY(guild, host, port, pass),
	UNIQUE(guild, prim)
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
	stmt, err := db.Connection.Prepare("INSERT INTO minecraft (guild, host, port, pass) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Println(err)
	}
	_, err = stmt.Exec(guild, host, port, pass)
	if err != nil {
		log.Println(err)
	}
}

func getServers(guild string) []Minecraft {
	servers := []Minecraft{}
	db.Connection.Select(&servers, "SELECT * FROM minecraft WHERE guild=? ORDER BY host ASC", guild)
	return servers
}

func setPrimary(guild string) {
}
