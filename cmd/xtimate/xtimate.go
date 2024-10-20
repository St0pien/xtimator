package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/St0pien/ratio-xtimator/internal/core"
	"github.com/St0pien/ratio-xtimator/internal/prompts"
	"github.com/St0pien/ratio-xtimator/internal/xtimators"
	"github.com/spf13/cobra"
)

type xtimationResult struct {
	Xtimation core.XtimatedPost
	err       error
}

func Xtimate(cmd *cobra.Command, args []string) {
	content, err := os.ReadFile(args[0])
	if err != nil {
		log.Fatalf("Coudn't read file: %v", err)
	}

	var input core.Thread
	err = json.Unmarshal(content, &input)
	if err != nil {
		log.Fatalf("Failed to parse input file: %v", err)
	}

	output := core.XtimatedThread{
		Start:  input.Start,
		Thread: make([]core.XtimatedPost, 0, len(input.Thread)),
	}

	ch := make(chan xtimationResult, len(input.Thread))
	var wg sync.WaitGroup

	xtimator, err := xtimators.NewOllamaXtimator(xtimators.OllamaXtimatorConfig{
		Model:  "mistral:v0.3",
		Prompt: prompts.BasicPrompt{},
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, post := range input.Thread {
		wg.Add(1)
		go func() {
			xtimation, err := xtimator.Xtimate(input.Start, post)
			ch <- xtimationResult{
				Xtimation: core.XtimatedPost{
					Post: post,
					Type: string(xtimation),
				},
				err: err,
			}

			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for result := range ch {
		output.Thread = append(output.Thread, result.Xtimation)
		log.Printf("%v/%v finished", len(output.Thread), len(input.Thread))
	}

	json, err := json.Marshal(output)
	if err != nil {
		log.Fatal(err)
	}

	os.Stdout.Write(json)
}
