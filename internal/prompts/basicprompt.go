package prompts

import (
	"fmt"

	"github.com/St0pien/ratio-xtimator/internal/core"
)

type BasicPrompt struct {
}

func (prompt BasicPrompt) Build(tweet core.Post, response core.Post) (string, error) {
	return fmt.Sprintf(`Estimate category of response:
%v wrote "%v" he got %v likes.

%v responded: "%v" and got %v likes

What type of response %v wrote? Reply only with one of possible options: positive, negative, joke/irony, spam/ad/scam`, tweet.Author.Name, tweet.Text, tweet.Likes, response.Author.Name, response.Text, response.Likes, response.Author.Name), nil
}

func (prompt BasicPrompt) GetSystemPrompt() string {
	return "You are an certified social media expert with experience in judging the overall connotation of internet posts. Your job is to estimate in which category certain responses fit into. Categories are: positive, negative, joke/irony, spam/ad/scam. You will receive a tweet and a response to this tweet. You will have to estimate the category of the tweet. Respond only with the name of category"
}
