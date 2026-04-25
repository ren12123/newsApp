package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
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

func main() {
	ctx := context.Background()
	connStr := "postgres://user:password@localhost:5432/mydatabase?sslmode=disable"

	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		fmt.Printf("DB接続失敗:%v\n", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	query := `
	CREATE TABLE IF NOT EXISTS articles (
		news_id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		date TIMESTAMP DEFAULT NOW(),
		category_id INT,
		url TEXT UNIQUE
	);

	CREATE TABLE IF NOT EXISTS summaries (
		summary_id SERIAL PRIMARY KEY,
		news_id INT,
		content TEXT,
		ai_model TEXT
	);

	CREATE TABLE IF NOT EXISTS users (
		user_id SERIAL PRIMARY KEY,
		name TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS user_categories (
		user_id INT,
		category_id INT
	);`

	_, err = conn.Exec(ctx, query)
	if err != nil {
		fmt.Printf("テーブル作成失敗: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("ニュースアプリのデータベース基盤が完成しました！")

	url := "https://news.google.com/rss?hl=ja&gl=JP&ceid=JP:ja"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("データ取得失敗:%v\n", err)
		return
	}
	fmt.Println("データ取得")
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("読み込み失敗:%v\n", err)
		return
	}
	fmt.Println("読み込み成功")

	var rss RSS
	err = xml.Unmarshal(body, &rss) //rssだけだったらコピーが渡されるだけで意味がない
	if err != nil {
		fmt.Printf("解析失敗: %v\n", err)
		return
	}
	fmt.Println("解析成功")

	for _, item := range rss.Items {
		_, err := conn.Exec(ctx, `
		INSERT INTO articles (title, url)
		VALUES ($1, $2)
		ON CONFLICT (url) DO NOTHING`,
			item.Title, item.Link,
		)
		if err != nil {
			fmt.Printf("転送失敗:%v\n", err)
			os.Exit(1)
		}
	}
	fmt.Println("転送成功")
}
