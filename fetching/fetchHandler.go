package fetching

import (
	"fpl_discord_bot/repository"
	"github.com/bwmarrin/discordgo"
	"time"
)

func HandleFetch(pr repository.PlayerRepository, tr repository.TeamRepository, s *discordgo.Session) {
	for {
		FetchAndUpdate(pr, tr, s)
		time.Sleep(time.Minute * 15)
	}
}
