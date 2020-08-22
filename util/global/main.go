package global

import (
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
)

var (
	// Route is the global command router
	Router *exrouter.Route
	// Session is the global discord session
	Session *discordgo.Session
	// ConfigPath is the path to the configuration file
	ConfigPath = "configs/config.toml"
	// TextConfigPath is the path to the test configuration file
	TextConfigPath = "configs/text.toml"
)
