package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/EmilyBjartskular/BlueBerryBunsBot/cmd"
	"github.com/EmilyBjartskular/BlueBerryBunsBot/cmd/minecraft"
	"github.com/EmilyBjartskular/BlueBerryBunsBot/db"
	"github.com/EmilyBjartskular/BlueBerryBunsBot/util/config"
	"github.com/EmilyBjartskular/BlueBerryBunsBot/util/emoji"
	"github.com/EmilyBjartskular/BlueBerryBunsBot/util/global"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
)

// Version is a constant that stores the Disgord version information.
const Version = "v0.0.0-alpha"

var ()

func main() {
	// Declare any variables needed later.
	var err error

	// Initialize global discord Session
	global.Session, err = discordgo.New("Bot " + config.Discord.Token)
	if err != nil {
		log.Fatalln(err)
	}

	// Initialize global command router
	global.Router = exrouter.New()

	// Parse command line arguments
	flag.Parse()

	// Verify a Token was provided
	if global.Session.Token == "" {
		log.Println("You must provide a Discord authentication token.")
		return
	}

	// Create a command controller
	cmd.Add(global.Router)

	// Initialize database
	db.Init(config.Database.User,
		config.Database.Pass,
		config.Database.Host,
		config.Database.Port,
		config.Database.Name)

	// Handle Discord Events
	global.Session.AddHandler(guildCreate)
	global.Session.AddHandler(addRouter)

	// Open a websocket connection to Discord
	err = global.Session.Open()
	if err != nil {
		log.Printf("error opening connection to Discord, %s\n", err)
		os.Exit(1)
	}

	// Initialize global emoji state
	emoji.Init(global.Session)

	// Get minecraft config and Add minecraft commands
	minecraft.Init()
	minecraft.Add(global.Router)

	// Wait for a CTRL-C
	log.Printf(`Now running. Press CTRL-C to exit.`)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Clean up
	global.Session.Close()

	// Exit Normally.
}
