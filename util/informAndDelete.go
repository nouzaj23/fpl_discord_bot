package util

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"time"
)

func InformAndDelete(s *discordgo.Session, m *discordgo.Message, content string) {
	res, err := s.ChannelMessageSend(m.ChannelID, content)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}
	time.Sleep(5 * time.Second)
	err = s.ChannelMessageDelete(m.ChannelID, m.ID)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}
	err = s.ChannelMessageDelete(res.ChannelID, res.ID)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}
}
