// PostgreSQL を利用するための関数をまとめたパッケージ
package postgres

import (
	"app/controller/log"
	"app/domain/model"
	"context"
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
func SaveCrawledData(domain, path, title, description, keywords, markdown, hash string) (err error) {
	// 文字列の長さが制限を超えている場合は切り詰める
	if len(domain) > 100 {
		domain = domain[:100]
	}
	if len(path) > 255 {
		path = path[:255]
	}
	if len(title) > 100 {
		title = title[:100]
	}
	if len(description) > 255 {
		description = description[:255]
	}
	if len(keywords) > 255 {
		keywords = keywords[:255]
	}

	// pages テーブルのモデルを作成
	page := &model.Page{
		Domain:      domain,
		Path:        path,
		Title:       title,
		Description: description,
		Keywords:    keywords,
		Markdown:    markdown,
		Hash:        hash,
	}

	// 同 URL の存在確認
	exists, err := db.NewSelect().
		Model(page).
		Where("domain = ? AND path = ?", domain, path).
		Exists(context.Background())
	if err != nil {
		log.Error(err)
		return err
	}

	// すでに存在する場合は更新、存在しない場合は新規挿入
	if exists {
		// 更新
		_, err = db.NewUpdate().
			Model(page).
			Where("domain = ? AND path = ?", domain, path).
			Set("title = ?", title).
			Set("description = ?", description).
			Set("keywords = ?", keywords).
			Set("markdown = ?", markdown).
			Set("hash = ?", hash).
			Set("updated_at = CURRENT_TIMESTAMP").
			Exec(context.Background())
	} else {
		// 新規挿入
		_, err = db.NewInsert().Model(page).Exec(context.Background())
	}
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
