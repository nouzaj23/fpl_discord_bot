package main

import (
	"fmt"
	"fpl_discord_bot/database"
	"fpl_discord_bot/fetching"
	"fpl_discord_bot/message/handler"
	"fpl_discord_bot/models"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	godotenv.Load()
	log.SetFlags(log.Ldate | log.Ltime)
	sess, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatalf("Failed to create new session: %v", err)
	}

	db := database.InitDB()
	err = db.AutoMigrate(&models.Player{}, &models.Team{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	sess.AddHandler(handler.MessageCreate)

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged
	err = sess.Open()
	if err != nil {
		log.Fatalf("Failed to open session: %v", err)
	}
	defer sess.Close()
	fmt.Println("The bot is running!")
	go fetching.HandleFetch(db, sess)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
