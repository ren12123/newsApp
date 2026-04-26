package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

type ITEM struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
}

type RSS struct {
	Items []ITEM `xml:"channel>item"` //channelの中にあるアイテムをすべて取る
}

type SummarizeAI interface {
	GetSummary(text string) (string, error)
}

func main() {
	ctx := context.Background()
	connStr := "postgres://user:password@localhost:5432/mydatabase?sslmode=disable"

	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		fmt.Printf("DB接続失敗:%v\n", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

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
