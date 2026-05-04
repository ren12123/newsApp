package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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
		description TEXT,
		fetched_at TIMESTAMP DEFAULT NOW(),
		published_at TIMESTAMP,
		category_id INT REFERENCES categories(category_id),
		url TEXT UNIQUE NOT NULL
	);

	CREATE TABLE IF NOT EXISTS summaries (
		summary_id SERIAL PRIMARY KEY,
		news_id INT UNIQUE REFERENCES articles(news_id) ON DELETE CASCADE,
		content TEXT NOT NULL,
		ai_model TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS users (
		user_id SERIAL PRIMARY KEY,
		name TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS categories (
    	category_id SERIAL PRIMARY KEY,
    	name TEXT NOT NULL UNIQUE
	);

	CREATE TABLE IF NOT EXISTS user_categories (
		user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
		category_id INT,
		PRIMARY KEY (user_id, category_id)
	);`

	_, err := conn.Exec(ctx, query)
	if err != nil {
		fmt.Println("テーブル作成失敗")
		os.Exit(1)
	}
	fmt.Println("テーブル作成成功")
}

func InsertArticles(ctx context.Context, conn *pgx.Conn, title string, url string, description string, published_at time.Time, category_ID int) {
	_, err := conn.Exec(ctx, `
		INSERT INTO articles (title, url, description, published_at, categoryID)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (url) DO NOTHING`,
		title, url, description, published_at, category_ID,
	)
	if err != nil {
		fmt.Printf("転送失敗:%v\n", err)
		log.Printf("記事保存失敗:%v\n", err)
		return
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
