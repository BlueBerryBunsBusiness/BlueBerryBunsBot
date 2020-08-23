package server

import (
	"fmt"

	"github.com/EmilyBjartskular/BlueBerryBunsBot/cmd/minecraft/db"
	"github.com/EmilyBjartskular/BlueBerryBunsBot/util/text"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bearbin/mcgorcon"
)

// Add server commands to minecraft commands
func Add(mc *exrouter.Route) {
	list := mc.On(text.Minecraft.Mc.List.Command, listFunc)
	list.Desc(text.Minecraft.Mc.List.Description)
}

func listFunc(ctx *exrouter.Context) {
	guild := ctx.Msg.GuildID
	server := ctx.Args.Get(1)

	id := 0
	if server == "[]" {
	}

	srv, err := db.GetServer(guild, id)
	if err != nil {
		ctx.Reply(text.Minecraft.Errors.NoServers)
		return
	}

	client, err := mcgorcon.Dial(srv.Host, int(srv.Port), srv.Pass)

	players, err := client.SendCommand("list")
	if err != nil {
		ctx.Reply(text.Minecraft.Errors.ServerSide)
		return
	}

	ctx.Reply(fmt.Sprintf(text.Minecraft.Mc.List.Reply, players))
}
