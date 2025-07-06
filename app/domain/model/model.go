package model

import (
	"time"

	"github.com/uptrace/bun"
)

// ページコンテンツ情報
type Page struct {
	bun.BaseModel `bun:"table:pages"`

	Id          int64     `bun:"id,pk,autoincrement,type:bigint"`                              // ID
	Domain      string    `bun:"domain,notnull,unique:url,type:varchar(100)"`                  // ドメイン
	Path        string    `bun:"path,notnull,unique:url,type:varchar(255)"`                    // パス
	Title       string    `bun:"title,notnull,type:varchar(100)"`                              // ページタイトル
	Description string    `bun:"description,notnull,type:varchar(255)"`                        // ページ説明
	Keywords    string    `bun:"keywords,notnull,type:varchar(255)"`                           // キーワード
	Markdown    string    `bun:"markdown,notnull,type:text"`                                   // Markdown コンテンツ
	Hash        string    `bun:"hash,notnull,type:char(64)"`                                   // コンテンツのハッシュ値
	CreatedAt   time.Time `bun:",nullzero,notnull,default:current_timestamp,type:timestamptz"` // 作成日時
	UpdatedAt   time.Time `bun:",nullzero,notnull,default:current_timestamp,type:timestamptz"` // 更新日時
	DeletedAt   time.Time `bun:",soft_delete,type:timestamptz"`                                // 削除日時
}
