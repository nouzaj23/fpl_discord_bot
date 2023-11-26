package fetching

import (
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
	"time"
)

func HandleFetch(db *gorm.DB, s *discordgo.Session) {
	for {
		FetchAndUpdate(db, s)
		time.Sleep(time.Minute * 15)
	}
}
