package model

import (
	"time"

	"github.com/uptrace/bun"
)

// ドメイン情報
// type Domain struct {
// 	bun.BaseModel `bun:"table:headings"`
// }

// ページコンテンツ情報
type Page struct {
	bun.BaseModel `bun:"table:pages"`

	Id          int64     `bun:"id,pk,autoincrement"`                                 // ID
	Domain      string    `bun:"domain,notnull,unique:url,type:varchar(100)"`         // ドメイン
	Path        string    `bun:"path,notnull,unique:url,type:varchar(255)"`           // パス
	Title       string    `bun:"title,notnull,type:varchar(100)"`                     // ページタイトル
	Description string    `bun:"description,notnull,type:varchar(255)"`               // ディスクリプション
	Keywords    string    `bun:"keywords,notnull,type:varchar(255)"`                  // キーワード
	Markdown    string    `bun:"markdown,notnull,type:text"`                          // Markdown コンテンツ
	Hash        string    `bun:"hash,notnull,type:char(64)"`                          // コンテンツのハッシュ値
	Vector      []float32 `bun:"vector,notnull,type:vector(384)"`                     // ベクトルデータ（モデルの次元数に合わせて変更）
	CreatedAt   time.Time `bun:",notnull,default:current_timestamp,type:timestamptz"` // 作成日時
	UpdatedAt   time.Time `bun:",notnull,default:current_timestamp,type:timestamptz"` // 更新日時
	DeletedAt   time.Time `bun:",soft_delete,type:timestamptz"`                       // 削除日時
}

// 見出し情報
// type Heading struct {
// 	bun.BaseModel `bun:"table:headings"`

// 	Id           int64     `bun:"id,pk,autoincrement"`                                 // ID
// 	PageId       int64     `bun:"page_id,notnull,type:int"`                            // ページID
// 	HeadingIndex int64     `bun:"heading_index,notnull,type:int"`                      // 見出しインデックス
// 	HeadingPath  string    `bun:"heading_path,notnull,type:varchar(255)"`              // 見出しパス (title > h1 > h2 > h3 ...)
// 	CreatedAt    time.Time `bun:",notnull,default:current_timestamp,type:timestamptz"` // 作成日時
// 	UpdatedAt    time.Time `bun:",notnull,default:current_timestamp,type:timestamptz"` // 更新日時
// 	DeletedAt    time.Time `bun:",soft_delete,type:timestamptz"`                       // 削除日時
// }

// ベクトル情報
// 一つの対象に複数モデルによって複数ベクトルが作られることが想定されるため、ベクトルのテーブルは独立しているべき
// type Embedding struct {
// 	bun.BaseModel `bun:"table:embeddings"`

// 	Id        int64     `bun:"id,pk,autoincrement"`                                 // ID
// 	ChunkId   int64     `bun:"chunk_id,notnull,type:int"`                           // チャンクID
// 	Vector    []float32 `bun:"vector,notnull,type:vector(1536)"`                    // ベクトルデータ
// 	ModelName string    `bun:"model_name,notnull,type:varchar(100)"`                // モデル名
// 	CreatedAt time.Time `bun:",notnull,default:current_timestamp,type:timestamptz"` // 作成日時
// 	UpdatedAt time.Time `bun:",notnull,default:current_timestamp,type:timestamptz"` // 更新日時
// 	DeletedAt time.Time `bun:",soft_delete,type:timestamptz"`                       // 削除日時
// }

// 検索履歴情報
// type SearchLog struct {
// 	bun.BaseModel `bun:"table:headings"`
// }
