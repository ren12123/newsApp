package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

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

}
