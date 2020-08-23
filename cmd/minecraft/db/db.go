package db

import (
	"database/sql"
	"errors"
	"log"
	"strings"

	"github.com/EmilyBjartskular/BlueBerryBunsBot/db"
)

var (
	// ErrInvID is an invalid id error
	ErrInvID = errors.New("Invalid ID")
)

var schema1 = `
CREATE TABLE IF NOT EXISTS minecraft (
	guild VARCHAR(25) NOT NULL,
	host VARCHAR(50) NOT NULL,
	port INT UNSIGNED NOT NULL,
	pass VARCHAR(20) NOT NULL,
	id INT,
	PRIMARY KEY(guild, host, port),
	FOREIGN KEY(guild) REFERENCES guild(id)
	ON DELETE CASCADE
	ON UPDATE CASCADE
);`

var schema2 = `
CREATE TABLE IF NOT EXISTS minecraft_prim (
	guild VARCHAR(25) NOT NULL,
	id INT,
	PRIMARY KEY (guild),
	FOREIGN KEY (guild) REFERENCES guild(id)
	ON DELETE CASCADE
	ON UPDATE CASCADE
)`

var schema3 = `
DROP PROCEDURE IF EXISTS UpdateServerIds;
DELIMITER $$
CREATE PROCEDURE
	UpdateServerIds
	(IN param_guild VARCHAR(25))
	MODIFIES SQL DATA
	BEGIN
		SET @i = 0;
		UPDATE minecraft SET id = @i:=@i+1 WHERE minecraft.guild = param_guild;
	END;
$$
DELIMITER ;
`

// Minecraft stores minecraft server settings for guilds
type Minecraft struct {
	Guild string `db:"guild"`
	Host  string `db:"host"`
	Port  uint   `db:"port"`
	Pass  string `db:"pass"`
	ID    int    `db:"id"`
}

// Init creates the Minecraft related tables and procedures.
func Init() {
	db.Connection.MustExec(schema1)
	db.Connection.MustExec(schema2)
	//db.Connection.MustExec(schema3)
}

// AddServer inserts a server for a guild to the database
func AddServer(guild, host, pass string, port int, prim string) error {
	res := db.Connection.QueryRow("SELECT COUNT(*) FROM minecraft WHERE guild=?", guild)
	var num int
	err := res.Scan(&num)
	if err != nil {
		log.Println(guild, err)
		return err
	}

	_, err = db.Connection.Exec("INSERT INTO minecraft (guild, host, port, pass, id) VALUES (?, ?, ?, ?, ?)", guild, host, port, pass, num)
	if err != nil {
		if !strings.Contains(err.Error(), "Error 1062") {
			log.Println(guild, err)
		}
		return err
	}

	if prim == "true" || num == 0 {
		_, err = db.Connection.Exec("INSERT INTO minecraft_prim(guild, id) VALUES(?, ?) ON DUPLICATE KEY UPDATE id=?", guild, num, num)
		if err != nil {
			log.Println(guild, err)
			return err
		}
	}
	return nil
}

// GetServer retrieves a server from a guild
func GetServer(guild string, id int) (*Minecraft, error) {
	row := db.Connection.QueryRowx("SELECT * FROM minecraft WHERE guild=? AND id=?", guild, id)
	srv := &Minecraft{}
	err := row.StructScan(srv)
	if err != nil {
		log.Println(guild, err)
		return nil, err
	}

	return srv, nil
}

// GetServers retreives all servers for a guild from the database
func GetServers(guild string) ([]*Minecraft, int, error) {
	servers := []*Minecraft{}
	err := db.Connection.Select(&servers, "SELECT * FROM minecraft WHERE guild=? ORDER BY id ASC", guild)
	if err != nil {
		log.Println(guild, err)
		return nil, -1, err
	}

	res := db.Connection.QueryRow("SELECT id FROM minecraft_prim WHERE guild=?", guild)
	var num int
	err = res.Scan(&num)
	if err == sql.ErrNoRows {
		return servers, -1, nil
	}
	if err != nil {
		log.Println(guild, err)
		return nil, -1, err
	}
	return servers, num, err
}

// DeleteServer deletes a server for a guild from the database
func DeleteServer(guild string, id int) (*Minecraft, error) {
	res := db.Connection.QueryRowx("SELECT * FROM minecraft WHERE guild = ? AND id = ?", guild, id)

	_, err := db.Connection.Exec("DELETE FROM minecraft WHERE guild = ? AND id = ?", guild, id)
	if err != nil {
		log.Println(guild, err)
		return nil, err
	}

	_, err = db.Connection.Exec("CALL UpdateServerIds(?)", guild)
	if err != nil {
		log.Println(guild, err)
		return nil, err
	}

	_, err = db.Connection.Exec("UPDATE minecraft_prim SET id = IFNULL((SELECT id FROM minecraft WHERE guild = ? ORDER BY id LIMIT 1), -1) WHERE guild = ?", guild, guild)
	if err != nil {
		log.Println(guild, err)
		return nil, err
	}

	srv := &Minecraft{}
	err = res.StructScan(srv)
	if err != nil {
		log.Println(guild, err)
		return nil, err
	}

	return srv, nil
}

// SetPrimary sets server as primary / default for a guild
func SetPrimary(guild string, id int) (*Minecraft, error) {
	res := db.Connection.QueryRowx("SELECT COUNT(*) FROM minecraft WHERE guild=?", guild)

	var num int
	err := res.Scan(&num)
	if err != nil {
		log.Println(guild, err)
		return nil, err
	}

	if id < 0 || id >= num {
		return nil, ErrInvID
	}

	log.Println(guild, id)
	res = db.Connection.QueryRowx("SELECT * FROM minecraft WHERE guild = ? AND id = ?", guild, id)

	srv := &Minecraft{}
	err = res.StructScan(srv)
	if err != nil {
		log.Println(srv)
		log.Println(guild, err, "????")
		return nil, err
	}

	_, err = db.Connection.Exec("UPDATE minecraft_prim SET id=? WHERE guild=?", id, guild)
	if err != nil {
		log.Println(guild, err)
		return nil, err
	}

	return srv, nil
}
