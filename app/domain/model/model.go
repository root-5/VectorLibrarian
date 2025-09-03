// ドメイン（業務知識）のデータ構造を定義するパッケージ
package model

/*
bun: は ORM が利用するタグ
json: は API のレスポンスを生成時に使用されるタグ
ドメイン層にはタグが無いべきだが、そうすると重複した記述が entity 側に生じ、ドメインを分離した意味が薄れるため妥協している
*/

// URLドメイン情報
type Domain struct {
	Domain string `bun:"domain,notnull,unique:url,type:varchar(100)" json:"domain"` // ドメイン
}

// ページコンテンツ情報
type Page struct {
	DomainId    int64  `bun:"domain_id,notnull,type:int"`                               // ドメインID
	Path        string `bun:"path,notnull,unique:url,type:varchar(255)" json:"path"`    // パス
	Title       string `bun:"title,notnull,type:varchar(100)" json:"title"`             // ページタイトル
	Description string `bun:"description,notnull,type:varchar(255)" json:"description"` // ディスクリプション
	Keywords    string `bun:"keywords,notnull,type:varchar(255)" json:"keywords"`       // キーワード
	Markdown    string `bun:"markdown,notnull,type:text" json:"markdown"`               // Markdown コンテンツ
	Hash        string `bun:"hash,notnull,type:char(64)" json:"-"`                      // コンテンツのハッシュ値
}

// チャンク情報
type Chunk struct {
	Chunk       string `bun:"chunk,notnull,type:text"`        // チャンク
	PageId      int64  `bun:"page_id,notnull,type:int"`       // ページID
	NlpConfigId int64  `bun:"nlp_config_id,notnull,type:int"` // NLP設定ID
}

// ベクトル情報
type Vector struct {
	Vector      []float32 `bun:"vector,notnull,type:vector(384)"` // ベクトルデータ（モデルの次元数に合わせて変更）
	ChunkId     int64     `bun:"chunk_id,notnull,type:int"`       // チャンクID
	NlpConfigId int64     `bun:"nlp_config_id,notnull,type:int"`  // NLP設定ID
}

// NLP設定情報
type NlpConfig struct {
	MaxTokenLength     int64  `bun:"max_token_length,notnull,type:int"`     // 最大トークン長
	OverlapTokenLength int64  `bun:"overlap_token_length,notnull,type:int"` // オーバーラップトークン長
	ModelName          string `bun:"model_name,notnull,type:varchar(100)"`  // モデル名
	ModelVectorLength  int64  `bun:"model_vector_length,notnull,type:int"`  // モデルのベクトル長
}

// 検索履歴情報
// type SearchLog struct {
// 	bun.BaseModel `bun:"table:headings"`
// }
