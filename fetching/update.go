package fetching

import (
	"encoding/json"
	"errors"
	"fmt"
	"fpl_discord_bot/models"
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const fetchURL = "https://fantasy.premierleague.com/api/bootstrap-static/"
const injuryNewsChannel = "1177893613649268776"
const priceChangesChannel = "1177893636252381214"
const newPlayersChannel = "1177919636398948444"

func FetchAndUpdate(db *gorm.DB, s *discordgo.Session) {
	resp, err := http.Get(fetchURL)
	if err != nil {
		log.Fatalf("Failed to fetch data: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to fetch data: %v", err)
	}

	var data TeamsAndPlayersData
	if err := json.Unmarshal(body, &data); err != nil {
		log.Printf("Failed to fetch data: %v", err)
	}

	updateTeams(db, data.Teams)
	updatePlayers(db, data.Players, s)
}

func updateTeams(db *gorm.DB, teams []FetchedTeam) {
	for _, fetchedTeam := range teams {
		var team models.Team
		newTeam := models.Team{
			ID:                  uint(fetchedTeam.ID),
			Name:                fetchedTeam.Name,
			ShortName:           fetchedTeam.ShortName,
			StrengthOverallHome: uint(fetchedTeam.StrengthOverallHome),
			StrengthOverallAway: uint(fetchedTeam.StrengthOverallAway),
			StrengthAttackHome:  uint(fetchedTeam.StrengthAttackHome),
			StrengthAttackAway:  uint(fetchedTeam.StrengthAttackAway),
			StrengthDefenceHome: uint(fetchedTeam.StrengthDefenceHome),
			StrengthDefenceAway: uint(fetchedTeam.StrengthDefenceAway),
		}

		dbSearch := db.First(&team, newTeam.ID)
		if errors.Is(dbSearch.Error, gorm.ErrRecordNotFound) {
			db.Create(&newTeam)
			log.Printf("New Team created: (%v) %v", newTeam.ID, newTeam.Name)
		} else {
			newTeam.Model = team.Model
			if newTeam != team {
				db.Save(&newTeam)
				log.Printf("Team updated: (%v) %v", newTeam.ID, newTeam.Name)
			}
		}
	}
}

func updatePlayers(db *gorm.DB, players []FetchedPlayer, s *discordgo.Session) {
	var injuryNewsBatch = map[uint][]string{}
	var newPlayersBatch = map[uint][]string{}
	var priceRisersBatch = map[uint][]string{}
	var priceFallersBatch = map[uint][]string{}
	for _, fetchedPlayer := range players {
		var player models.Player
		newPlayer := models.Player{
			ID:                 uint(fetchedPlayer.ID),
			Name:               fmt.Sprintf("%s %s", fetchedPlayer.FirstName, fetchedPlayer.SecondName),
			WebName:            fetchedPlayer.WebName,
			Position:           getPosition(fetchedPlayer.ElementType),
			Nationality:        "",
			TotalPoints:        fetchedPlayer.TotalPoints,
			TeamID:             uint(fetchedPlayer.Team),
			Price:              fetchedPlayer.NowCost,
			ChanceOfNextRound:  fetchedPlayer.ChanceOfNextRound,
			CostChange:         fetchedPlayer.CostChange,
			News:               fetchedPlayer.News,
			Minutes:            uint(fetchedPlayer.Minutes),
			Goals:              uint(fetchedPlayer.GoalsScored),
			Assists:            uint(fetchedPlayer.Assists),
			CleanSheets:        uint(fetchedPlayer.CleanSheets),
			GoalsConceded:      uint(fetchedPlayer.GoalsConceded),
			Saves:              uint(fetchedPlayer.Saves),
			PenaltiesSaved:     uint(fetchedPlayer.PenaltiesSaved),
			PenaltiesMissed:    uint(fetchedPlayer.PenaltiesMissed),
			YellowCards:        uint(fetchedPlayer.YellowCards),
			RedCards:           uint(fetchedPlayer.RedCards),
			PointsPerGame:      fetchedPlayer.PointsPerGame,
			SelectedByPercent:  fetchedPlayer.SelectedByPercent,
			Bonus:              uint(fetchedPlayer.Bonus),
			Bps:                fetchedPlayer.Bps,
			Starts:             uint(fetchedPlayer.Starts),
			XG:                 stringToFloat(fetchedPlayer.XG),
			XA:                 stringToFloat(fetchedPlayer.XA),
			XGI:                stringToFloat(fetchedPlayer.XGI),
			XGC:                stringToFloat(fetchedPlayer.XGC),
			XGper90:            fetchedPlayer.XGper90,
			SavesPer90:         fetchedPlayer.SavesPer90,
			XAper90:            fetchedPlayer.XAper90,
			XGIper90:           fetchedPlayer.XGIper90,
			XGCper90:           fetchedPlayer.XGCper90,
			GoalsConcededPer90: fetchedPlayer.GoalsConcededPer90,
			CleanSheetsPer90:   fetchedPlayer.CleanSheetsPer90,
		}
		var team models.Team
		db.First(&team, newPlayer.TeamID)
		dbSearch := db.First(&player, newPlayer.ID)
		if errors.Is(dbSearch.Error, gorm.ErrRecordNotFound) {
			db.Create(&newPlayer)
			newPlayersBatch[newPlayer.TeamID] = append(newPlayersBatch[newPlayer.TeamID],
				fmt.Sprintf("- %s (%s) - %s $%.1f",
					newPlayer.WebName,
					team.ShortName,
					newPlayer.Position,
					float32(newPlayer.Price)/10))
			log.Printf("New player created: (%v) %v", newPlayer.ID, newPlayer.Name)
		} else {
			newPlayer.Model = player.Model
			if newPlayer != player {
				if newPlayer.News != player.News {
					var news string
					if newPlayer.News == "" {
						news = "Available"
					} else {
						news = newPlayer.News
					}
					injuryNewsBatch[newPlayer.TeamID] = append(injuryNewsBatch[newPlayer.TeamID],
						fmt.Sprintf("- %s (%s) - %s", newPlayer.WebName, team.ShortName, news))
				}

				if newPlayer.Price < player.Price {
					priceFallersBatch[newPlayer.TeamID] = append(priceFallersBatch[newPlayer.TeamID],
						fmt.Sprintf("- %s (%s): £%.1f -> £%.1f", newPlayer.WebName, team.ShortName, float32(player.Price)/10, float32(newPlayer.Price)/10))
				}

				if newPlayer.Price > player.Price {
					priceRisersBatch[newPlayer.TeamID] = append(priceRisersBatch[newPlayer.TeamID],
						fmt.Sprintf("- %s (%s): £%.1f -> £%.1f", newPlayer.WebName, team.ShortName, float32(player.Price)/10, float32(newPlayer.Price)/10))
				}
				db.Save(&newPlayer)
				log.Printf("Player updated: (%v) %v", newPlayer.ID, newPlayer.Name)
			}
		}
	}
	go exportInjuryNews(injuryNewsBatch, s)
	go exportPriceChanges(priceRisersBatch, priceFallersBatch, s)
	go exportNewPlayers(newPlayersBatch, s)
}

func getPosition(elementType int) string {
	switch elementType {
	case 1:
		return "Goalkeeper"
	case 2:
		return "Defender"
	case 3:
		return "Midfielder"
	case 4:
		return "Forward"
	default:
		return "Unknown"
	}
}

func stringToFloat(s string) float32 {
	result, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0.0
	}
	return float32(result)
}

