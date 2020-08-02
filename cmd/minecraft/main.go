package minecraft

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/EmilyBjartskular/BlueBerryBunsBot/util/emoji"
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

	cmdPrefix       string
	cmdServerAdd    string
	cmdServerList   string
	cmdServerRemove string
)

func init() {
	helpPrefix = "```asciidoc\n"
	helpSuffix = "```"
	cmdServerAdd = "mc server add HOST RCONPORT PASS [PRIMARY]\n"
	cmdServerList = "mc server list\n"
	cmdServerRemove = "mc server remove ID\n"
}

// Add adds the minecraft command package to a CommandController
func Add(r *exrouter.Route) {
	mc := r.On("mc", mcFunc)
	mc.Desc("Minecraft commands")
	mc.Cat(fmt.Sprintf("%s %s", emoji.Emojis["mc"].MessageFormat(), "Minecraft"))

	server := mc.On("server", serverFunc)
	server.Desc("Minecraft server commands")

	serverAdd := server.On("add", serverAddFunc)
	serverAdd.Desc("Add a minecraft server to the guild")

	serverList := server.On("list", serverListFunc)
	serverList.Desc("List minecraft servers in the guild")

	serverRemove := server.On("remove", serverRemoveFunc)
	serverRemove.Desc("Remove minecraft server from the guild")

	serverPrim := server.On("primary", serverPrimFunc)
	serverPrim.Desc("Set primary server in the guild")

	for _, v := range mc.Routes {
		helpMc += fmt.Sprintf("%s %-15s:: %s\n", v.Category, v.Name, v.Description)
	}

	for _, v := range server.Routes {
		helpServer += fmt.Sprintf("%s %-15s:: %s\n", v.Category, v.Name, v.Description)
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
	port, err := strconv.Atoi(ctx.Args.Get(3))
	pass := ctx.Args.Get(2)
	prim := ctx.Args.Get(4)
	if host == "" || pass == "" {
		ctx.Reply(helpPrefix + cmdPrefix + cmdServerAdd + helpSuffix)
		return
	}
	if err != nil || port < 0 || port > 65535 {
		log.Println(err)
		ctx.Reply("Port has to be a number between 0 and 65536.")
		return
	}
	err = addServer(guild, host, pass, port, prim)
	if err != nil {
		if strings.Contains(err.Error(), "Error 1062") {
			ctx.Reply(fmt.Sprintf("Minecraft server `%s` already exists", host))
			return
		}
		ctx.Reply("Something went wrong on the server side, please contact the bot developer with your Guild ID")
		return
	}

	ctx.Reply(fmt.Sprintf("Added Minecraft server `%s`", host))
}

func serverListFunc(ctx *exrouter.Context) {
	guild := ctx.Msg.GuildID
	servers, prim, err := getServers(guild)
	if err != nil {
		ctx.Reply("Something went wrong on the server side, please contact the bot developer with your Guild ID")
		return
	}
	if len(servers) == 0 {
		ctx.Reply("Guild has no servers, maybe add one? :sweat_smile:")
		return
	}

	headerID := "ID"
	headerHost := "Host"

	lengthHost := len(headerHost)
	for _, value := range servers {
		if l := len(value.Host); l > lengthHost {
			lengthHost = l
		}
	}

	lengthID := len(headerID)
	if l := len(strconv.Itoa(servers[len(servers)-1].ID)); l > lengthID {
		lengthID = l
	}

	lineID := strings.Repeat("-", lengthID+2)
	lineHost := strings.Repeat("-", lengthHost+2)

	reply := "```\n"
	reply += "+---+" + lineID + "+" + lineHost + "+\n"
	reply += fmt.Sprintf("| %s | %*s | %-*s |\n", "P", lengthID, headerID, lengthHost, headerHost)
	reply += "+---+" + lineID + "+" + lineHost + "+\n"
	for _, value := range servers {
		p := " "
		if value.ID == prim {
			p = "*"
		}
		reply += fmt.Sprintf("| %s | %*d | %-*s |\n", p, lengthID, value.ID, lengthHost, value.Host)
	}
	reply += "+---+" + lineID + "+" + lineHost + "+"
	reply += "```"
	ctx.Reply(reply)
}

func serverRemoveFunc(ctx *exrouter.Context) {
	guild := ctx.Msg.GuildID
	id, err := strconv.Atoi(ctx.Args.Get(1))
	if err != nil {
		reply := "Invalid server id.\n"
		reply += helpPrefix + cmdPrefix + cmdServerRemove + helpSuffix
		ctx.Reply(reply)
	}
	srv, err := deleteServer(guild, id)
	if err != nil {
		ctx.Reply("Something went wrong on the server side, please contact the bot developer with your Guild ID")
		return
	}
	ctx.Reply("Server `" + srv.Host + "` deleted")
}

func serverPrimFunc(ctx *exrouter.Context) {
	guild := ctx.Msg.GuildID
	id, err := strconv.Atoi(ctx.Args.Get(1))
	if err != nil {
		reply := "Invalid server id.\n"
		reply += helpPrefix + cmdPrefix + cmdServerRemove + helpSuffix
		ctx.Reply(reply)
		return
	}
	srv, err := setPrimary(guild, id)
	if err != nil {
		if err == ErrInvID {
			reply := "Invalid server id.\n"
			reply += helpPrefix + cmdPrefix + cmdServerRemove + helpSuffix
			ctx.Reply(reply)
			return
		}
		ctx.Reply("Something went wrong on the server side, please contact the bot developer with your Guild ID")
		return
	}

	reply := "Set `" + srv.Host + "` as primary server."
	ctx.Reply(reply)
}
