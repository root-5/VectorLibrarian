// PostgreSQL を利用するための関数をまとめたパッケージ
package postgres

import (
	"app/controller/log"
	"app/controller/nlp"
	"app/domain/model"

	"app/usecase/entity"
	"context"
	"unicode/utf8"
)

/*
ドメイン情報を取得する関数
  - return) domains	ドメイン情報のスライス
  - return) err		エラー
*/
func GetDomains() (dbDomains []entity.DBDomain, err error) {
	// ドメイン情報を取得
	err = db.NewSelect().
		Model(&dbDomains).
		Scan(context.Background())
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return dbDomains, nil
}

/*
同一のハッシュ値を持つページが存在するか確認する関数
  - hash			確認するページのハッシュ値
  - return) exists	同一のハッシュ値を持つページが存在するかどうか
  - return) err		エラー
*/
func CheckHashExists(hash string) (exists bool, err error) {
	// ハッシュ値を持つページが存在するか確認
	page := &entity.DBPage{PageInfo: model.PageInfo{Hash: hash}}
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
  - pageInfo		保存するページ情報
  - convertResult	nlp サーバーからの変換結果
  - return) err	エラー
*/
func SaveCrawledData(page model.PageInfo, convertResult nlp.ConvertResponse) (err error) {
	// ページ情報
	// 文字列の長さが制限を超えている場合は UTF-8 安全に切り詰める
	page.Path = truncateRunes(page.Path, 255)
	page.Title = truncateRunes(page.Title, 100)
	page.Description = truncateRunes(page.Description, 255)
	page.Keywords = truncateRunes(page.Keywords, 255)
	log.Info(">> Saving:" + page.Title)

	// チャンク情報、ベクトル情報、NLP設定情報
	chunks := convertResult.Chunks
	vectors := convertResult.Vectors
	NlpConfigInfo := convertResult.NlpConfigInfo

	// トランザクション開始
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		log.Error(err)
		return err
	}

	// NLP設定を保存（存在しない場合のみ挿入）
	nlpConfig := &entity.DBNlpConfig{
		NlpConfigInfo: NlpConfigInfo,
	}
	_, err = db.NewInsert().
		Model(nlpConfig).
		On("CONFLICT (max_token_length,overlap_token_length,model_name,model_vector_length) DO NOTHING").
		Exec(context.Background())
	if err != nil {
		log.Error(err)
		return err
	}

	// ページを保存（存在しない場合は挿入、存在する場合は更新）
	dbPage := &entity.DBPage{
		PageInfo: page,
	}
	_, err = db.NewInsert().
		Model(dbPage).
		On("CONFLICT (domain_id, path) DO UPDATE").
		Set("title = ?", page.Title).
		Set("description = ?", page.Description).
		Set("keywords = ?", page.Keywords).
		Set("markdown = ?", page.Markdown).
		Set("hash = ?", page.Hash).
		Set("updated_at = CURRENT_TIMESTAMP").
		Exec(context.Background())
	if err != nil {
		log.Error(err)
		return err
	}

	for i, chunk := range chunks {
		// チャンク情報を保存
		chunkData := model.ChunkInfo{
			Chunk:       chunk,
			PageID:      dbPage.ID,
			NlpConfigID: nlpConfig.ID,
		}
		chunkInfo := &entity.DBChunk{
			ChunkInfo: chunkData,
		}
		_, err = db.NewInsert().
			Model(chunkInfo).
			On("CONFLICT (page_id) DO NOTHING").
			Exec(context.Background())
		if err != nil {
			log.Error(err)
			return err
		}

		// ベクトル情報を保存
		vectorData := model.VectorInfo{
			Vector:      vectors[i],
			ChunkID:     chunkInfo.ID,
			NlpConfigID: nlpConfig.ID,
		}
		vectorInfo := &entity.DBVector{
			VectorInfo: vectorData,
		}
		_, err = db.NewInsert().
			Model(vectorInfo).
			On("CONFLICT (chunk_id) DO UPDATE").
			Set("vector = ?", vectors[i]).
			Set("nlp_config_id = ?", nlpConfig.ID).
			Exec(context.Background())
		if err != nil {
			log.Error(err)
			return err
		}
	}

	// トランザクションコミット
	if err := tx.Commit(); err != nil {
		log.Error(err)
		return err
	}

	// 正常終了
	log.Info(">> OK !!")

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
