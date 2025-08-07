// PostgreSQL を利用するための関数をまとめたパッケージ
package postgres

import (
	"app/controller/log"
	"app/domain/model"
	"context"
)

/*
同一のハッシュ値を持つページが存在するか確認する関数
  - hash			確認するページのハッシュ値
  - return) exists	同一のハッシュ値を持つページが存在するかどうか
  - return) err		エラー
*/
func CheckHashExists(hash string) (exists bool, err error) {
	// ハッシュ値を持つページが存在するか確認
	page := &model.Page{Hash: hash}
	exists, err = db.NewSelect().
		Model(page).
		Where("hash = ?", hash).
		Exists(context.Background())
	if err != nil {
		log.Error(err)
		return false, err
	}

	return exists, nil
}

/*
クロールしたページデータを保存する関数
  - url			クロールしたページの URL
  - title			ページのタイトル
  - description	ページの説明
  - keywords		ページのキーワード
  - markdown		ページのマークダウンコンテンツ
  - return) err	エラー
*/
func SaveCrawledData(domain, path, title, description, keywords, markdown, hash string, vector []float32) (err error) {
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
		Vector:      vector,
	}

	// すでに存在する場合は更新、存在しない場合は新規挿入
	_, err = db.NewInsert().
		Model(page).
		On("CONFLICT (domain, path) DO UPDATE").
		Set("title = ?", title).
		Set("description = ?", description).
		Set("keywords = ?", keywords).
		Set("markdown = ?", markdown).
		Set("hash = ?", hash).
		Set("vector = ?", vector).
		Set("updated_at = CURRENT_TIMESTAMP").
		Exec(context.Background())
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
