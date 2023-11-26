package handler

import (
	"fpl_discord_bot/message/cmd"
	"fpl_discord_bot/message/commands"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
	"time"
)

const prefix string = "!fpl"

var allowedCommands = map[string][]string{
	"1174296046155866112": {cmd.Hello},
}

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	args := strings.Split(m.Content, " ")
	if len(args) == 0 || args[0] != prefix || len(args) == 1 {
		return
	}

	command := args[1]
	channelID := m.ChannelID

	if isCommandAllowedInChannel(command, channelID) {
		switch command {
		case cmd.Hello:
			commands.HandleHello(s, m)
		default:
			InformAndDelete(s, m.Message, "Unknown command")
		}
	} else {
		InformAndDelete(s, m.Message, "This command is not allowed in this channel")
	}
}

func isCommandAllowedInChannel(command, channelID string) bool {
	allowedCommands, exists := allowedCommands[channelID]
	if !exists {
		return false
	}

	for _, allowedCommand := range allowedCommands {
		if allowedCommand == command {
			return true
		}
	}
	return false
}

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
