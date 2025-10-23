// 各コントローラーへの処理をまとめ、動作単位にまとめた関数を定義するパッケージ
package usecase

import (
	"app/controller/log"
	"app/controller/nlp"
	"app/controller/postgres"
)

// 検索結果用のページ情報（ドメイン文字列を含む）
type PageWithDomain struct {
	Domain      string `json:"domain"`
	Path        string `json:"path"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Keywords    string `json:"keywords"`
	Markdown    string `json:"markdown"`
}

/*
ページデータのベクトル検索を行う関数
  - query					検索クエリ
  - resultLimit				返却する件数
  - return)	similarPages	コサイン類似度が上位のページデータ
  - return) err				エラー
*/
func VectorSearch(query string, resultLimit int) (similarPagesWithDomain []PageWithDomain, err error) {
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

	similarPages, err := postgres.GetSimilarVectors(vector, resultLimit)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// 検索結果を PageWithDomain に変換
	similarPagesWithDomain = make([]PageWithDomain, 0, len(similarPages))
	for _, page := range similarPages {
		domainStr := ""
		if page.Domain != nil {
			domainStr = page.Domain.Domain
		}
		similarPagesWithDomain = append(similarPagesWithDomain, PageWithDomain{
			Domain:      domainStr,
			Path:        page.Path,
			Title:       page.Title,
			Description: page.Description,
			Keywords:    page.Keywords,
			Markdown:    page.Markdown,
		})
	}

	return similarPagesWithDomain, nil
}
