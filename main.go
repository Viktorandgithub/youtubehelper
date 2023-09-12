package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	ApiKey string
}

func main() {

godotenv.Load(".env")

apiKey := os.Getenv("YOUR_API_KEY")
if apiKey == "" {
    log.Fatal("YOUR_API_KEY environment variable is not set")
}
cfg := config{
    ApiKey : apiKey,    

}
scanner := bufio.NewScanner(os.Stdin)

    for {
        // Ask the user for the creator's name
        fmt.Print("Enter the name of the YouTube creator: ")
        
        scanner.Scan()
        creatorName := scanner.Text()
        cleanedCreatorName := cleanInput(creatorName)
        if len(cleanedCreatorName) == 0 {
            continue
        }
        query = &cleanedCreatorName

        videos := handleYoutube(&cfg, query)

        printIDs("Videos", videos)

        
    }
}

func printIDs(sectionName string, matches map[string]string) {
    fmt.Printf("%v:\n", sectionName)
    for id := range matches {
            videoLink := fmt.Sprintf("https://www.youtube.com/watch?v=%s", id,)
            fmt.Printf("First video on the channel: %s\n", videoLink)
    }
    fmt.Printf("\n\n")
}