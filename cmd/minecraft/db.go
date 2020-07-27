package minecraft

import (
	"log"

	"github.com/EmilyBjartskular/BlueBerryBunsBot/db"
)

var schema1 = `
CREATE TABLE IF NOT EXISTS minecraft (
	guild VARCHAR(25) NOT NULL,
	host VARCHAR(50) NOT NULL,
	port INT UNSIGNED NOT NULL,
	pass VARCHAR(20) NOT NULL,
	id INT UNSIGNED,
	PRIMARY KEY(guild, host, port),
	FOREIGN KEY(guild) REFERENCES guild(id)
	ON DELETE CASCADE
	ON UPDATE CASCADE
);`

var schema2 = `
CREATE TABLE IF NOT EXISTS minecraft_prim (
	guild VARCHAR(25) NOT NULL,
	id INT UNSIGNED,
	PRIMARY KEY (guild),
	FOREIGN KEY (guild) REFERENCES guild(id)
	ON DELETE CASCADE
	ON UPDATE CASCADE
)`

var schema3 = `
DELIMITER //
CREATE OR REPLACE PROCEDURE
	UpdateServerIds
	(param_guild VARCHAR(25))
	MODIFIES SQL DATA
	BEGIN
		SET @i = -1;
		UPDATE minecraft SET id = @i:=@i+1 WHERE minecraft.guild = param_guild;
	END;
//
DELIMITER ;
`

// Minecraft stores minecraft server settings for guilds
type Minecraft struct {
	Guild string `db:"guild"`
	Host  string `db:"host"`
	Port  uint   `db:"port"`
	Pass  string `db:"pass"`
	id    uint   `db:"id"`
}

// Init creates the Minecraft related tables and procedures.
func Init() {
	db.Connection.MustExec(schema1)
	db.Connection.MustExec(schema2)
	db.Connection.MustExec(schema3)
}

// addServer inserts a server for a guild to the database
func addServer(guild, host, pass string, port int, prim string) error {
	res := db.Connection.QueryRow("SELECT COUNT(*) FROM minecraft WHERE guild=?", guild)
	var num int
	err := res.Scan(&num)
	if err != nil {
		log.Println(guild, err)
		return err
	}

	tx, err := db.Connection.Begin()
	if err != nil {
		log.Println(guild, err)
		return err
	}
	tx.Exec("INSERT INTO minecraft (guild, host, port, pass, id) VALUES (?, ?, ?, ?, ?)", guild, host, port, pass, num)
	if prim == "true" || num == 0 {
		tx.Exec("INSERT INTO minecraft_prim(guild, id) VALUES(?, ?) ON DUPLICATE KEY UPDATE id=?", guild, num, num)
	}
	err = tx.Commit()

	if err != nil {
		log.Println(guild, err)
		return err
	}
	return nil
}

// getServers retreives all servers for a guild from the database
func getServers(guild string) []Minecraft {
	servers := []Minecraft{}
	db.Connection.Select(&servers, "SELECT * FROM minecraft WHERE guild=? ORDER BY host ASC", guild)
	return servers
}

// deleteServer deletes a server for a guild from the database
func deleteServer(guild, string, id int) error {
	stmt, err := db.Connection.Prepare("DELETE FROM minecraft WHERE guild = ? AND id = ?")
	if err != nil {
		log.Println(guild, err)
		return err
	}
	_, err = stmt.Exec(guild, id)
	if err != nil {
		log.Println(guild, err)
		return err
	}
	return nil
}

// setPrimary sets server as primary / default for a guild
func setPrimary(guild string) {

}
