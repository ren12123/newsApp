package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

type Article struct {
	ID    int
	Title string
}

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

func GetUnsummarizedArticles(ctx context.Context, conn *pgx.Conn) ([]Article, error) {
	query := `
	SELECT news_id, title
	FROM articles
	WHERE news_id NOT IN (SELECT news_id FROM summarise)
	LIMIT 5`

	rows, err := conn.Query(ctx, query)
	if err != nil {
		fmt.Printf("取得失敗:%v\n", err)
	}
	defer rows.Close()

	var targets []Article

	for rows.Next() {
		var a Article
		if err := rows.Scan(&a.ID, &a.Title); err != nil {
			return nil, err
		}
		targets = append(targets, a)
	}

	return targets, nil
}

func InsertSummary(ctx context.Context, conn *pgx.Conn, newsID int, content string, aiModel string) {
	_, err := conn.Exec(ctx, `
			  INSERT INTO summaries (news_id, content, ai_model) 
			  VALUES ($1, $2, $3)`, newsID, content, aiModel)
	if err != nil {
		fmt.Printf("要約の保存失敗: %v\n", err)
		return
	}
	fmt.Printf("要約の保存成功")

}
