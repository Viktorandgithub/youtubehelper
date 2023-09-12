package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/PaulSonOfLars/gotgbot/handlers/Filters"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

type config struct {
    ApiKey        string
    TelegramToken string
}

var (
    query      = flag.String("query", "Google", "Search term")
    maxResults = flag.Int64("max-results", 5, "Max YouTube results")
)

func main() {
    godotenv.Load(".env")

    apiKey := os.Getenv("YOUR_API_KEY")
    if apiKey == "" {
        log.Fatal("YOUR_API_KEY environment variable is not set")
    }

    telegramToken := os.Getenv("TELEGRAM_BOT_TOKEN")
    if telegramToken == "" {
        log.Fatal("TELEGRAM_BOT_TOKEN environment variable is not set")
    }

    cfg := config{
        ApiKey:        apiKey,
        TelegramToken: telegramToken,
    }

    logCfg := zap.NewProductionEncoderConfig()
    logCfg.EncodeLevel = zapcore.CapitalLevelEncoder
    logCfg.EncodeTime = zapcore.RFC3339TimeEncoder

    logger := zap.New(zapcore.NewCore(zapcore.NewConsoleEncoder(logCfg), os.Stdout, zap.InfoLevel))

    updater, err := gotgbot.NewUpdater(logger, cfg.TelegramToken)
    if err != nil {
        logger.Panic("UPDATER FAILED TO START")
        return
    }
    logger.Sugar().Info("UPDATER STARTED SUCCESSFULLY")
    updater.StartCleanPolling()

    updater.Dispatcher.AddHandler(handlers.NewMessage(Filters.Text, func(b ext.Bot, u *gotgbot.Update) error {
        handleTelegramInput(b, u, &cfg)
        return nil
    }))

    updater.Idle()

    handleYouTubeSearch(&cfg)
}

func handleTelegramInput(b ext.Bot, u *gotgbot.Update, cfg *config) {
    // Handle user input from Telegram
    userMessage := u.EffectiveMessage.Text
    if strings.ToLower(userMessage) == "search" {
        // Ask the user for the search query
        b.SendMessage(u.EffectiveChat.Id, "Enter the search query:")
    } else {
        // Perform the YouTube search
        query := cleanInput(userMessage)
        videos := handleYoutube(cfg, &query)
        videoLink := printIDs("Videos", videos)
        b.SendMessage(u.EffectiveChat.Id, videoLink)
    }
}

func handleYouTubeSearch(cfg *config) {
    for {
        // Ask the user for the creator's name
        fmt.Print("Enter the name of the YouTube creator: ")
        scanner := bufio.NewScanner(os.Stdin)
        scanner.Scan()
        creatorName := scanner.Text()
        cleanedCreatorName := cleanInput(creatorName)
        if len(cleanedCreatorName) == 0 {
            continue
        }
        query = &cleanedCreatorName

        videos := handleYoutube(cfg, query)

        printIDs("Videos", videos)
    }
}

func cleanInput(str string) string {
    lowered := strings.ToLower(str)
    words := lowered
    return words
}

func printIDs(sectionName string, matches map[string]string) string{
    fmt.Printf("%v:\n", sectionName)
    videoLink := "https://www.youtube.com/watch?v"
    for id := range matches {
        videoLink = fmt.Sprintf("https://www.youtube.com/watch?v=%s", id)
        fmt.Printf("First video on the channel: %s\n", videoLink)
    }
    fmt.Printf("\n\n")
    return videoLink
}

func handleYoutube(cfg *config, query *string) map[string]string {
    flag.Parse()

    client := &http.Client{
        Transport: &transport.APIKey{Key: cfg.ApiKey},
    }

    service, err := youtube.New(client)
    if err != nil {
        log.Fatalf("Error creating new YouTube client: %v", err)
    }
    part := []string{"id", "snippet"}
    // Make the API call to YouTube.
    call := service.Search.List(part).
        Q(*query).
        MaxResults(*maxResults)
    response, err := call.Do()
    if err != nil {
        fmt.Println(err)
    }

    // Group video, channel, and playlist results in separate lists.
    videos := make(map[string]string)

    for _, item := range response.Items {
        if item.Id.Kind == "youtube#video" {
            videos[item.Id.VideoId] = item.Snippet.Title
            break
        }
    }

    return videos
}
