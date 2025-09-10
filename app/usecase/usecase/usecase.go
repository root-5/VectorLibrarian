// 各コントローラーへの処理をまとめ、動作単位にまとめた関数を定義するパッケージ
package usecase

import (
	"app/controller/log"
	"app/controller/nlp"
	"app/controller/postgres"
	"app/domain/model"
)

/*
ページデータのベクトル検索を行う関数
  - query			検索クエリ
  - limit			返却する件数
  - return)	similarPages	コサイン類似度が上位のページデータ
  - return) err		エラー
*/
func VectorSearch(query string, limit int) (similarPages []model.VectorInfo, err error) {
	resp, _ := nlp.ConvertToVector(query, false)

	// 検索用にベクトルを一つにまとめる（平均を取る）
	vector := resp.Vectors[0]
	for i := 1; i < len(resp.Vectors); i++ {
		for j := 0; j < len(resp.Vectors[i]); j++ {
			vector[j] += resp.Vectors[i][j]
		}
	}
	for i := 0; i < len(vector); i++ {
		vector[i] /= float32(len(resp.Vectors))
	}

	similarPages, err = postgres.GetSimilarVectors(vector, limit)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return similarPages, nil
}
