// PostgreSQL を利用するための関数をまとめたパッケージ
package postgres

import (
	"app/controller/log"
	"crypto/sha1"
	"encoding/hex"
	"net/url"
)

/*
DB にクロールしたデータを保存する関数
  - url			クロールしたページの URL
  - title			ページのタイトル
  - description	ページの説明
  - keywords		ページのキーワード
  - markdown		ページのマークダウンコンテンツ
  - return) err	エラー
*/
func SaveCrawledData(urlStr, title, description, keywords, markdown string) (err error) {
	// URLをパース
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		log.Error(err)
		return err
	}
	domain := parsedURL.Host
	path := parsedURL.Path

	// markdownのハッシュを計算
	hash := sha1.Sum([]byte(markdown))
	hashStr := hex.EncodeToString(hash[:])

	// すでに存在するかどうかをドメインとパスで確認
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM pages WHERE domain = $1 AND path = $2)", domain, path).Scan(&exists)
	if err != nil {
		log.Error(err)
		return err
	}

	// すでに存在する場合は更新、存在しない場合は新規挿入
	if exists {
		updateSQL := `
		UPDATE pages
		SET title = $1, description = $2, keywords = $3, markdown = $4, hash = $5, updated_at = CURRENT_TIMESTAMP
		WHERE domain = $6 AND path = $7;
		`
		_, err = db.Exec(updateSQL, title, description, keywords, markdown, hashStr, domain, path)
		if err != nil {
			log.Error(err)
			return err
		}
		return nil
	} else {
		insertSQL := `
		INSERT INTO pages (domain, path, title, description, keywords, markdown, created_at, updated_at, hash)
		VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, $7)
		`
		_, err = db.Exec(insertSQL, domain, path, title, description, keywords, markdown, hashStr)
		if err != nil {
			log.Error(err)
			return err
		}
		return nil
	}
}
