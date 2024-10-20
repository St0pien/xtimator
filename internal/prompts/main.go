package prompts

import "github.com/St0pien/ratio-xtimator/internal/core"

type Prompt interface {
	GetSystemPrompt() string
	Build(tweet core.Post, response core.Post) (string, error)
}
