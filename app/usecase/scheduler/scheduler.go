// 定期実行を行う関数をまとめたパッケージ
package scheduler

import (
	"app/controller/log"
	"time"
)

// 型定義
type Job struct {
	Name        string
	Duration    time.Duration
	Function    func() error
	ExecuteFlag bool
}
type Jobs []Job

// 定期実行を行う関数
func schedulerExec(jobs Jobs) {
	for _, job := range jobs {

		// ExecuteFlag が true の場合のみ実行
		if job.ExecuteFlag {
			go func(job Job) {
				for {
					log.Info("  └─ " + job.Name)
					job.Function()
					time.Sleep(job.Duration)
				}
			}(job)
		}
		// Jobs を確実に上から実行するために1秒待機
		time.Sleep(1 * time.Second)
	}
}

// 定期実行を開始する関数
func SchedulerStart() {
	schedulerExec(jobs)
}
