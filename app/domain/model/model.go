// ドメイン（業務知識）のデータ構造を定義するパッケージ
package model

/*
## ドメイン分離の考え方
データ構造のうち、ドメイン領域のものを domain/model、そうでないものを usecase/entity に分離する構成をとっている。

## Page に DomainId を持たせることについて
ID という存在そのものはドメイン領域ではないが、Page と Domain の関係性はドメイン領域であり、ドメイン領域でも関係性を明示できるためしている。
ドメイン層には本来タグが無いべきだが、そうすると重複した記述が entity 側に生じ、ドメインを分離した意味が薄れるため妥協している。


bun: は ORM が利用するタグ
json: は API のレスポンスを生成時に使用されるタグ
*/

// URLドメイン情報
type DomainInfo struct {
	Domain string `bun:"domain,notnull,unique,type:varchar(100)" json:"domain"` // ドメイン
}

// ページコンテンツ情報
type PageInfo struct {
	DomainID    int64  `bun:"domain_id,notnull,unique"`                                 // ドメインID
	Path        string `bun:"path,notnull,unique,type:varchar(255)" json:"path"`        // パス
	Title       string `bun:"title,notnull,type:varchar(100)" json:"title"`             // ページタイトル
	Description string `bun:"description,notnull,type:varchar(255)" json:"description"` // ディスクリプション
	Keywords    string `bun:"keywords,notnull,type:varchar(255)" json:"keywords"`       // キーワード
	Markdown    string `bun:"markdown,notnull,type:text" json:"markdown"`               // Markdown コンテンツ
	Hash        string `bun:"hash,notnull,type:char(64)" json:"-"`                      // コンテンツのハッシュ値
}

// チャンク情報
type ChunkInfo struct {
	Chunk       string `bun:"chunk,notnull,type:text"` // チャンク
	PageID      int64  `bun:"page_id,notnull"`         // ページID
	NlpConfigID int64  `bun:"nlp_config_id,notnull"`   // NLP設定ID
}

// ベクトル情報
type VectorInfo struct {
	Vector      []float32 `bun:"vector,notnull,type:vector(384)"` // ベクトルデータ（モデルの次元数に合わせて変更）
	ChunkID     int64     `bun:"chunk_id,notnull"`                // チャンクID
	NlpConfigID int64     `bun:"nlp_config_id,notnull"`           // NLP設定ID
}

// NLP設定情報
type NlpConfigInfo struct {
	MaxTokenLength     int64  `bun:"max_token_length,unique:configname,notnull" json:"max_token_length"`
	OverlapTokenLength int64  `bun:"overlap_token_length,unique:configname,notnull" json:"overlap_token_length"`
	ModelName          string `bun:"model_name,unique:configname,notnull,type:varchar(100)" json:"model_name"`
	ModelVectorLength  int64  `bun:"model_vector_length,unique:configname,notnull" json:"model_vector_length"`
}

// 検索履歴情報
// type SearchLog struct {
// 	bun.BaseModel `bun:"table:headings"`
// }
