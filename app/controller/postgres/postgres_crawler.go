// PostgreSQL を利用するための関数をまとめたパッケージ
package postgres

import (
	"app/controller/log"
	"app/domain/model"
	"context"
	"unicode/utf8"
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
	// 文字列の長さが制限を超えている場合は UTF-8 安全に切り詰める
	domain = truncateRunes(domain, 100)
	path = truncateRunes(path, 255)
	title = truncateRunes(title, 100)
	description = truncateRunes(description, 255)
	keywords = truncateRunes(keywords, 255)

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

// UTF-8 安全に文字数で切り詰めるヘルパー関数
func truncateRunes(s string, max int) string {
	if utf8.RuneCountInString(s) <= max {
		return s
	}
	r := []rune(s)
	return string(r[:max])
}
