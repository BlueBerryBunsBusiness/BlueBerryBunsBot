package minecraft

import (
	"fmt"
	"log"
	"strconv"

	"github.com/EmilyBjartskular/BlueBerryBunsBot/cmd"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bearbin/mcgorcon"
	"github.com/bwmarrin/discordgo"
)

var _ = mcgorcon.Dial // For debugging; delete when done

var (
	// Help is the minecraft command help message
	Help *discordgo.MessageEmbed
)

func init() {
	Help.Author.IconURL = "https://i.imgur.com/ipJTwRA.png"
}

// Add adds the minecraft command package to a CommandController
func Add(c *cmd.CommandController) {
	c.Router.On("mc", func(ctx *exrouter.Context) {
		command := ctx.Args.Get(1)
		switch command {
		case "server":
			server(ctx)
		default:

		}
		ctx.ReplyEmbed(Help)
	}).Desc("Well, it's Minecraft y'all")
}

func server(ctx *exrouter.Context) {
	command := ctx.Args.Get(2)
	switch command {
	case "add":
		guild := ctx.Msg.GuildID
		host := ctx.Args.Get(3)
		pass := ctx.Args.Get(4)
		port, err := strconv.Atoi(ctx.Args.Get(5))
		if err != nil || port < 0 || port > 65535 {
			log.Fatalln(err)
			ctx.Reply("Port has to be a number between 0 and 65536.")
		}
		addServer(guild, host, pass, port)
		ctx.Reply(fmt.Sprintf("Added Minecraft server `%s`", host))
	case "list":
		guild := ctx.Msg.GuildID
		servers := getServers(guild)
		reply := "```\n"
		for id, value := range servers {
			reply += fmt.Sprintf("%3d %20s %6d %5t\n", id, value.Host, value.Port, value.Prim)
		}
		reply += "```"
		ctx.Reply(reply)
	}
}
