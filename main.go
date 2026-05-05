package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()
	connStr := "postgres://user:password@localhost:5432/mydatabase?sslmode=disable"
	modePtr := flag.String("mode", "all", "実行モード (fetch, summarize, all)")
	flag.Parse()
	mode := *modePtr

	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		fmt.Printf("DB接続失敗:%v\n", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	if mode == "fetch" || mode == "all" {
		url := "https://news.google.com/rss?hl=ja&gl=JP&ceid=JP:ja"

		articles, err := Getnewsforhttp(ctx, url)
		if err != nil {
			fmt.Printf("受け取り失敗:%v\n", err)
			return
		}

		for _, item := range articles {
			InsertArticles(ctx, conn, item.Title, item.Link)
		}
	}

	if mode == "summarize" || mode == "all" {
		unsummarize, err := GetUnsummarizedArticles(ctx, conn)
		if err != nil {
			fmt.Printf("受け取り失敗:%v\n", err)
			return
		}

		s, err := NewAISummarizer(ctx, os.Getenv("Gemini_API_KEY"))
		if err != nil {
			fmt.Printf("apikeyが正しくありません:%v\n", err)
			os.Exit(1)
		}

		for _, a := range unsummarize {
			result, err := s.GetSummary(ctx, a.Title)
			if err != nil {
				continue
			}

			InsertSummary(ctx, conn, a.ID, result, "gemini-2.5-flash")
		}
	}
}
