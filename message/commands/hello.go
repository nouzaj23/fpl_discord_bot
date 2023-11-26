package commands

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

func HelloCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	_, err := s.ChannelMessageSend(m.ChannelID, "world!")
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}
}
