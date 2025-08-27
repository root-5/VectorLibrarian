package main

import (
	"app/controller/api"
	"app/controller/log"
	"app/controller/postgres"
	"app/test"
	"app/usecase/scheduler"
	"flag"

	_ "github.com/lib/pq"
)

func main() {
	// flag パッケージを使ってモードを指定できるようにする
	mode := flag.String("mode", "normal", "execution mode: normal or test")
	flag.Parse()

	switch *mode {
	case "test":
		// -mode=test を指定した場合の処理
		runTestMode()
	default:
		run()
	}
}

func runTestMode() {
	log.Info("テストモード起動")
	test.Start()
}

func run() {
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
	// スケジューラーを起動
	// =======================================================================
	log.Info("スケジューラー起動")
	scheduler.SchedulerStart()

	// =======================================================================
	// API サーバーを起動
	// =======================================================================
	log.Info("API サーバー起動")
	api.StartServer()
}
