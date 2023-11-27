package commands

import (
	"fmt"
	"fpl_discord_bot/models"
	"fpl_discord_bot/util"
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
	"log"
	"strings"
)

func HandlePlayer(s *discordgo.Session, m *discordgo.MessageCreate, db *gorm.DB, args []string) {
	if len(args) == 0 {
		util.InformAndDelete(s, m.Message, "Missing player name. Usage:\n!fpl player <player_name>")
		return
	}

	var name string
	if len(args) > 1 {
		name = strings.Join(args, " ")
	} else {
		name = args[0]
	}
	var players []models.Player
	db.Where("Name LIKE ?", "%"+name+"%").Find(&players)

	if len(players) == 0 {
		util.InformAndDelete(s, m.Message, fmt.Sprintf("No player with name '%s' found, please try again", name))
		return
	}

	if len(players) > 1 {
		response := fmt.Sprintf("**Multiple players found matching the name %s**.\n", name)
		for _, player := range players {
			var team models.Team
			db.First(&team, player.TeamID)
			response += fmt.Sprintf("- %s (%s) - %s\n", player.Name, team.ShortName, player.Position)
		}
		response += "Please try again with a more precise name"
		util.InformAndDelete(s, m.Message, response)
		return
	}

	player := players[0]
	var team models.Team
	db.First(&team, player.TeamID)
	response := buildResponse(player, team)

	_, err := s.ChannelMessageSendReply(m.ChannelID, response, m.Reference())
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}
}

func buildResponse(player models.Player, team models.Team) string {
	response := fmt.Sprintf("# Player Info for %s\n", player.Name)
	response += fmt.Sprintf("**Position:** %s\n", player.Position)
	response += fmt.Sprintf("**Team:** %s\n", team.Name)
	if player.Nationality != "" {
		response += fmt.Sprintf("**Nationality:** %s\n", player.Nationality)
	}
	response += "**Selected By:** " + player.SelectedByPercent + "%\n"
	response += "\n**Main FPL Stats:**\n"
	response += fmt.Sprintf("- **Price:** £%.1f\n", float32(player.Price)/10)
	response += fmt.Sprintf("- **Total Points:** %d\n", player.TotalPoints)
	response += fmt.Sprintf("- **Points Per Game:** %s\n", player.PointsPerGame)
	response += fmt.Sprintf("- **Goals:** %d\n", player.Goals)
	response += fmt.Sprintf("- **Assists:** %d\n", player.Assists)
	response += fmt.Sprintf("- **Games in Starting 11:** %d\n", player.Starts)
	response += fmt.Sprintf("- **Minutes Played:** %d\n", player.Minutes)

	if player.Position != "Forward" && player.Position != "Midfielder" {
		response += fmt.Sprintf("- **Clean Sheets:** %d\n", player.CleanSheets)
		response += fmt.Sprintf("- **Clean Sheets Per 90:** %.2f\n", util.RoundFloat(float64(player.CleanSheetsPer90), 2))
		response += fmt.Sprintf("- **Goals Conceded:** %d\n", player.GoalsConceded)
		response += fmt.Sprintf("- **Goals Conceded Per 90:** %.2f\n", util.RoundFloat(float64(player.GoalsConcededPer90), 2))
	}

	if player.Position == "Goalkeeper" {
		response += fmt.Sprintf("- **Penalties Saved:** %d\n", player.PenaltiesSaved)
	}
	response += "\n**Disciplinary Records:**\n"
	response += fmt.Sprintf("- **Yellow Cards:** %d\n", player.YellowCards)
	response += fmt.Sprintf("- **Red Cards:** %d\n", player.RedCards)
	response += "\n**Advanced Stats:**\n"
	response += fmt.Sprintf("- **xG (Expected Goals):** %.2f\n", util.RoundFloat(float64(player.XG), 2))
	response += fmt.Sprintf("- **xG (Expected Goals) Per 90:** %.2f\n", util.RoundFloat(float64(player.XGper90), 2))
	response += fmt.Sprintf("- **xA (Expected Assists):** %.2f\n", util.RoundFloat(float64(player.XA), 2))
	response += fmt.Sprintf("- **xA (Expected Assists) Per 90:** %.2f\n", util.RoundFloat(float64(player.XAper90), 2))
	response += fmt.Sprintf("- **xGI (Expected Goal Involvement):** %.2f\n", util.RoundFloat(float64(player.XGI), 2))
	response += fmt.Sprintf("- **xGI (Expected Goal Involvement) Per 90:** %.2f\n", util.RoundFloat(float64(player.XGIper90), 2))

	if player.Position == "Goalkeeper" || player.Position == "Defender" {
		response += fmt.Sprintf("- **xGC (Expected Goals Conceded):** %.2f\n", util.RoundFloat(float64(player.XGC), 2))
		response += fmt.Sprintf("- **XGC (Expected Goals Conceded) Per 90:** %.2f\n", util.RoundFloat(float64(player.XGCper90), 2))
	}

	response += "\n**Miscellaneous Information:**\n"
	response += fmt.Sprintf("- **News:** %s\n", player.News)
	response += fmt.Sprintf("- **Cost Change:** %.1f\n", float32(player.CostChange)/10)
	response += fmt.Sprintf("- **Bonus:** %d\n", player.Bonus)
	response += fmt.Sprintf("- **BPS (Bonus Points System):** %d\n", player.Bps)
	return response
}
