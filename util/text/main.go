package text

import (
	"fmt"
	"log"

	"github.com/EmilyBjartskular/BlueBerryBunsBot/util/global"
	"github.com/pelletier/go-toml"
)

var (
	// Generic stores generic text
	Generic generic
	// Minecraft stores minecraft text
	Minecraft minecraft
)

type command struct {
	Command     string
	Description string
	Reply       string
}

type generic struct {
	EmptyHelp        string
	EmbedTitle       string
	EmbedDescription string
}

type minecraft struct {
	Mc     minecraftMc
	Errors minecraftErrors
}

type minecraftErrors struct {
	InvalidID    string
	InvalidPort  string
	ServerExists string
	ServerSide   string
	NoServers    string
}

type minecraftMc struct {
	command
	Server minecraftMcServer
}

type minecraftMcServer struct {
	command
	Add     minecraftMcServerAdd
	List    minecraftMcServerList
	Remove  minecraftMcServerRemove
	Primary minecraftMcServerPrimary
}

type minecraftMcServerAdd struct {
	command
}

type minecraftMcServerList struct {
	command
	TableFrame  string
	TableHeader string
	TableRow    string
}

type minecraftMcServerRemove struct {
	command
}

type minecraftMcServerPrimary struct {
	command
}

func init() {
	tomlFile, err := toml.LoadFile(global.TextConfigPath)
	if err != nil {
		log.Fatalln(err)
	}

	// GENERIC
	Generic.EmptyHelp = getProperty("Generic.emptyHelp", tomlFile).(string)
	Generic.EmbedTitle = getProperty("Generic.embedTitle", tomlFile).(string)
	Generic.EmbedDescription = getProperty("Generic.embedDescription", tomlFile).(string)

	// MINECRAFT
	Minecraft.Errors.InvalidID = getProperty("Minecraft.Errors.invalidId", tomlFile).(string)
	Minecraft.Errors.InvalidPort = getProperty("Minecraft.Errors.invalidPort", tomlFile).(string)
	Minecraft.Errors.ServerExists = getProperty("Minecraft.Errors.serverExists", tomlFile).(string)
	Minecraft.Errors.ServerSide = getProperty("Minecraft.Errors.serverSide", tomlFile).(string)
	Minecraft.Errors.NoServers = getProperty("Minecraft.Errors.noServers", tomlFile).(string)
	Minecraft.Mc.Command = getProperty("Minecraft.Mc.command", tomlFile).(string)
	Minecraft.Mc.Description = getProperty("Minecraft.Mc.description", tomlFile).(string)
	Minecraft.Mc.Server.Command = getProperty("Minecraft.Mc.Server.command", tomlFile).(string)
	Minecraft.Mc.Server.Description = getProperty("Minecraft.Mc.Server.description", tomlFile).(string)
	Minecraft.Mc.Server.Add.Command = getProperty("Minecraft.Mc.Server.Add.command", tomlFile).(string)
	Minecraft.Mc.Server.Add.Description = getProperty("Minecraft.Mc.Server.Add.description", tomlFile).(string)
	Minecraft.Mc.Server.List.Command = getProperty("Minecraft.Mc.Server.List.command", tomlFile).(string)
	Minecraft.Mc.Server.List.Description = getProperty("Minecraft.Mc.Server.List.description", tomlFile).(string)
	Minecraft.Mc.Server.Remove.Command = getProperty("Minecraft.Mc.Server.Remove.command", tomlFile).(string)
	Minecraft.Mc.Server.Remove.Description = getProperty("Minecraft.Mc.Server.Remove.description", tomlFile).(string)
	Minecraft.Mc.Server.Primary.Command = getProperty("Minecraft.Mc.Server.Primary.command", tomlFile).(string)
	Minecraft.Mc.Server.Primary.Description = getProperty("Minecraft.Mc.Server.Primary.description", tomlFile).(string)
}

func getProperty(prop string, t *toml.Tree) interface{} {
	if !t.Has(prop) {
		log.Fatalln(fmt.Sprintf("Missing `%s` in `%s`", prop, global.TextConfigPath))
	}
	return t.Get(prop)
}
