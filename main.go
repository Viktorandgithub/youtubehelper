package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

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

        flag.Parse()

        client := &http.Client{
                Transport: &transport.APIKey{Key: apiKey},
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

func cleanInput(str string) string{
	lowered := strings.ToLower(str)
	words := lowered
	return words
}
