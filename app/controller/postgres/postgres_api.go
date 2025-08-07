// PostgreSQL を利用するための関数をまとめたパッケージ
package postgres

import (
	"app/controller/log"
	"app/domain/model"
	"context"
)

/*
ベクトルを入力して、コサイン類似度が上位10件までのデータを返却する関数
  - vector		入力するベクトル
  - return)		コサイン類似度が上位10件までのページデータ
  - return) err	エラー
*/
func GetTop10SimilarPages(vector []float32) (pages []model.Page, err error) {
	err = db.NewSelect().
		Model(&pages).
		Where("SELECT *, 1 - (vector <=> ?) AS similarity FROM pages ORDER BY similarity LIMIT 10", vector).
		Scan(context.Background())
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return pages, nil
}
