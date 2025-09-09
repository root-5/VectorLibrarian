// PostgreSQL を利用するための関数をまとめたパッケージ
package postgres

import (
	"context"
	"database/sql"
	"os"

	"app/controller/log"
	"app/usecase/entity"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

// 型定義
var db *bun.DB

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
	port := "5432"
	sslmode := "disable"
	dsn := "postgres://" + user + ":" + password + "@" + host + ":" + port + "/" + dbname + "?sslmode=" + sslmode

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db = bun.NewDB(sqldb, pgdialect.New())
	return nil
}

/*
DB の初期化をする関数
  - return) err	エラー
*/
func InitTable() (err error) {
	_, err = db.NewCreateTable().
		Model((*entity.DBDomain)(nil)).
		IfNotExists().
		Exec(context.Background())
	if err != nil {
		log.Error(err)
		return
	}

	_, err = db.NewCreateTable().
		Model((*entity.DBPage)(nil)).
		IfNotExists().
		Exec(context.Background())
	if err != nil {
		log.Error(err)
		return
	}

	_, err = db.NewCreateTable().
		Model((*entity.DBChunk)(nil)).
		IfNotExists().
		Exec(context.Background())
	if err != nil {
		log.Error(err)
		return
	}

	_, err = db.NewCreateTable().
		Model((*entity.DBVector)(nil)).
		IfNotExists().
		Exec(context.Background())
	if err != nil {
		log.Error(err)
		return
	}

	_, err = db.NewCreateTable().
		Model((*entity.DBNlpConfig)(nil)).
		IfNotExists().
		Exec(context.Background())
	if err != nil {
		log.Error(err)
		return
	}

	return nil
}
