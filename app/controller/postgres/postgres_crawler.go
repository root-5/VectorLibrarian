// PostgreSQL を利用するための関数をまとめたパッケージ
package postgres

import (
	"app/controller/log"
	"app/controller/model"
	"crypto/sha1"
	"encoding/hex"
	"net/url"
	"strings"
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
		// 更新処理
		// カラム名をmodelから取得
		columnNames := model.GetColumnNames(&model.PageContentAtUpdate{})
		columnNamesStr := strings.Join(columnNames, ", ")
		placeholders := model.GetPlaceholders(columnNames)
		placeholdersStr := strings.Join(placeholders, ", ")
		whereTargets := []interface{}{model.PageContentAtUpdate{}.Domain, model.PageContentAtUpdate{}.Path}
		whereTargetsStr := model.GetWhereTargetStr(columnNames, placeholders, whereTargets)

		updateSQL := `
		UPDATE pages (` + columnNamesStr + `)
		VALUES (` + placeholdersStr + `)
		WHERE ` + whereTargetsStr + `;
		`
		// updateSQL := `
		// UPDATE pages (` + columnNamesStr + `)
		// VALUES (` + placeholdersStr + `)
		// WHERE domain = $5 AND path = $6;
		// `
		_, err = db.Exec(updateSQL, domain, path, title, description, keywords, markdown, hashStr)
		if err != nil {
			log.Error(err)
			return err
		}
		return nil
	} else {
		// 存在しない場合は新規挿入
		columnNames := model.GetColumnNames(&model.PageContentAtInsert{})
		columnNamesStr := strings.Join(columnNames, ", ")
		placeholders := model.GetPlaceholders(columnNames)
		placeholdersStr := strings.Join(placeholders, ", ")

		insertSQL := `
		INSERT INTO pages (` + columnNamesStr + `)
		VALUES (` + placeholdersStr + `)
		`
		_, err = db.Exec(insertSQL, domain, path, title, description, keywords, markdown, hashStr)
		if err != nil {
			log.Error(err)
			return err
		}
		return nil
	}
}
