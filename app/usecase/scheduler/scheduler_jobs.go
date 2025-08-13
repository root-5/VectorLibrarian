package scheduler

import (
	"app/controller/crawler"
	"time"
)

// 定期実行する関数とその設定をまとめた構造体
var jobs = Jobs{
	{
		Name:        "crawler.Start",
		Duration:    6 * time.Hour,
		Function:    crawler.Start, // クローリングを開始する関数
		ExecuteFlag: true,
	},
}
