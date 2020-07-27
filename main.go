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
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/pelletier/go-toml"
)

// Version is a constant that stores the Disgord version information.
const Version = "v0.0.0-alpha"

var (
	// Secrets contains all key/value pairs defined in the secrets toml file.
	Secrets, _ = toml.LoadFile("secrets.toml")

	// Config contains server specific config
	Config, _ = toml.LoadFile("config.toml")

	// Session is declared in the global space so it can be easily used
	// throughout this program.
	// In this use case, there is no error that would be returned.
	Session, _ = discordgo.New("Bot " + Secrets.Get("Discord.token").(string))

	// Router handles bot commands
	Router = exrouter.New()

	// Prefix is the command prefix for the bot.
	Prefix = Config.Get("Discord.prefix").(string)

	dbname = Secrets.Get("Database.name").(string)
	dbhost = Secrets.Get("Database.host").(string)
	dbuser = Secrets.Get("Database.user").(string)
	dbpass = Secrets.Get("Database.pass").(string)
	dbport = Secrets.Get("Database.port").(int64)
)

// Read in all configuration options from both environment variables and
// command line arguments.
func init() {
	if Session.Token == "" {
		flag.StringVar(&Session.Token, "t", "", "Discord Authentication Token")
	}
}

func main() {
	// Declare any variables needed later.
	var err error

	// Parse command line arguments
	flag.Parse()

	// Verify a Token was provided
	if Session.Token == "" {
		log.Println("You must provide a Discord authentication token.")
		return
	}

	// Create a command controller
	cc := cmd.New(Router, Prefix)
	cmd.Add(cc)

	db.Init(dbuser, dbpass, dbhost, dbport, dbname)

	// Get minecraft config and Add minecraft commands
	minecraft.Init()
	minecraft.Add(cc)

	// Handle Discord Events
	Session.AddHandler(guildCreate)
	Session.AddHandler(addRouter)

	// Open a websocket connection to Discord
	err = Session.Open()
	if err != nil {
		log.Printf("error opening connection to Discord, %s\n", err)
		os.Exit(1)
	}

	// Wait for a CTRL-C
	log.Printf(`Now running. Press CTRL-C to exit.`)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Clean up
	Session.Close()

	// Exit Normally.
}
