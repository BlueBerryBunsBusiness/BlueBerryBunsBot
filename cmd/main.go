package cmd

import (
	"fmt"
	"log"

	"github.com/EmilyBjartskular/BlueBerryBunsBot/util/config"
	"github.com/EmilyBjartskular/BlueBerryBunsBot/util/text"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
)

var (
	helpEmbed *discordgo.MessageEmbed
)

// TODO: Rewamp b!help and b!help <command> as well as calling empty commands

// Add adds set and get for testing management of context variables
func Add(r *exrouter.Route) {
	help := r.On("help", helpFunc).Desc("Prints this help menu").Cat("Help")
	r.Default = help

	for _, v := range r.Routes {
		if v.Name != help.Name {
			log.Println(v.Category)
		}
	}
}

func helpFunc(ctx *exrouter.Context) {
	guild := ctx.Msg.GuildID
	channel := ctx.Msg.ChannelID
	// Initialize empty embed
	embed := &discordgo.MessageEmbed{}

	// Create embedfields slize
	fields := make([]*discordgo.MessageEmbedField, len(ctx.Route.Parent.Routes)-1)

	embed.Title = text.Generic.EmbedTitle

	desc := text.Generic.EmbedDescription
	usage := fmt.Sprintf("```\n%s%s <command>\n```", config.Discord.Prefix, ctx.Route.Name)
	embed.Description = fmt.Sprintf("%s\n%s", desc, usage)

	// Create inline embed fields
	i := 0
	for _, v := range ctx.Route.Parent.Routes {
		if v.Name != ctx.Route.Name {
			field := &discordgo.MessageEmbedField{
				Name:   fmt.Sprintf("%s", v.Category),
				Value:  fmt.Sprintf("`%s`", v.Name),
				Inline: true,
			}
			fields[i] = field
			i++
		}
	}
	embed.Fields = fields
	_, err := ctx.Ses.ChannelMessageSendEmbed(channel, embed)
	if err != nil {
		log.Println(guild, err)
	}
}
