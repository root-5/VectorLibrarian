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

	model.DomainInfo
	Id        int64     `bun:"id,pk,autoincrement"`                                 // ID
	CreatedAt time.Time `bun:",notnull,default:current_timestamp,type:timestamptz"` // 作成日時
	UpdatedAt time.Time `bun:",notnull,default:current_timestamp,type:timestamptz"` // 更新日時
	DeletedAt time.Time `bun:",soft_delete,type:timestamptz"`                       // 削除日時
}

// DB 用ページコンテンツ情報
type DBPage struct {
	bun.BaseModel `bun:"table:pages"`

	model.PageInfo
	Id        int64     `bun:"id,pk,autoincrement" json:"-"`                                          // ID
	CreatedAt time.Time `bun:",notnull,default:current_timestamp,type:timestamptz" json:"-"`          // 作成日時
	UpdatedAt time.Time `bun:",notnull,default:current_timestamp,type:timestamptz" json:"updated_at"` // 更新日時
	DeletedAt time.Time `bun:",soft_delete,type:timestamptz" json:"-"`                                // 削除日時
}

// DB 用チャンク情報
type DBChunk struct {
	bun.BaseModel `bun:"table:chunks"`

	model.ChunkInfo
	Id        int64     `bun:"id,pk,autoincrement"`                                 // ID
	CreatedAt time.Time `bun:",notnull,default:current_timestamp,type:timestamptz"` // 作成日時
	UpdatedAt time.Time `bun:",notnull,default:current_timestamp,type:timestamptz"` // 更新日時
	DeletedAt time.Time `bun:",soft_delete,type:timestamptz"`                       // 削除日時
}

// DB 用ベクトル情報
// 一つの対象に複数モデルによって複数ベクトルが作られることが想定されるため、ベクトルのテーブルは独立しているべき
type DBVector struct {
	bun.BaseModel `bun:"table:vectors"`

	model.VectorInfo
	Id        int64     `bun:"id,pk,autoincrement"`                                 // ID
	CreatedAt time.Time `bun:",notnull,default:current_timestamp,type:timestamptz"` // 作成日時
	UpdatedAt time.Time `bun:",notnull,default:current_timestamp,type:timestamptz"` // 更新日時
	DeletedAt time.Time `bun:",soft_delete,type:timestamptz"`                       // 削除日時
}

// DB 用NLP設定情報
type DBNlpConfig struct {
	bun.BaseModel `bun:"table:nlp_configs"`

	model.NlpConfigInfo
	Id        int64     `bun:"id,pk,autoincrement"`                                 // ID
	CreatedAt time.Time `bun:",notnull,default:current_timestamp,type:timestamptz"` // 作成日時
	UpdatedAt time.Time `bun:",notnull,default:current_timestamp,type:timestamptz"` // 更新日時
	DeletedAt time.Time `bun:",soft_delete,type:timestamptz"`                       // 削除日時
}
