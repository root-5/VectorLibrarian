// PostgreSQL を利用するための関数をまとめたパッケージ
package postgres

import (
	"app/controller/log"
	"app/domain/model"
	"context"
	"fmt"
	"strings"
)

/*
ベクトルを入力して、コサイン類似度が上位のデータを指定の件数返却する関数
  - vector		入力するベクトル
  - limit		返却する件数
  - return)		コサイン類似度が上位のページデータ
  - return) err	エラー
*/
func GetSimilarVectors(vector []float32, limit int) (similarVectors []model.VectorInfo, err error) {
	vectorStr := vectorToString(vector)

	err = db.NewSelect().
		Model(&similarVectors).
		Relation("Chunk.Page").
		OrderExpr("vector <=> ?", vectorStr).
		Limit(limit).
		Scan(context.Background())
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return similarVectors, nil
}

// float32スライスをPostgreSQLのベクトル形式の文字列に変換
func vectorToString(vector []float32) string {
	strSlice := make([]string, len(vector))
	for i, v := range vector {
		strSlice[i] = fmt.Sprintf("%g", v)
	}
	return "[" + strings.Join(strSlice, ",") + "]"
}
