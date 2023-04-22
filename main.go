package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

func main() {

	// Get ENV from file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("ERROR loading .env file")
	}

	// Set apiToken from ENV
	apiToken := os.Getenv("TELEGRAM_BOT_TOKEN")

	// Load channels IDs to resend
	channelIDs, err := loadChannelIDs("channels")
	if err != nil {
		log.Fatalf("Couldn`t load channels` IDs: %v", err)
	}

	// Set bot
	bot, err := tgbotapi.NewBotAPI(apiToken)
	if err != nil {
		log.Fatalf("Couldn`t create bot: %v", err)
	}

	log.Printf("Autorized on %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Get updates
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatalf("Couldn`t get update: %v", err)
	}

	for update := range updates {
		mainTelegramChannel, _ := strconv.Atoi(os.Getenv("MAIN_TELEGRAM_CHANNEL"))
		if update.ChannelPost != nil && update.ChannelPost.Chat.ID == int64(mainTelegramChannel) {
			forwardToChannels(bot, update.ChannelPost, channelIDs)
		}
	}
}

func loadChannelIDs(filepath string) ([]int64, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	channelIDs := make([]int64, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var id int64
		_, err := fmt.Sscanf(scanner.Text(), "%d", &id)
		if err == nil {
			channelIDs = append(channelIDs, id)
		}
	}

	return channelIDs, scanner.Err()
}

func forwardToChannels(bot *tgbotapi.BotAPI, message *tgbotapi.Message, channelIDs []int64) {
	for _, id := range channelIDs {
		forward := tgbotapi.NewForward(id, message.Chat.ID, message.MessageID)
		_, err := bot.Send(forward)
		if err != nil {
			log.Printf("Couldn`t send message to channel %d: %v", id, err)
		}
	}
}
