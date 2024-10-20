package test

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"testing"

	"github.com/St0pien/ratio-xtimator/internal/core"
	"github.com/St0pien/ratio-xtimator/internal/prompts"
	"github.com/St0pien/ratio-xtimator/internal/xtimators"
	"github.com/stretchr/testify/assert"
)

type TestCases struct {
	Start    core.Post `json:"start"`
	Comments []struct {
		Post       core.Post                 `json:"post"`
		Impression xtimators.TweetImpression `json:"impression"`
	} `json:"comments"`
}

type TestResult struct {
	Input    core.Post
	Expected xtimators.TweetImpression
	Actual   xtimators.TweetImpression
	err      error
}

func loadTestData() TestCases {
	content, err := os.ReadFile("data/1.json")
	if err != nil {
		log.Fatal(err)
	}

	var cases TestCases

	err = json.Unmarshal(content, &cases)
	if err != nil {
		log.Fatal(err)
	}

	return cases
}

func TestPrompt(t *testing.T) {
	assert := assert.New(t)
	tests := loadTestData()

	xtimator, err := xtimators.NewOllamaXtimator(xtimators.OllamaXtimatorConfig{
		Model:  "mistral:v0.3",
		Prompt: prompts.BasicPrompt{},
	})
	if err != nil {
		log.Fatal(err)
	}

	ch := make(chan TestResult, len(tests.Comments))
	var wg sync.WaitGroup

	for _, test := range tests.Comments {
		wg.Add(1)
		go func() {
			actual, err := xtimator.Xtimate(tests.Start, test.Post)
			ch <- TestResult{
				Input:    test.Post,
				Expected: test.Impression,
				Actual:   actual,
				err:      err,
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for result := range ch {
		assert.NoError(result.err)
		assert.Equal(result.Expected, result.Actual, result.Input)
	}
}
