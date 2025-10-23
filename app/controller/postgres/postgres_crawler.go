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

	// チャンク情報、ベクトル情報、NLP設定情報
	chunks := convertResult.Chunks
	vectors := convertResult.Vectors
	nlpConfigInfo := convertResult.NlpConfigInfo

	// トランザクション開始
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Error(err)
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// NLP設定を保存（UPDATE 文以下はコンフリクト時に既存のレコードのIDを取得するためのもの）
	nlpConfig := &entity.DBNlpConfig{NlpConfigInfo: nlpConfigInfo}
	_, err = tx.NewInsert().
		Model(nlpConfig).
		On("CONFLICT (max_token_length,overlap_token_length,model_name,model_vector_length) DO UPDATE SET max_token_length = EXCLUDED.max_token_length").
		Returning("id").
		Exec(ctx)
	if err != nil {
		log.Error(err)
		return err
	}

	// ページを保存（存在しない場合は挿入、存在する場合は更新）
	dbPage := &entity.DBPage{PageInfo: page}
	_, err = tx.NewInsert().
		Model(dbPage).
		On("CONFLICT (domain_id, path) DO UPDATE").
		Set("title = EXCLUDED.title").
		Set("description = EXCLUDED.description").
		Set("keywords = EXCLUDED.keywords").
		Set("markdown = EXCLUDED.markdown").
		Set("hash = EXCLUDED.hash").
		Set("updated_at = CURRENT_TIMESTAMP").
		Returning("id").
		Exec(ctx)
	if err != nil {
		log.Error(err)
		return err
	}

	// チャンクとベクトルを一括保存
	for i, chunk := range chunks {
		// チャンク情報を保存（UPDATE 文以下はコンフリクト時に既存のレコードのIDを取得するためのもの）
		chunkData := model.ChunkInfo{
			NlpConfigID: nlpConfig.ID,
			PageID:      dbPage.ID,
			Chunk:       chunk,
		}
		chunkInfo := &entity.DBChunk{
			ChunkInfo: chunkData,
		}
		_, err = tx.NewInsert().
			Model(chunkInfo).
			On("CONFLICT (nlp_config_id, page_id, chunk) DO UPDATE SET chunk = EXCLUDED.chunk").
			Returning("id").
			Exec(ctx)
		if err != nil {
			log.Error(err)
			return err
		}

		// ベクトル情報を保存
		vectorData := model.VectorInfo{
			NlpConfigID: nlpConfig.ID,
			ChunkID:     chunkInfo.ID,
			Vector:      vectors[i],
		}
		vectorInfo := &entity.DBVector{
			VectorInfo: vectorData,
		}
		_, err = tx.NewInsert().
			Model(vectorInfo).
			On("CONFLICT (chunk_id) DO NOTHING").
			Exec(ctx)
		if err != nil {
			log.Error(err)
			return err
		}
	}

	// トランザクションコミット
	if err = tx.Commit(); err != nil {
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
