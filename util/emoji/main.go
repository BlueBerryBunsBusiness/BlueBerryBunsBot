package emoji

import (
	"log"

	"github.com/EmilyBjartskular/BlueBerryBunsBot/util/config"
	"github.com/bwmarrin/discordgo"
)

var (
	// Emojis stores custom emojis from a source guild
	Emojis map[string]*discordgo.Emoji
)

// Init initializes the global emoji list.
// Init can also be used to update the list.
func Init(session *discordgo.Session) {
	emojis, err := session.GuildEmojis(config.Emoji.Guild)
	if err != nil {
		log.Fatalln(config.Emoji.Guild, err)
	}

	Emojis = map[string]*discordgo.Emoji{}
	for _, v := range emojis {
		Emojis[v.Name] = v
		log.Println(config.Emoji.Guild, v.Name)
	}
}
