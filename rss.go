package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

func Getnewsforhttp(ctx context.Context, url string) ([]ITEM, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("データ取得失敗:%v\n", err)
		return nil, err
	}
	fmt.Println("データ取得")
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("読み込み失敗:%v\n", err)
		return nil, err
	}
	fmt.Println("読み込み成功")

	var rss RSS
	err = xml.Unmarshal(body, &rss) //rssだけだったらコピーが渡されるだけで意味がない
	if err != nil {
		return nil, err
	}

	return rss.Items, nil
}
