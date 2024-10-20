package parsers

import (
	"encoding/json"
	"fmt"

	"github.com/St0pien/ratio-xtimator/internal/core"
)

type AuthorResponse struct {
	Name         string `json:"name"`
	ScreenName   string `json:"screen_name"`
	Image        string `json:"image"`
	BlueVerified bool   `json:"blue_verified"`
}

type PostResponse struct {
	Likes  int            `json:"likes"`
	Text   string         `json:"text"`
	Id     string         `json:"id"`
	Author AuthorResponse `json:"author"`
}

type InitialResponse struct {
	Likes   int            `json:"likes"`
	Replies int            `json:"replies"`
	Text    string         `json:"text"`
	Author  AuthorResponse `json:"author"`
	Id      string         `json:"id"`
	Thread  []PostResponse `json:"thread"`
	Cursor  string         `json:"cursor"`
}

type FollowupResponse struct {
	Id     string         `json:"id"`
	Thread []PostResponse `json:"thread"`
	Cursor string         `json:"cursor"`
}

type XRapidApiParser struct {
	Parsed core.Thread
}

func transformAuthor(a AuthorResponse) core.Author {
	return core.Author{
		Name:         a.Name,
		Handle:       a.ScreenName,
		Image:        a.Image,
		BlueVerified: a.BlueVerified,
	}
}

func transformPost(p PostResponse) core.Post {
	return core.Post{
		Id:     p.Id,
		Author: transformAuthor(p.Author),
		Text:   p.Text,
		Likes:  p.Likes,
	}
}

func (parser *XRapidApiParser) ParseInitialSegment(content []byte) (string, error) {
	var data InitialResponse

	err := json.Unmarshal(content, &data)

	if err != nil {
		return "", fmt.Errorf("[XRapidApiParser] Error: %w", err)
	}

	parser.Parsed = core.Thread{
		Start: core.Post{
			Id:     data.Id,
			Author: transformAuthor(data.Author),
			Text:   data.Text,
			Likes:  data.Likes,
		},
		Thread: make([]core.Post, 0, data.Replies),
	}

	for _, post := range data.Thread {
		parser.Parsed.Thread = append(parser.Parsed.Thread, transformPost(post))
	}

	return data.Cursor, nil
}

func (parser *XRapidApiParser) ParseFollowupSegments(content []byte) (string, error) {
	var data FollowupResponse

	err := json.Unmarshal(content, &data)
	if err != nil {
		return "", fmt.Errorf("[XRadpidApiParser] Error: %w", err)
	}

	for _, post := range data.Thread {
		parser.Parsed.Thread = append(parser.Parsed.Thread, transformPost(post))
	}

	return data.Cursor, nil
}
