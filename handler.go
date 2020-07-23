package main

import (
	"github.com/EmilyBjartskular/BlueBerryBunsBot/db"
	"github.com/bwmarrin/discordgo"
)

func guildCreate(_ *discordgo.Session, m *discordgo.GuildCreate) {
	db.AddGuild(m.Guild.ID)
}

func addRouter(_ *discordgo.Session, m *discordgo.MessageCreate) {
	Router.FindAndExecute(Session, Prefix, Session.State.User.ID, m.Message)
}
