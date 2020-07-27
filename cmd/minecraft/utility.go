package minecraft

import (
	"fmt"
	"log"
	"strconv"

	"github.com/EmilyBjartskular/BlueBerryBunsBot/cmd"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bearbin/mcgorcon"
)

var _ = mcgorcon.Dial // For debugging; delete when done

var (
	// Help is the minecraft command help message
	helpPrefix string
	helpSuffix string
	helpMc     string
	helpServer string
	helpAdd    string
	helpList   string
)

func init() {
	helpPrefix = "```asciidoc\n"
	helpSuffix = "```"
	helpServer = ""
	helpAdd = ""
	helpList = ""
	helpMc = ""
}

// Add adds the minecraft command package to a CommandController
func Add(c *cmd.CommandController) {
	mc := c.Router.On("mc", mcFunc)
	mc.Desc("Minecraft commands")

	server := mc.On("server", serverFunc)
	server.Desc("Minecraft server commands")

	serverAdd := server.On("add", serverAddFunc)
	serverAdd.Desc("Add a minecraft server to the guild")

	serverList := server.On("list", serverListFunc)
	serverList.Desc("List minecraft servers in the guild")

	serverRemove := server.On("remove", serverRemoveFunc)
	serverRemove.Desc("Remove minecraft server from the guild")

	for _, v := range c.Router.Routes {
		helpMc += fmt.Sprintf("* %-15s:: %s\n", v.Name, v.Description)
	}
}

func mcFunc(ctx *exrouter.Context) {
	ctx.Reply(helpPrefix + helpMc + helpSuffix)
}

func serverFunc(ctx *exrouter.Context) {
	ctx.Reply(helpPrefix + helpServer + helpSuffix)
}

func serverAddFunc(ctx *exrouter.Context) {
	guild := ctx.Msg.GuildID
	host := ctx.Args.Get(1)
	pass := ctx.Args.Get(2)
	port, err := strconv.Atoi(ctx.Args.Get(3))
	prim := ctx.Args.Get(4)
	if host == "" || pass == "" {
		ctx.Reply(helpPrefix + helpAdd + helpSuffix)
		return
	}
	if err != nil || port < 0 || port > 65535 {
		log.Println(err)
		ctx.Reply("Port has to be a number between 0 and 65536.")
		return
	}
	addServer(guild, host, pass, port, prim)
	ctx.Reply(fmt.Sprintf("Added Minecraft server `%s`", host))
}

func serverListFunc(ctx *exrouter.Context) {
	guild := ctx.Msg.GuildID
	servers := getServers(guild)
	if len(servers) == 0 {
		ctx.Reply("Guild has no servers, maybe add one? :sweat_smile:")
		return
	}
	reply := "```\n"
	for _, value := range servers {
		reply += fmt.Sprintf("%3d %20s %6d\n", value.id, value.Host, value.Port)
	}
	reply += "```"
	ctx.Reply(reply)
}

func serverRemoveFunc(ctx *exrouter.Context) {
}
