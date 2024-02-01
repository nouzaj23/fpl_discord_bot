package fetching

import (
	"fpl_discord_bot/repository"
	"github.com/bwmarrin/discordgo"
	"time"
)

func HandleFetch(pr repository.PlayerRepository, tr repository.TeamRepository, s *discordgo.Session) {
	FetchAndUpdate(pr, tr, s)
	ticker := time.NewTicker(time.Minute * 15)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			FetchAndUpdate(pr, tr, s)
		}
	}
}
