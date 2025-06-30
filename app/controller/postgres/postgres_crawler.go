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

	insertSQL := `
		INSERT INTO pages (domain, path, title, description, keywords, markdown, hash, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP)
		ON CONFLICT (domain, path) DO UPDATE SET
			title = EXCLUDED.title,
			description = EXCLUDED.description,
			keywords = EXCLUDED.keywords,
			markdown = EXCLUDED.markdown,
			hash = EXCLUDED.hash,
			updated_at = CURRENT_TIMESTAMP;`

	_, err = db.Exec(insertSQL, domain, path, title, description, keywords, markdown, hashStr)
	if err != nil {
		log.Error(err)
	}
	return err
}
