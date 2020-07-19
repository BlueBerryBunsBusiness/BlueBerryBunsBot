package minecraft

import (
	"github.com/EmilyBjartskular/BlueBerryBunsBot/cmd"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bearbin/mcgorcon"
)

// Add adds the minecraft command package to a CommandController
func Add(c *cmd.CommandController, host string, port int, pass string) {
	c.Router.On("mc", func(ctx *exrouter.Context) {
		client, _ := mcgorcon.Dial(host, port, pass)
		reply, _ := client.SendCommand("list")
		ctx.Reply(reply)
	}).Desc("Well, it's Minecraft y'all")
}
