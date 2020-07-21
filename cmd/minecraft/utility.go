package minecraft

import (
	"strconv"

	"github.com/EmilyBjartskular/BlueBerryBunsBot/cmd"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bearbin/mcgorcon"
)

var _ = mcgorcon.Dial // For debugging; delete when done

// Add adds the minecraft command package to a CommandController
func Add(c *cmd.CommandController) {
	c.Router.On("mc", func(ctx *exrouter.Context) {
		arg1 := ctx.Args.Get(1)
		switch arg1 {
		case "add":
			guild := ctx.Msg.GuildID
			host := ctx.Args.Get(2)
			port := ctx.Args.Get(3)
			pass, err := strconv.Atoi(ctx.Args.Get(4))
			if err != nil || pass < 0 {
				ctx.Reply("Port has to be a number between 0 and 65536.")
			}
			addServer(guild, host, port, pass)
		default:
		}
		ctx.Reply("")
	}).Desc("Well, it's Minecraft y'all")
}
