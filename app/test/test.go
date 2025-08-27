package test

import (
	"app/controller/crawler"
	"app/controller/log"
	"app/controller/postgres"

	_ "github.com/lib/pq"
)

func Start() {
	log.Info("テストモード起動")

	// =======================================================================
	// データベース接続とテーブル初期化
	// =======================================================================
	err := postgres.Connect()
	if err != nil {
		return
	}
	err = postgres.InitTable()
	if err != nil {
		return
	}
	log.Info("データベース接続とテーブル初期化完了")

	// =======================================================================
	// テスト内容
	// =======================================================================
	// 初期設定・定数
	targetDomain := "www.city.hamura.tokyo.jp"
	startPath := "/prsite/0000000440.html"
	allowedPaths := []string{
		"/prsite/0000000440.html",
	}
	maxScrapeDepth := 7
	isTest := true

	// クロールを開始
	err = crawler.CrawlDomain(targetDomain, startPath, allowedPaths, maxScrapeDepth, isTest)
	if err != nil {
		log.Error(err)
		return
	}
}
