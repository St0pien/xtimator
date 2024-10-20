package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/St0pien/ratio-xtimator/internal/parsers"
	"github.com/St0pien/ratio-xtimator/internal/scrappers"
	"github.com/spf13/cobra"
)

func PreProcess(cmd *cobra.Command, args []string) {
	api := scrappers.NewXRapidApiClient(cmd.Flag("api-key").Value.String())
	tweetId := args[0]

	content, err := api.FetchSegment(scrappers.ThreadSegment{
		TweetId: tweetId,
	})
	if err != nil {
		log.Fatalf("Error %v", err)
	}

	var apiParser parsers.XRapidApiParser

	cursor, err := apiParser.ParseInitialSegment(content)
	if err != nil {
		log.Fatalf("Error %v", err)
	}

	for cursor != "" {
		content, err = api.FetchSegment(scrappers.ThreadSegment{
			TweetId: tweetId,
			Cursor:  cursor,
		})
		if err != nil {
			log.Fatalf("Error %v", err)
		}

		cursor, err = apiParser.ParseFollowupSegments(content)
		if err != nil {
			log.Fatalf("Error %v", err)
		}
	}

	encoded, _ := json.Marshal(apiParser.Parsed)
	err = os.WriteFile("output.json", encoded, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
