package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func CreateTables(ctx context.Context, conn *pgx.Conn) {
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

	_, err := conn.Exec(ctx, query)
	if err != nil {
		fmt.Println("テーブル作成失敗")
		os.Exit(1)
	}
	fmt.Println("テーブル作成成功")
}

func InsertArticles(ctx context.Context, conn *pgx.Conn, title string, url string) {
	_, err := conn.Exec(ctx, `
		INSERT INTO articles (title, url)
		VALUES ($1, $2)
		ON CONFLICT (url) DO NOTHING`,
		title, url,
	)
	if err != nil {
		fmt.Printf("転送失敗:%v\n", err)
		os.Exit(1)
	}
	fmt.Println("転送成功")

}