func exportInjuryNews(injuryNewsBatch map[uint][]string, s *discordgo.Session) {
	if len(injuryNewsBatch) == 0 {
		return
	}
	var result string
	result += "# New injury updates ⚠️ \n"
	for _, teamNews := range injuryNewsBatch {
		result += strings.Join(teamNews, "\n")
		result += "\n"
	}
	result = strings.TrimSuffix(result, "\n")
	_, err := s.ChannelMessageSend(injuryNewsChannel, result)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}
}

func exportPriceChanges(priceRisersBatch map[uint][]string, priceFallersBatch map[uint][]string, s *discordgo.Session) {
	if len(priceRisersBatch) == 0 && len(priceFallersBatch) == 0 {
		return
	}
	var result string
	if len(priceRisersBatch) > 0 {
		result += "# Price risers 📈 \n"
		for _, priceRisers := range priceRisersBatch {
			result += strings.Join(priceRisers, "\n")
			result += "\n"
		}
	}
	if len(priceFallersBatch) > 0 {
		result += "# Price fallers 📉 \n"
		for _, priceFallers := range priceFallersBatch {
			result += strings.Join(priceFallers, "\n")
			result += "\n"
		}
	}
	result = strings.TrimSuffix(result, "\n")
	_, err := s.ChannelMessageSend(priceChangesChannel, result)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}
}

func exportNewPlayers(newPlayersBatch map[uint][]string, s *discordgo.Session) {
	if len(newPlayersBatch) == 0 {
		return
	}
	var result string
	result += "# New players 🆕 \n"
	for _, newPlayers := range newPlayersBatch {
		result += strings.Join(newPlayers, "\n")
		result += "\n"
	}
	result = strings.TrimSuffix(result, "\n")
	_, err := s.ChannelMessageSend(newPlayersChannel, result)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}
}
