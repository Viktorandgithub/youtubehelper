package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)


var (
    query      = flag.String("query", "Google", "Search term")
    maxResults = flag.Int64("max-results", 5, "Max YouTube results")
)

func handleYoutube(cfg *config, query *string) map[string]string{
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



func cleanInput(str string) string{
	lowered := strings.ToLower(str)
	words := lowered
	return words
}