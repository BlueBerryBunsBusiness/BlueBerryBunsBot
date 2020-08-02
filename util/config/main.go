package config

import (
	"fmt"
	"log"

	"github.com/pelletier/go-toml"
)

var (
	configPath = "config.toml"
	// Discord contains discordbot settings
	Discord = &confDiscord{}
	// Emoji contains bot emoji settings
	Emoji = &confEmoji{}
	// Database contains database settings
	Database = &confDatabase{}
	// Minecraft contains minecraft settings
	Minecraft = &confMinecraft{}
)

type confDiscord struct {
	Prefix string
	Token  string
}

type confEmoji struct {
	Guild string
}

type confDatabase struct {
	Name string
	Host string
	User string
	Pass string
	Port int64
}

type confMinecraft struct {
	Host string
	Pass string
	Port int64
}

func init() {
	config, err := toml.LoadFile("config.toml")
	if err != nil {
		log.Fatalln(err)
	}

	// DISCORD SECTION
	Discord.Prefix = getProperty("Discord.prefix", config).(string)
	Discord.Token = getProperty("Discord.token", config).(string)
	// EMOJI SECTION
	Emoji.Guild = getProperty("Emoji.srcGuild", config).(string)
	// DATABASE SECTION
	Database.Name = getProperty("Database.name", config).(string)
	Database.Host = getProperty("Database.host", config).(string)
	Database.User = getProperty("Database.user", config).(string)
	Database.Pass = getProperty("Database.pass", config).(string)
	Database.Port = getProperty("Database.port", config).(int64)
	// MINECRAFT SECTION
	Minecraft.Host = getProperty("Minecraft.host", config).(string)
	Minecraft.Pass = getProperty("Minecraft.pass", config).(string)
	Minecraft.Port = getProperty("Minecraft.port", config).(int64)
}

func getProperty(prop string, t *toml.Tree) interface{} {
	if !t.Has(prop) {
		log.Fatalln(fmt.Sprintf("Missing `%s` in `%s`", prop, configPath))
	}
	return t.Get(prop)
}
