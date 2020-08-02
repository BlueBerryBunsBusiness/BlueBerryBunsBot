package main

import (
	"github.com/EmilyBjartskular/BlueBerryBunsBot/db"
	"github.com/EmilyBjartskular/BlueBerryBunsBot/util/config"
	"github.com/EmilyBjartskular/BlueBerryBunsBot/util/global"
	"github.com/bwmarrin/discordgo"
)

func guildCreate(_ *discordgo.Session, m *discordgo.GuildCreate) {
	db.AddGuild(m.Guild.ID)
}

func addRouter(s *discordgo.Session, m *discordgo.MessageCreate) {
	global.Router.FindAndExecute(s, config.Discord.Prefix, s.State.User.ID, m.Message)
}
