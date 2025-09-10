// 永続化（DB等）や外部システム連携のためのデータ構造を定義するパッケージ
package entity

import (
	"app/domain/model"
	"time"

	"github.com/uptrace/bun"
)

/*
bun: は ORM が利用するタグ
json: は API のレスポンスを生成時に使用されるタグ
*/

type DBDomain struct {
	bun.BaseModel `bun:"table:domains"`

	ID        int64     `bun:"id,pk,autoincrement"`                                 // ID
	model.DomainInfo
	CreatedAt time.Time `bun:",notnull,default:current_timestamp,type:timestamptz"` // 作成日時
	UpdatedAt time.Time `bun:",notnull,default:current_timestamp,type:timestamptz"` // 更新日時
	DeletedAt time.Time `bun:",soft_delete,type:timestamptz"`                       // 削除日時
}

// DB 用ページコンテンツ情報
type DBPage struct {
	bun.BaseModel `bun:"table:pages"`

	ID        int64     `bun:"id,pk,autoincrement" json:"-"`                                          // ID
	model.PageInfo
	Domain    *DBDomain `bun:"rel:belongs-to,join:domain_id=id" json:"domain_info"`                   // ドメイン情報
	CreatedAt time.Time `bun:",notnull,default:current_timestamp,type:timestamptz" json:"-"`          // 作成日時
	UpdatedAt time.Time `bun:",notnull,default:current_timestamp,type:timestamptz" json:"updated_at"` // 更新日時
	DeletedAt time.Time `bun:",soft_delete,type:timestamptz" json:"-"`                                // 削除日時
}

// DB 用チャンク情報
type DBChunk struct {
	bun.BaseModel `bun:"table:chunks"`

	ID        int64     `bun:"id,pk,autoincrement"`                                 // ID
	model.ChunkInfo
	Page      *DBPage   `bun:"rel:belongs-to,join:page_id=id"`                      // ページ情報
	CreatedAt time.Time `bun:",notnull,default:current_timestamp,type:timestamptz"` // 作成日時
	UpdatedAt time.Time `bun:",notnull,default:current_timestamp,type:timestamptz"` // 更新日時
	DeletedAt time.Time `bun:",soft_delete,type:timestamptz"`                       // 削除日時
}

// DB 用ベクトル情報
// 一つの対象に複数モデルによって複数ベクトルが作られることが想定されるため、ベクトルのテーブルは独立しているべき
type DBVector struct {
	bun.BaseModel `bun:"table:vectors"`

	ID        int64     `bun:"id,pk,autoincrement"`                                 // ID
	model.VectorInfo
	Chunk     *DBChunk  `bun:"rel:belongs-to,join:chunk_id=id"`                     // チャンク情報
	CreatedAt time.Time `bun:",notnull,default:current_timestamp,type:timestamptz"` // 作成日時
	UpdatedAt time.Time `bun:",notnull,default:current_timestamp,type:timestamptz"` // 更新日時
	DeletedAt time.Time `bun:",soft_delete,type:timestamptz"`                       // 削除日時
}

// DB 用NLP設定情報
type DBNlpConfig struct {
	bun.BaseModel `bun:"table:nlp_configs"`

	model.NlpConfigInfo
	ID        int64     `bun:"id,pk,autoincrement"`                                 // ID
	CreatedAt time.Time `bun:",notnull,default:current_timestamp,type:timestamptz"` // 作成日時
	UpdatedAt time.Time `bun:",notnull,default:current_timestamp,type:timestamptz"` // 更新日時
	DeletedAt time.Time `bun:",soft_delete,type:timestamptz"`                       // 削除日時
}
