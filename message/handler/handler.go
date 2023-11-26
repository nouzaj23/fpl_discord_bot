package handler

import (
	"fpl_discord_bot/database"
	"fpl_discord_bot/message/cmd"
	"fpl_discord_bot/message/commands"
	"fpl_discord_bot/util"
	"github.com/bwmarrin/discordgo"
	"strings"
)

const prefix string = "!fpl"

var allowedCommands = map[string][]string{
	"1174296046155866112": {cmd.Hello},
	"1178335542962827275": {cmd.Player},
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
	db := database.GetDB()

	if isCommandAllowedInChannel(command, channelID) {
		switch command {
		case cmd.Hello:
			commands.HandleHello(s, m)
		case cmd.Player:
			commands.HandlePlayer(s, m, db, args[2:])
		default:
			util.InformAndDelete(s, m.Message, "Unknown command")
		}
	} else {
		util.InformAndDelete(s, m.Message, "This command is not allowed in this channel")
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
