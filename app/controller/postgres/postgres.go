// PostgreSQL を利用するための関数をまとめたパッケージ
package postgres

import (
	"database/sql"
	"os"

	"app/controller/log"

	_ "github.com/lib/pq"
)

// 型定義
var db *sql.DB

/*
DB の接続をする関数
  - return) err	エラー
*/
func Connect() (err error) {
	// 環境変数から接続情報を取得
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=5432" + " sslmode=disable TimeZone=Asia/Tokyo"

	// DB に接続
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Error(err)
		return err
	}

	// DB の接続を確認
	err = db.Ping()
	if err != nil {
		log.Error(err)
		return
	}

	return nil
}

/*
DB の初期化をする関数
  - return) err	エラー
*/
func InitTable() (err error) {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS pages (
		id SERIAL PRIMARY KEY,
		domain TEXT NOT NULL,
		path TEXT NOT NULL,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		keywords TEXT NOT NULL,
		markdown TEXT NOT NULL,
		hash TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMP WITH TIME ZONE,
		UNIQUE (domain, path)
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Error(err)
		return
	}

	return nil
}
