package global

import (
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
)

var (
	// Route is the global command router
	Router  *exrouter.Route
	Session *discordgo.Session
)
