package xtimators

import "github.com/St0pien/ratio-xtimator/internal/core"

const (
	Positive = "positive"
	Negative = "negative"
	Joke     = "joke/irony"
	Spam     = "spam/ad/scam"
)

type TweetImpression string

type Xtimator interface {
	Xtimate(core.Post, core.Post) (TweetImpression, error)
}
