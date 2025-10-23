// PostgreSQL を利用するための関数をまとめたパッケージ
package postgres

import (
	"app/controller/log"
	"app/usecase/entity"
	"context"
	"fmt"
	"strings"
)

/*
ベクトルを入力して、コサイン類似度が上位のデータを指定の件数返却する関数
  - vector		入力するベクトル
  - resultLimit	返却する件数
  - return)		コサイン類似度が上位のページデータ
  - return) err	エラー
*/
func GetSimilarVectors(vector []float32, resultLimit int) (similarPages []entity.DBPage, err error) {
	vectorStr := vectorToString(vector)

	var dbVectors []entity.DBVector
	err = db.NewSelect().
		Model(&dbVectors).
		Relation("Chunk.Page.Domain").
		OrderExpr("vector <=> ?", vectorStr).
		Limit(resultLimit).
		Scan(context.Background())
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// entity.DBVector から entity.DBPage に変換（ドメイン情報を含む）
	similarPages = make([]entity.DBPage, 0, len(dbVectors))
	for _, dbVector := range dbVectors {
		if dbVector.Chunk != nil && dbVector.Chunk.Page != nil {
			similarPages = append(similarPages, *dbVector.Chunk.Page)
		}
	}

	return similarPages, nil
}

// float32スライスをPostgreSQLのベクトル形式の文字列に変換
func vectorToString(vector []float32) string {
	strSlice := make([]string, len(vector))
	for i, v := range vector {
		strSlice[i] = fmt.Sprintf("%g", v)
	}
	return "[" + strings.Join(strSlice, ",") + "]"
}
