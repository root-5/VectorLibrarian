// PostgreSQL を利用するための関数をまとめたパッケージ
package postgres

import (
	"app/controller/log"
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
func SaveCrawledData(url, title, description, keywords, markdown string) (err error) {
	insertSQL := `
		INSERT INTO pages (url, title, description, keywords, markdown_content)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (url) DO UPDATE SET
			title = EXCLUDED.title,
			description = EXCLUDED.description,
			keywords = EXCLUDED.keywords,
			markdown_content = EXCLUDED.markdown_content,
			created_at = CURRENT_TIMESTAMP;`

	_, err = db.Exec(insertSQL, url, title, description, keywords, markdown)
	if err != nil {
		log.Error(err)
	}
	return err
}
