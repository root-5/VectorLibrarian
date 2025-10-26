package scheduler

import (
	"app/controller/crawler"
	"time"
)

// 定期実行する関数とその設定をまとめた構造体
var jobs = Jobs{
	{
		Name:        "クローリング開始",
		Duration:    7 * 24 * time.Hour, // 7日ごとに実行
		Function:    crawler.Start,      // クローリングを開始する関数
		ExecuteFlag: true,
	},
}
