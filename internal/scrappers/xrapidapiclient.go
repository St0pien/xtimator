package scrappers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const baseUrl string = "https://twitter-api45.p.rapidapi.com/tweet_thread.php"
const rapidApiHost string = "twitter-api45.p.rapidapi.com"

type ThreadSegment struct {
	TweetId string
	Cursor  string
}

type XRapidApiClient struct {
	ApiKey string
	cache  ApiCache
}

func NewXRapidApiClient(apiKey string) XRapidApiClient {
	cache, err := NewApiCache(".cache", "x_com")
	if err != nil {
		log.Fatal(err)
	}

	return XRapidApiClient{
		ApiKey: apiKey,
		cache:  cache,
	}
}

func (client *XRapidApiClient) fetch(segment ThreadSegment) ([]byte, error) {
	url := fmt.Sprintf("%s?id=%s", baseUrl, segment.TweetId)

	if segment.Cursor != "" {
		url += "&cursor=" + segment.Cursor
	}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, fmt.Errorf("Error while fetching: %w", err)
	}

	req.Header.Add("x-rapidapi-key", client.ApiKey)
	req.Header.Add("x-rapidapi-host", rapidApiHost)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("Error while fetching: %w", err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	return body, nil

}

func (client *XRapidApiClient) FetchSegment(segment ThreadSegment) ([]byte, error) {
	log.Printf("[FetchFromXRapidAPI] Checking cache for %q", segment)

	content, err := client.cache.Get(segment.TweetId + segment.Cursor)
	if err != nil {
		log.Printf("[FetchFromXRapidAPI] Cache error: %v", err)
		err = nil
	}

	if len(content) != 0 {
		log.Printf("[FetchFromXRapidAPI] Cache hit %q", segment)
		return content, nil
	}

	log.Printf("[FetchFromXRapidAPI] cache miss, fetching %q", segment)

	for i := 0; i < 10; i++ {
		content, err = client.fetch(segment)

		if err != nil {
			return nil, fmt.Errorf("[XRapidApiClient] Error while fetching: %w", err)
		}

		if len(content) == 0 {
			log.Printf("[XRapidApiClient] Rate limit hit - pausing for %v seconds", 30*i)
			time.Sleep(time.Second * 30 * time.Duration(i))
		} else {
			break
		}
	}

	log.Printf("[XRapidApiClient] Fetch success %v", segment)

	if client.cache.Save(segment.TweetId+segment.Cursor, content) != nil {
		log.Printf("[FetchFromXRapidAPI] Error : %v", err)
	}

	return content, err

}
