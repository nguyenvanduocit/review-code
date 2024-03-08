package gemini

import (
	"context"
	"github.com/google/generative-ai-go/genai"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
	"os"
)

func NewGeminiClient(ctx context.Context) (*genai.GenerativeModel, func(), error) {
	apiKey := os.Getenv("VERTEX_API_KEY")
	if apiKey == "" {
		return nil, nil, errors.New("missing VERTEX_API_KEY")
	}
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, nil, errors.WithMessage(err, "failed to create genai client")
	}

	cleanup := func() {
		if err := client.Close(); err != nil {
			panic(err)
		}
	}

	model := client.GenerativeModel("gemini-1.0-pro-latest")
	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: 0,
		},
	}

	candidateCount := int32(1)
	model.CandidateCount = &candidateCount

	return model, cleanup, nil
}
