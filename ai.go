package main

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type AISummarizer struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

func NewAISummarizer(ctx context.Context, apikey string) (*AISummarizer, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apikey))
	if err != nil {
		return nil, err
	}

	model := client.GenerativeModel("models/gemini-2.5-flash")

	return &AISummarizer{
		client: client,
		model:  model,
	}, nil
}

func (s *AISummarizer) GetSummary(ctx context.Context, title string) (string, error) {
	prompt := fmt.Sprintf("以下のニュース記事を15文字程度で短く要約してください。挨拶は不要です:\n%s", title)

	resp, err := s.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return "", fmt.Errorf("AIからの返答無")
	}

	summary := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])
	return summary, nil

}
