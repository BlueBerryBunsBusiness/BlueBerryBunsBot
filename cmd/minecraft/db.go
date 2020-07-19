package minecraft

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

type Minecraft struct {
	Guild	string	`db:"guild"`
	Host	string  `db:"host"`
	Port 	uint 		`db:"port"`
	Pass 	string 	`db:"pass"`
}
