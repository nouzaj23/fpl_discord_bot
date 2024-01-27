package commands

import (
	"fmt"
	"fpl_discord_bot/message"
	"fpl_discord_bot/models"
	"fpl_discord_bot/repository"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

func HandlePlayer(s *discordgo.Session, m *discordgo.MessageCreate, pr repository.PlayerRepository, tr repository.TeamRepository, args []string) {
	if len(args) == 0 {
		message.InformAndDelete(s, m.Message, "Missing player name. Usage:\n!fpl player <player_name>")
		return
	}

	var name string
	if len(args) > 1 {
		name = strings.Join(args, " ")
	} else {
		name = args[0]
	}
	players, _ := pr.FindByName(name)

	if len(players) == 0 {
		message.InformAndDelete(s, m.Message, fmt.Sprintf("No player with name '%s' found, please try again", name))
		return
	}

	if len(players) > 1 {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("**Multiple players found matching the name %s**.\n", name))
		for _, player := range players {
			team, _ := tr.Find(player.TeamID)
			sb.WriteString(fmt.Sprintf("- %s (%s) - %s\n", player.Name, team.ShortName, player.Position))
		}
		sb.WriteString("Please try again with a more precise name")
		message.InformAndDelete(s, m.Message, sb.String())
		return
	}

	player := players[0]
	team, _ := tr.Find(player.TeamID)
	response := buildResponse(player, *team)

	_, err := s.ChannelMessageSendReply(m.ChannelID, response, m.Reference())
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}
}

func buildResponse(player models.Player, team models.Team) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Player Info for %s\n", player.Name))
	sb.WriteString(fmt.Sprintf("**Position:** %s\n", player.Position))
	sb.WriteString(fmt.Sprintf("**Team:** %s\n", team.Name))
	if player.Nationality != "" {
		sb.WriteString(fmt.Sprintf("**Nationality:** %s\n", player.Nationality))
	}
	sb.WriteString("**Selected By:** " + player.SelectedByPercent + "%\n")
	sb.WriteString("\n**Main FPL Stats:**\n")
	sb.WriteString(fmt.Sprintf("- **Price:** £%.1f\n", float32(player.Price)/10))
	sb.WriteString(fmt.Sprintf("- **Total Points:** %d\n", player.TotalPoints))
	sb.WriteString(fmt.Sprintf("- **Points Per Game:** %s\n", player.PointsPerGame))
	sb.WriteString(fmt.Sprintf("- **Goals:** %d\n", player.Goals))
	sb.WriteString(fmt.Sprintf("- **Assists:** %d\n", player.Assists))
	sb.WriteString(fmt.Sprintf("- **Games in Starting 11:** %d\n", player.Starts))
	sb.WriteString(fmt.Sprintf("- **Minutes Played:** %d\n", player.Minutes))

	if player.Position != "Forward" && player.Position != "Midfielder" {
		sb.WriteString(fmt.Sprintf("- **Clean Sheets:** %d\n", player.CleanSheets))
		sb.WriteString(fmt.Sprintf("- **Clean Sheets Per 90:** %.2f\n", float64(player.CleanSheetsPer90)))
		sb.WriteString(fmt.Sprintf("- **Goals Conceded:** %d\n", player.GoalsConceded))
		sb.WriteString(fmt.Sprintf("- **Goals Conceded Per 90:** %.2f\n", float64(player.GoalsConcededPer90)))
	}

	if player.Position == "Goalkeeper" {
		sb.WriteString(fmt.Sprintf("- **Penalties Saved:** %d\n", player.PenaltiesSaved))
	}
	sb.WriteString("\n**Disciplinary Records:**\n")
	sb.WriteString(fmt.Sprintf("- **Yellow Cards:** %d\n", player.YellowCards))
	sb.WriteString(fmt.Sprintf("- **Red Cards:** %d\n", player.RedCards))
	sb.WriteString("\n**Advanced Stats:**\n")
	sb.WriteString(fmt.Sprintf("- **xG (Expected Goals):** %.2f\n", float64(player.XG)))
	sb.WriteString(fmt.Sprintf("- **xG (Expected Goals) Per 90:** %.2f\n", float64(player.XGper90)))
	sb.WriteString(fmt.Sprintf("- **xA (Expected Assists):** %.2f\n", float64(player.XA)))
	sb.WriteString(fmt.Sprintf("- **xA (Expected Assists) Per 90:** %.2f\n", float64(player.XAper90)))
	sb.WriteString(fmt.Sprintf("- **xGI (Expected Goal Involvement):** %.2f\n", float64(player.XGI)))
	sb.WriteString(fmt.Sprintf("- **xGI (Expected Goal Involvement) Per 90:** %.2f\n", float64(player.XGIper90)))

	if player.Position == "Goalkeeper" || player.Position == "Defender" {
		sb.WriteString(fmt.Sprintf("- **xGC (Expected Goals Conceded):** %.2f\n", float64(player.XGC)))
		sb.WriteString(fmt.Sprintf("- **XGC (Expected Goals Conceded) Per 90:** %.2f\n", float64(player.XGCper90)))
	}

	sb.WriteString("\n**Miscellaneous Information:**\n")
	sb.WriteString(fmt.Sprintf("- **News:** %s\n", player.News))
	sb.WriteString(fmt.Sprintf("- **Cost Change:** %.1f\n", float32(player.CostChange)/10))
	sb.WriteString(fmt.Sprintf("- **Bonus:** %d\n", player.Bonus))
	sb.WriteString(fmt.Sprintf("- **BPS (Bonus Points System):** %d\n", player.Bps))
	return sb.String()
}
