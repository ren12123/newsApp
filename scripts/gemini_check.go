package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()

	apikey := os.Getenv("Gemini_API_KEY")
	if apikey == "" {
		log.Fatal("API_KEYが見つからない")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apikey))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("models/gemini-2.5-flash")

	text := "今日のトップページにある２つの記事"
	prompt := fmt.Sprintf("%s\nを15文字程度で要約して", text)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Fatal(err)
	}

	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			fmt.Println("AIの要約結果", cand.Content.Parts[0])
		}
	}

}
