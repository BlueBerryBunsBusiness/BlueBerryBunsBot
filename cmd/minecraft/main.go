package minecraft

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/EmilyBjartskular/BlueBerryBunsBot/cmd/minecraft/db"
	mcsrv "github.com/EmilyBjartskular/BlueBerryBunsBot/cmd/minecraft/server"
	"github.com/EmilyBjartskular/BlueBerryBunsBot/util/emoji"
	"github.com/EmilyBjartskular/BlueBerryBunsBot/util/text"
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

// Init Initializes minecraft module
func Init() {
	db.Init()
}

// Add adds the minecraft command package to a CommandController
func Add(r *exrouter.Route) {
	mc := r.On(text.Minecraft.Mc.Command, mcFunc)
	mc.Desc(text.Minecraft.Mc.Description)
	mc.Cat(fmt.Sprintf("%s %s", emoji.Emojis["mc"].MessageFormat(), "Minecraft"))

	server := mc.On(text.Minecraft.Mc.Server.Command, serverFunc)
	server.Desc(text.Minecraft.Mc.Server.Description)

	serverAdd := server.On(text.Minecraft.Mc.Server.Add.Command, serverAddFunc)
	serverAdd.Desc(text.Minecraft.Mc.Server.Add.Description)

	serverList := server.On(text.Minecraft.Mc.Server.List.Command, serverListFunc)
	serverList.Desc(text.Minecraft.Mc.Server.List.Description)

	serverRemove := server.On(text.Minecraft.Mc.Server.Remove.Command, serverRemoveFunc)
	serverRemove.Desc(text.Minecraft.Mc.Server.Remove.Description)

	serverPrim := server.On(text.Minecraft.Mc.Server.Primary.Command, serverPrimFunc)
	serverPrim.Desc(text.Minecraft.Mc.Server.Primary.Description)

	mcsrv.Add(mc)

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
	port, err := strconv.Atoi(ctx.Args.Get(2))
	pass := ctx.Args.Get(3)
	prim := ctx.Args.Get(4)
	if host == "" || pass == "" {
		ctx.Reply(helpPrefix + cmdPrefix + cmdServerAdd + helpSuffix)
		return
	}
	if err != nil || port < 0 || port > 65535 {
		log.Println(err)
		ctx.Reply(text.Minecraft.Errors.InvalidPort)
		return
	}
	err = db.AddServer(guild, host, pass, port, prim)
	if err != nil {
		if strings.Contains(err.Error(), "Error 1062") {
			ctx.Reply(fmt.Sprintf(text.Minecraft.Errors.ServerExists, host))
			return
		}
		ctx.Reply(text.Minecraft.Errors.ServerSide)
		return
	}

	ctx.Reply(fmt.Sprintf(text.Minecraft.Mc.Server.Add.Reply, host))
}

func serverListFunc(ctx *exrouter.Context) {
	guild := ctx.Msg.GuildID
	servers, prim, err := db.GetServers(guild)
	if err != nil {
		ctx.Reply(text.Minecraft.Errors.ServerSide)
		return
	}
	if len(servers) == 0 {
		ctx.Reply(text.Minecraft.Errors.NoServers)
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
	reply += fmt.Sprintf(text.Minecraft.Mc.Server.List.TableFrame, lineID, lineHost)
	reply += fmt.Sprintf(text.Minecraft.Mc.Server.List.TableHeader, "P", lengthID, headerID, lengthHost, headerHost)
	reply += fmt.Sprintf(text.Minecraft.Mc.Server.List.TableFrame, lineID, lineHost)
	for _, value := range servers {
		p := " "
		if value.ID == prim {
			p = "*"
		}
		reply += fmt.Sprintf(text.Minecraft.Mc.Server.List.TableRow, p, lengthID, value.ID, lengthHost, value.Host)
	}
	reply += fmt.Sprintf(text.Minecraft.Mc.Server.List.TableFrame, lineID, lineHost)
	reply += "```"
	ctx.Reply(reply)
}

func serverRemoveFunc(ctx *exrouter.Context) {
	guild := ctx.Msg.GuildID
	id, err := strconv.Atoi(ctx.Args.Get(1))
	if err != nil {
		reply := text.Minecraft.Errors.InvalidID + "\n"
		reply += helpPrefix + cmdPrefix + cmdServerRemove + helpSuffix
		ctx.Reply(reply)
	}
	srv, err := db.DeleteServer(guild, id)
	if err != nil {
		ctx.Reply(text.Minecraft.Errors.ServerSide)
		return
	}
	ctx.Reply(fmt.Sprintf(text.Minecraft.Mc.Server.Remove.Reply, srv.Host))
}

func serverPrimFunc(ctx *exrouter.Context) {
	guild := ctx.Msg.GuildID
	id, err := strconv.Atoi(ctx.Args.Get(1))
	if err != nil {
		reply := text.Minecraft.Errors.InvalidID + "\n"
		reply += helpPrefix + cmdPrefix + cmdServerRemove + helpSuffix
		ctx.Reply(reply)
		return
	}
	srv, err := db.SetPrimary(guild, id)
	if err != nil {
		if err == db.ErrInvID {
			reply := text.Minecraft.Errors.InvalidID + "\n"
			reply += helpPrefix + cmdPrefix + cmdServerRemove + helpSuffix
			ctx.Reply(reply)
			return
		}
		ctx.Reply(text.Minecraft.Errors.ServerSide)
		return
	}

	reply := fmt.Sprintf(text.Minecraft.Mc.Server.Primary.Reply, srv.Host)
	ctx.Reply(reply)
}
